// This file contains a way to read binary data from the FFT ISO.
// It should be expanded to also read the FFT bin file.
//
// It contains the low level methods for different sized ints/uints as well has
// some simple geometry parsing. The higher level iso parsing happens in mesh.go.
// The split is somewhat arbitrary.
package main

import (
	"encoding/binary"
	"log"
	"math"
	"os"
)

const (
	sectorSize int64 = 2048

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
	to := sector * sectorSize
	_, err := r.file.Seek(to, 0)
	if err != nil {
		log.Fatalf("seek to sector: %v", err)
	}
}

// seekPointer will seek to the specified sector, plus a little more, of the iso
// file. This is useful when using MeshFileHeader intra-file pointers.
func (r Reader) seekPointer(sector int64, ptr int64) {
	to := sector*sectorSize + ptr
	_, err := r.file.Seek(to, 0)
	if err != nil {
		log.Fatalf("seek to pointer: %v", err)
	}
}

func (r Reader) readUint8() uint8 {
	size := 1
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return data[0]
}

func (r Reader) readUint16() uint16 {
	size := 2
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint16(data)
}

func (r Reader) readUint32() uint32 {
	size := 4
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint32(data)
}

func (r Reader) readInt8() int8   { return int8(r.readUint8()) }
func (r Reader) readInt16() int16 { return int16(r.readUint16()) }
func (r Reader) readInt32() int32 { return int32(r.readUint32()) }

func (r Reader) readVertex() Vec3 {
	x := float64(r.readInt16())
	y := float64(r.readInt16())
	z := float64(r.readInt16())
	return Vec3{x: x, y: -y, z: z}
}

func (r Reader) readTriangle() Triangle {
	a := r.readVertex()
	b := r.readVertex()
	c := r.readVertex()
	return Triangle{vertices: [3]Vec3{a, b, c}, color: White}
}

func (r Reader) readQuad() []Triangle {
	a := r.readVertex()
	b := r.readVertex()
	c := r.readVertex()
	d := r.readVertex()
	return []Triangle{
		{vertices: [3]Vec3{a, b, c}, color: White},
		{vertices: [3]Vec3{b, d, c}, color: White},
	}
}

func (r Reader) readNormal() Vec3 {
	x := r.readF1x3x12()
	y := r.readF1x3x12()
	z := r.readF1x3x12()
	return Vec3{x: x, y: -y, z: z}
}

func (r Reader) readTriNormal() []Vec3 {
	a := r.readNormal()
	b := r.readNormal()
	c := r.readNormal()
	return []Vec3{a, b, c}
}

func (r Reader) readQuadNormal() [][]Vec3 {
	a := r.readNormal()
	b := r.readNormal()
	c := r.readNormal()
	d := r.readNormal()
	return [][]Vec3{
		{a, b, c},
		{b, d, c},
	}
}

func (r Reader) readUV() Tex {
	x := float64(r.readUint8())
	y := float64(r.readUint8())
	return Tex{U: x, V: y}
}

func (r Reader) readTriUV() ([3]Tex, int) {
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

func (r Reader) readQuadUV() ([2][3]Tex, int) {
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

func (r Reader) readF1x3x12() float64 {
	return float64(r.readInt16()) / 4096.0
}

func (r Reader) readRGB8() Color {
	return Color{
		R: r.readUint8(),
		G: r.readUint8(),
		B: r.readUint8(),
		A: 255,
	}
}

func (mr Reader) readRGB15() Color {
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

func (r Reader) readLightColor() uint8 {
	val := r.readF1x3x12()
	return uint8(255 * math.Min(math.Max(0.0, val), 1.0))
}

func (r Reader) readDirectionalLights() []DirectionalLight {
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

func (r Reader) readAmbientLight() AmbientLight {
	color := r.readRGB8()
	return AmbientLight{color: color}

}

func (r Reader) readBackground() Background {
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

//
// Mesh File Header
//

// Table of pointers contained in the meshFileHeader.
const (
	ptrPrimaryMesh          = 0x40
	ptrTexturePalettesColor = 0x44
	ptrUnknown              = 0x4c // Only non-zero in MAP000.5
	ptrLightsAndBackground  = 0x64 // Light colors/positions, bg gradient colors
	ptrTerrain              = 0x68 // Tile heights, slopes, and surface types
	ptrTextureAnimInst      = 0x6c
	ptrPaletteAnimInst      = 0x70
	ptrTexturePalettesGray  = 0x7c
	ptrMeshAnimInst         = 0x8c
	ptrAnimatedMesh1        = 0x90
	ptrAnimatedMesh2        = 0x94
	ptrAnimatedMesh3        = 0x98
	ptrAnimatedMesh4        = 0x9c
	ptrAnimatedMesh5        = 0xa0
	ptrAnimatedMesh6        = 0xa4
	ptrAnimatedMesh7        = 0xa8
	ptrAnimatedMesh8        = 0xac
	ptrVisibilityAngles     = 0xb0
)

// meshFileHeader contains 32-bit unsigned little-endian pointers to an area of
// the mesh data. Zero is returned if there is no pointer.
type meshFileHeader []byte

// meshFileHeaderLen is the length in bytes.
const meshFileHeaderLen = 196

// Return the intra-file pointer for different parts of the mesh data.
// All pointers are converted to int64 since thats what seek functions take
func (h meshFileHeader) ptr(location int32) int64 {
	const ptrLen = 4 // Intra-file pointers are always 32bit
	return int64(binary.LittleEndian.Uint32(h[location : location+ptrLen]))
}

func (h meshFileHeader) PrimaryMesh() int64          { return h.ptr(ptrPrimaryMesh) }
func (h meshFileHeader) TexturePalettesColor() int64 { return h.ptr(ptrTexturePalettesColor) }
func (h meshFileHeader) Unknown() int64              { return h.ptr(ptrUnknown) }
func (h meshFileHeader) LightsAndBackground() int64  { return h.ptr(ptrLightsAndBackground) }
func (h meshFileHeader) Terrain() int64              { return h.ptr(ptrTerrain) }
func (h meshFileHeader) TextureAnimInst() int64      { return h.ptr(ptrTextureAnimInst) }
func (h meshFileHeader) PaletteAnimInst() int64      { return h.ptr(ptrPaletteAnimInst) }
func (h meshFileHeader) TexturePalettesGray() int64  { return h.ptr(ptrTexturePalettesGray) }
func (h meshFileHeader) MeshAnimInst() int64         { return h.ptr(ptrMeshAnimInst) }
func (h meshFileHeader) AnimatedMesh1() int64        { return h.ptr(ptrAnimatedMesh1) }
func (h meshFileHeader) AnimatedMesh2() int64        { return h.ptr(ptrAnimatedMesh2) }
func (h meshFileHeader) AnimatedMesh3() int64        { return h.ptr(ptrAnimatedMesh3) }
func (h meshFileHeader) AnimatedMesh4() int64        { return h.ptr(ptrAnimatedMesh4) }
func (h meshFileHeader) AnimatedMesh5() int64        { return h.ptr(ptrAnimatedMesh5) }
func (h meshFileHeader) AnimatedMesh6() int64        { return h.ptr(ptrAnimatedMesh6) }
func (h meshFileHeader) AnimatedMesh7() int64        { return h.ptr(ptrAnimatedMesh7) }
func (h meshFileHeader) AnimatedMesh8() int64        { return h.ptr(ptrAnimatedMesh8) }
func (h meshFileHeader) VisibilityAngles() int64     { return h.ptr(ptrVisibilityAngles) }

//
// Mesh Header
//

// meshHeader contains the number of triangles and quads, textured and
// untextured. The numbers are represented by 16-bit unsigned integers.
//
// These method names have been used as a references to FFHacktics naming.
type meshHeader []byte

// meshHeaderLen is the length in bytes.
const meshHeaderLen = 8

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
		if record.Type() == RecordTypeTexture {
			texture := r.parseTexture(record)
			textures = append(textures, texture)
		} else if record.Type() == RecordTypeMeshPrimary {
			mesh = r.parseMesh(record)
		} else if record.Type() == RecordTypeMeshAlt {
			// Sometimes there is no primary mesh (ie MAP002.GNS),
			// there is only an alternate. I'm not sure why. So we
			// treat this one as the primary, only if the primary
			// hasn't been set. Kinda Hacky until we start treating
			// each GNS Record as a Scenario.
			if len(mesh.triangles) == 0 {
				mesh = r.parseMesh(record)
			}
		}
	}

	mesh.scale = Vec3{x: 5, y: 5, z: 5}
	if len(textures) > 0 {
		mesh.texture = textures[0]
	}

	mesh.normalizeCoordinates()
	mesh.centerCoordinates()
	return mesh
}

// parseTexture reads and returns an FFT texture as an engine Texture.
func (r Reader) parseTexture(record GNSRecord) Texture {
	r.seekSector(record.Sector())
	data := make([]byte, record.Len())
	n, err := r.file.Read(data)
	if err != nil || int64(n) != record.Len() {
		log.Fatalf("read texture data: %v", err)
	}
	pixels := textureSplitPixels(data)
	return NewTexture(textureWidth, textureHeight, pixels)
}

// parseMesh reads mesh data for primary and alternate meshes.
//
func (r Reader) parseMesh(record GNSRecord) Mesh {
	r.seekSector(record.Sector())

	// File header contains intra-file pointers to areas of mesh data.
	fileHeader := make(meshFileHeader, meshFileHeaderLen)
	n, err := r.file.Read(fileHeader)
	if err != nil || int64(n) != meshFileHeaderLen {
		log.Fatalf("read mesh file header: %v", err)
	}

	// Primary mesh pointer tells us where the primary mesh data is.  I
	// think this is always 196 as it starts directly after the header,
	// which has a size of 196. Keep dynamic here as pointer access will
	// grow and this keeps it consistent.
	primaryMeshPointer := fileHeader.PrimaryMesh()

	// Previously we did these pointer checks on every map. But some maps
	// (ie MAP002.GNS) don't have a primary mesh, only alternative. The location of that mesh
	if record.Type() == RecordTypeMeshPrimary {
		if primaryMeshPointer == 0 || primaryMeshPointer != 196 {
			log.Fatal("missing primary mesh pointer")
		}
	}

	// Skip ahead to color palettes
	r.seekPointer(record.Sector(), fileHeader.TexturePalettesColor())

	palettes := make([]Palette, 16)
	for i := 0; i < 16; i++ {
		palette := make(Palette, 16)
		for j := 0; j < 16; j++ {
			palette[j] = r.readRGB15()
		}
		palettes[i] = palette
	}

	// Seek to the mesh data.
	r.seekPointer(record.Sector(), primaryMeshPointer)

	// Mesh header contains the number of triangles and quads that exist.
	header := make(meshHeader, meshHeaderLen)
	n, err = r.file.Read(header)
	if err != nil || int64(n) != meshHeaderLen {
		log.Fatalf("read mesh file header: %v", err)
	}

	// FIXME: Change capacity from TT to total with untextured.
	triangles := make([]Triangle, 0, header.TT())
	for i := 0; i < header.N(); i++ {
		triangles = append(triangles, r.readTriangle())
	}
	for i := 0; i < header.P(); i++ {
		triangles = append(triangles, r.readQuad()...)
	}
	for i := 0; i < header.Q(); i++ {
		triangles = append(triangles, r.readTriangle())
	}
	for i := 0; i < header.R(); i++ {
		triangles = append(triangles, r.readQuad()...)
	}

	// Normals
	// Nothing is actually collected. They are just read here so the
	// iso read position moves forward, so we can read polygon texture data
	// next.  This could be cleaned up as a seek, but we may eventually use
	// the normal data here.
	for i := 0; i < header.N(); i++ {
		r.readTriNormal()
	}
	for i := 0; i < header.P(); i++ {
		r.readQuadNormal()
	}

	// Polygon texture data
	for i := 0; i < header.N(); i++ {
		uv, palette := r.readTriUV()
		triangles[i].texcoords = uv
		triangles[i].palette = palettes[palette]
	}
	for i := header.N(); i < header.TT(); i = i + 2 {
		uvs, palette := r.readQuadUV()
		triangles[i].texcoords = uvs[0]
		triangles[i].palette = palettes[palette]

		triangles[i+1].texcoords = uvs[1]
		triangles[i+1].palette = palettes[palette]
	}

	// Skip ahead to lights
	r.seekPointer(record.Sector(), fileHeader.LightsAndBackground())

	directionalLights := r.readDirectionalLights()
	ambientLight := r.readAmbientLight()
	background := r.readBackground()

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

// textureSplitPixels takes the ISO's raw bytes and splits each of them into two
// bytes. The ISO has two pixels per byte to save space. We want each pixel independent,
// so we split them here. The pixel values are just an index into a color palette so the
// values are 0-15.
func textureSplitPixels(buf []byte) []Color {
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
