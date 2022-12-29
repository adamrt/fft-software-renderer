// This file contains a way to read binary data from the FFT bin file.
package main

import (
	"encoding/binary"
	"log"
	"math"
	"os"
)

const (
	sectorSize       = 2048
	sectorRawSize    = 2352
	sectorHeaderSize = 24

	// This contains the pointers below.
	meshFileHeaderLen = 196
	// This contains the info about the geometry section of the mesh.
	meshHeaderLen = 8

	// Each of these values are the location of an intra-file pointer (32-bit unsigned
	// little-endian).  The location's value will be the offset in bytes to the
	// beginning of that data.  The mesh file If the locations value is a zero there
	// is not intra-file data for that type.
	locPrimaryMesh          = 0x40
	locTexturePalettesColor = 0x44
	locUnknown              = 0x4c // Only non-zero in MAP000.5
	locLightsAndBackground  = 0x64 // Light colors/positions, bg gradient colors
	locTerrain              = 0x68 // Tile heights, slopes, and surface types
	locTextureAnimInst      = 0x6c
	locPaletteAnimInst      = 0x70
	locTexturePalettesGray  = 0x7c
	locMeshAnimInst         = 0x8c
	locAnimatedMesh1        = 0x90
	locAnimatedMesh2        = 0x94
	locAnimatedMesh3        = 0x98
	locAnimatedMesh4        = 0x9c
	locAnimatedMesh5        = 0xa0
	locAnimatedMesh6        = 0xa4
	locAnimatedMesh7        = 0xa8
	locAnimatedMesh8        = 0xac
	locVisibilityAngles     = 0xb0

	// These are FFT texture specific.
	textureWidth  int = 256
	textureHeight int = 1024
	textureLen    int = textureWidth * textureHeight
	textureRawLen int = textureLen / 2
)

type Reader struct {
	file *os.File
}

func NewReader(filename string) *Reader {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return &Reader{f}
}

func (r Reader) Close() {
	r.file.Close()
}

// seekSector will seek to the specified sector of the file.
func (r Reader) seekSector(sector int64) {
	to := sector*sectorRawSize + sectorHeaderSize
	_, err := r.file.Seek(to, 0)
	if err != nil {
		log.Fatalf("seek to sector: %v", err)
	}
}

func (r Reader) readSector(sector int64) []byte {
	r.seekSector(sector)
	data := make([]byte, sectorSize)
	n, err := r.file.Read(data)
	if err != nil || n != sectorSize {
		log.Fatal("failed to read sector data", err)
	}
	return data
}

func (r Reader) readFile(sector int64, size int64) []byte {
	occupiedSectors := int64(math.Ceil(float64(size) / float64(sectorSize)))
	data := make([]byte, 0)
	for i := int64(0); i < occupiedSectors; i++ {
		sectorData := r.readSector(sector + i)
		data = append(data, sectorData...)
	}
	return data[0:size]
}

//
// Mesh File Header
//

//
// Mesh Header
//

// meshHeader contains the number of triangles and quads, textured and
// untextured. The numbers are represented by 16-bit unsigned integers.
//
// These method names have been used as a references to FFHacktics naming.
type meshHeader []byte

func (h meshHeader) N() int {
	return int(binary.LittleEndian.Uint16(h[0:2]))
}

func (h meshHeader) P() int {
	return int(binary.LittleEndian.Uint16(h[2:4]))
}

func (h meshHeader) Q() int {
	return int(binary.LittleEndian.Uint16(h[4:6]))
}

func (h meshHeader) R() int {
	return int(binary.LittleEndian.Uint16(h[6:8]))
}

// totalTexTris returns the count of all textured triangles after quads have
// been split.
func (h meshHeader) TT() int {
	return h.N() + h.P()*2
}

func (r Reader) ReadMesh(mapNum int) Mesh {
	records := r.readGNSRecords(mapNum)

	textures := []Texture{}
	var mesh Mesh

	for _, record := range records {
		switch record.Type() {
		case RecordTypeTexture:
			textures = append(textures, r.parseTexture(record))
		case RecordTypeMeshPrimary:
			mesh = r.parseMesh(record)
		case RecordTypeMeshOverride:
			// Sometimes there is no primary mesh (ie MAP002.GNS), there is
			// only an override. Usually a non-battle map. So we treat this
			// one as the primary, only if the primary hasn't been set. Kinda
			// Hacky until we start treating each GNS Record as a Scenario.
			if len(mesh.triangles) == 0 {
				mesh = r.parseMesh(record)
			}
		}
	}

	if len(textures) > 0 {
		mesh.texture = textures[0]
	}

	s := 15.0
	mesh.scale = Vec3{s, s, s}
	// mesh.normalizeCoordinates()
	mesh.centerCoordinates()
	return mesh
}

// parseTexture reads and returns an FFT texture as an engine Texture.
func (r Reader) parseTexture(record GNSRecord) Texture {
	data := r.readFile(record.Sector(), record.Len())
	pixels := textureSplitPixels(data)
	return NewTexture(textureWidth, textureHeight, pixels)
}

// parseMesh reads mesh data for primary and alternate meshes.
func (r Reader) parseMesh(record GNSRecord) Mesh {
	data := r.readFile(record.Sector(), record.Len())
	f := MeshFile{data, 0}

	// Primary mesh pointer tells us where the primary mesh data is.  I
	// think this is always 196 as it starts directly after the header,
	// which has a size of 196. Keep dynamic here as pointer access will
	// grow and this keeps it consistent.
	primaryMeshPointer := f.PtrPrimaryMesh()

	// Previously we did these pointer checks on every map. But some maps
	// (ie MAP002.GNS) don't have a primary mesh, only alternative. The location of that mesh
	if record.Type() == RecordTypeMeshPrimary {
		if primaryMeshPointer == 0 || primaryMeshPointer != 196 {
			log.Fatal("missing primary mesh pointer")
		}
	}

	// Skip ahead to color palettes
	f.seekPointer(f.PtrTexturePalettesColor())

	palettes := make([]Palette, 16)
	for i := 0; i < 16; i++ {
		palette := make(Palette, 16)
		for j := 0; j < 16; j++ {
			palette[j] = f.readRGB15()
		}
		palettes[i] = palette
	}

	// Seek to the mesh data.
	f.seekPointer(primaryMeshPointer)

	// Mesh header contains the number of triangles and quads that exist.
	header := meshHeader(f.data[f.offset : f.offset+meshHeaderLen])
	f.offset += meshHeaderLen

	// FIXME: Change capacity from TT to total with untextured.
	triangles := make([]Triangle, 0, header.TT())
	for i := 0; i < header.N(); i++ {
		triangles = append(triangles, f.readTriangle())
	}
	for i := 0; i < header.P(); i++ {
		triangles = append(triangles, f.readQuad()...)
	}
	for i := 0; i < header.Q(); i++ {
		triangles = append(triangles, f.readTriangle())
	}
	for i := 0; i < header.R(); i++ {
		triangles = append(triangles, f.readQuad()...)
	}

	// Normals
	// Nothing is actually collected. They are just read here so the
	// iso read position moves forward, so we can read polygon texture data
	// next.  This could be cleaned up as a seek, but we may eventually use
	// the normal data here.
	for i := 0; i < header.N(); i++ {
		triangles[i].normals = f.readTriNormal()
	}
	for i := header.N(); i < header.TT(); i = i + 2 {
		qns := f.readQuadNormal()
		triangles[i].normals = qns[0]
		triangles[i+1].normals = qns[1]
	}

	// Polygon texture data
	for i := 0; i < header.N(); i++ {
		uv, palette := f.readTriUV()
		triangles[i].texcoords = uv
		triangles[i].palette = palettes[palette]
	}
	for i := header.N(); i < header.TT(); i = i + 2 {
		uvs, palette := f.readQuadUV()
		triangles[i].texcoords = uvs[0]
		triangles[i].palette = palettes[palette]

		triangles[i+1].texcoords = uvs[1]
		triangles[i+1].palette = palettes[palette]
	}

	// Skip ahead to lights
	f.seekPointer(f.PtrLightsAndBackground())

	directionalLights := f.readDirectionalLights()
	ambientLight := f.readAmbientLight()
	background := f.readBackground()

	return Mesh{
		triangles:         triangles,
		ambientLight:      ambientLight,
		directionalLights: directionalLights,
		background:        background,
	}
}

func (r Reader) readGNSRecords(mapNum int) []GNSRecord {
	sector := GNSSectors[mapNum]
	r.seekSector(sector)

	records := []GNSRecord{}
	for {
		record := make(GNSRecord, GNSRecordLen)
		n, err := r.file.Read(record)
		if err != nil || n != GNSRecordLen {
			log.Fatalf("read gns record: %v", err)
		}
		if record.Type() == RecordTypeEnd {
			break
		}
		records = append(records, record)
	}
	return records
}

type MeshFile struct {
	data   []byte
	offset int64
}

func (f *MeshFile) readUint8() uint8 {
	var size int64 = 1
	data := f.data[f.offset]
	f.offset += size
	return data
}

func (f *MeshFile) readUint16() uint16 {
	var size int64 = 2
	value := binary.LittleEndian.Uint16(f.data[f.offset : f.offset+size])
	f.offset += size
	return value
}

func (f *MeshFile) readUint32() uint32 {
	var size int64 = 4
	value := binary.LittleEndian.Uint32(f.data[f.offset : f.offset+size])
	f.offset += size
	return value
}

func (f *MeshFile) readInt8() int8   { return int8(f.readUint8()) }
func (f *MeshFile) readInt16() int16 { return int16(f.readUint16()) }
func (f *MeshFile) readInt32() int32 { return int32(f.readUint32()) }

// seekPointer will set the offset to an intrafile pointer. This is useful when using
// MeshFile intra-file pointers.
func (f *MeshFile) seekPointer(ptr int64) {
	f.offset = ptr
}

// Return the intra-file pointer for different parts of the mesh data.
// All pointers are converted to int64 since thats what seek functions take
func (h MeshFile) pointer(loc int32) int64 {
	return int64(binary.LittleEndian.Uint32(h.data[loc : loc+4]))
}

func (h MeshFile) PtrPrimaryMesh() int64          { return h.pointer(locPrimaryMesh) }
func (h MeshFile) PtrTexturePalettesColor() int64 { return h.pointer(locTexturePalettesColor) }
func (h MeshFile) PtrUnknown() int64              { return h.pointer(locUnknown) }
func (h MeshFile) PtrLightsAndBackground() int64  { return h.pointer(locLightsAndBackground) }
func (h MeshFile) PtrTerrain() int64              { return h.pointer(locTerrain) }
func (h MeshFile) PtrTextureAnimInst() int64      { return h.pointer(locTextureAnimInst) }
func (h MeshFile) PtrPaletteAnimInst() int64      { return h.pointer(locPaletteAnimInst) }
func (h MeshFile) PtrTexturePalettesGray() int64  { return h.pointer(locTexturePalettesGray) }
func (h MeshFile) PtrMeshAnimInst() int64         { return h.pointer(locMeshAnimInst) }
func (h MeshFile) PtrAnimatedMesh1() int64        { return h.pointer(locAnimatedMesh1) }
func (h MeshFile) PtrAnimatedMesh2() int64        { return h.pointer(locAnimatedMesh2) }
func (h MeshFile) PtrAnimatedMesh3() int64        { return h.pointer(locAnimatedMesh3) }
func (h MeshFile) PtrAnimatedMesh4() int64        { return h.pointer(locAnimatedMesh4) }
func (h MeshFile) PtrAnimatedMesh5() int64        { return h.pointer(locAnimatedMesh5) }
func (h MeshFile) PtrAnimatedMesh6() int64        { return h.pointer(locAnimatedMesh6) }
func (h MeshFile) PtrAnimatedMesh7() int64        { return h.pointer(locAnimatedMesh7) }
func (h MeshFile) PtrAnimatedMesh8() int64        { return h.pointer(locAnimatedMesh8) }
func (h MeshFile) PtrVisibilityAngles() int64     { return h.pointer(locVisibilityAngles) }

func (r *MeshFile) readVertex() Vec3 {
	// Normals and light direction need to be normalized after this.
	x := float64(r.readInt16()) / 4096.0
	y := float64(r.readInt16()) / 4096.0
	z := float64(r.readInt16()) / 4096.0
	return Vec3{x: x, y: -y, z: z}
}

func (r *MeshFile) readTriangle() Triangle {
	a := r.readVertex()
	b := r.readVertex()
	c := r.readVertex()
	return Triangle{vertices: [3]Vec3{a, b, c}, color: White}
}

func (r *MeshFile) readQuad() []Triangle {
	a := r.readVertex()
	b := r.readVertex()
	c := r.readVertex()
	d := r.readVertex()
	return []Triangle{
		{vertices: [3]Vec3{a, b, c}, color: White},
		{vertices: [3]Vec3{b, d, c}, color: White},
	}
}

func (r *MeshFile) readNormal() Vec3 {
	x := r.readF1x3x12()
	y := r.readF1x3x12()
	z := r.readF1x3x12()
	return Vec3{x: x, y: -y, z: z}
}

func (r *MeshFile) readTriNormal() [3]Vec3 {
	a := r.readNormal()
	b := r.readNormal()
	c := r.readNormal()
	return [3]Vec3{a, b, c}
}

func (r *MeshFile) readQuadNormal() [][3]Vec3 {
	a := r.readNormal()
	b := r.readNormal()
	c := r.readNormal()
	d := r.readNormal()
	return [][3]Vec3{
		{a, b, c},
		{b, d, c},
	}
}

func (r *MeshFile) readUV() Tex {
	x := float64(r.readUint8())
	y := float64(r.readUint8())
	return Tex{U: x, V: y}
}

func (r *MeshFile) readTriUV() ([3]Tex, int) {
	a := r.readUV()
	palette := int(r.readUint8() & 0b1111)
	r.readUint8() // padding
	b := r.readUV()
	page := int(r.readUint8() & 0b11) // only 2 bits
	r.readUint8()                     // padding
	c := r.readUV()

	a = processTexCoords(a, page)
	b = processTexCoords(b, page)
	c = processTexCoords(c, page)

	return [3]Tex{a, b, c}, palette
}

func (r *MeshFile) readQuadUV() ([2][3]Tex, int) {
	a := r.readUV()
	palette := int(r.readUint8() & 0b1111)
	r.readUint8() // padding
	b := r.readUV()
	page := int(r.readUint8() & 0b11) // only 2 bits
	r.readUint8()                     // padding
	c := r.readUV()
	d := r.readUV()

	a = processTexCoords(a, page)
	b = processTexCoords(b, page)
	c = processTexCoords(c, page)
	d = processTexCoords(d, page)

	return [2][3]Tex{{a, b, c}, {b, d, c}}, palette

}

func (r *MeshFile) readF1x3x12() float64 {
	return float64(r.readInt16()) / 4096.0
}

func (r *MeshFile) readRGB8() Color {
	return Color{
		R: r.readUint8(),
		G: r.readUint8(),
		B: r.readUint8(),
		A: 255,
	}
}

func (mr *MeshFile) readRGB15() Color {
	val := mr.readUint16()
	var a uint8
	if val == 0 {
		a = 0x00
	} else {
		a = 0xFF
	}

	b := uint8((val & 0b01111100_00000000) >> 7)
	g := uint8((val & 0b00000011_11100000) >> 2)
	r := uint8((val & 0b00000000_00011111) << 3)
	return Color{R: r, G: g, B: b, A: a}
}

func (r *MeshFile) readLightColor() uint8 {
	val := r.readF1x3x12()
	return uint8(255 * math.Min(math.Max(0.0, val), 1.0))
}

func (r *MeshFile) readDirectionalLights() []DirectionalLight {
	l1r, l2r, l3r := r.readLightColor(), r.readLightColor(), r.readLightColor()
	l1g, l2g, l3g := r.readLightColor(), r.readLightColor(), r.readLightColor()
	l1b, l2b, l3b := r.readLightColor(), r.readLightColor(), r.readLightColor()

	l1color := Color{R: l1r, G: l1g, B: l1b, A: 255}
	l2color := Color{R: l2r, G: l2g, B: l2b, A: 255}
	l3color := Color{R: l3r, G: l3g, B: l3b, A: 255}

	l1pos, l2pos, l3pos := r.readVertex(), r.readVertex(), r.readVertex()

	return []DirectionalLight{
		{position: l1pos, color: l1color},
		{position: l2pos, color: l2color},
		{position: l3pos, color: l3color},
	}
}

func (r *MeshFile) readAmbientLight() AmbientLight {
	color := r.readRGB8()
	return AmbientLight{color: color}

}

func (r *MeshFile) readBackground() Background {
	top := r.readRGB8()
	bottom := r.readRGB8()
	return Background{Top: top, Bottom: bottom}
}

// processTexCoords has two functions:
//
// 1. Update the V coordinate to the specific page of the texture. FFT Textures
// have 4 pages (256x1024) and the original V specifies the pixel on one of the
// 4 pages. Multiple the page by the height of a single page (256).
//
// 2. Normalize the coordinates that can be U: 0-255 and V: 0-1023. Just divide
// them by their max to get a 0.0-1.0 value.
func processTexCoords(uv Tex, page int) Tex {
	v := float64(int(uv.V) + page*256)
	return Tex{U: uv.U / 255, V: v / 1023.0}
}

// textureSplitPixels takes the ISO's raw bytes and splits each of them into two
// bytes. The ISO has two pixels per byte to save space. We want each pixel independent,
// so we split them here. The pixel values are just an index into a color palette so the
// values are 0-15.
func textureSplitPixels(buf []uint8) []Color {
	data := make([]Color, 0)
	for i := 0; i < textureRawLen; i++ {
		colorA := uint8(buf[i] & 0x0F)
		colorB := uint8((buf[i] & 0xF0) >> 4)

		// We dont care about RGB here.
		// This is just an index to the palette.
		data = append(data,
			Color{R: colorA, G: colorA, B: colorA, A: 255},
			Color{R: colorB, G: colorB, B: colorB, A: 255},
		)
	}
	return data
}
