package main

import "math"

type Tex struct {
	u, v float64
}

type Triangle struct {
	vertices  [3]Vec3
	normals   [3]Vec3
	texcoords [3]Tex
	palette   Palette

	// Computed during render
	points     [3]Vec2
	avgDepth   float64
	lightColor Color

	// Color of untextured triangle
	color Color
}

type Mesh struct {
	triangles   []Triangle
	texture     Texture
	scale       Vec3
	rotation    Vec3
	translation Vec3

	ambientLight      AmbientLight
	directionalLights []DirectionalLight
	background        Background
}

func NewMesh() Mesh {
	return Mesh{scale: Vec3{1, 1, 1}}
}

// centerTranslation returns a translation vector that will center the mesh.
func (m *Mesh) coordCenter() Vec3 {
	var minx float64 = math.MaxInt16
	var maxx float64 = math.MinInt16
	var miny float64 = math.MaxInt16
	var maxy float64 = math.MinInt16
	var minz float64 = math.MaxInt16
	var maxz float64 = math.MinInt16

	for _, t := range m.triangles {
		for i := 0; i < 3; i++ {
			// Min
			minx = math.Min(t.vertices[i].x, minx)
			miny = math.Min(t.vertices[i].y, miny)
			minz = math.Min(t.vertices[i].z, minz)
			// Max
			maxx = math.Max(t.vertices[i].x, maxx)
			maxy = math.Max(t.vertices[i].y, maxy)
			maxz = math.Max(t.vertices[i].z, maxz)
		}
	}

	x := (maxx + minx) / 2.0
	y := (maxy + miny) / 2.0
	z := (maxz + minz) / 2.0

	return Vec3{x, y, z}
}
