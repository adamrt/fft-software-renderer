package main

import "math"

type Tex struct {
	U, V float64
}

type Triangle struct {
	vertices [3]Vec3
	points   [3]Vec2

	texcoords [3]Tex
	palette   Palette

	// color of untextured triangle
	color    Color
	avgDepth float64
}

type Mesh struct {
	triangles   []Triangle
	texture     Texture
	scale       Vec3
	rotation    Vec3
	translation Vec3
}

func NewMesh() Mesh {
	return Mesh{scale: Vec3{1, 1, 1}}
}

// normalizeCoordinates normalizes all vertex coordinates between 0 and 1. This
// scales down large models during import.  This is primary used for loading FFT
// maps since they have very large coordinates.  The min/max values should be
// the min and max of
func (m *Mesh) normalizeCoordinates() {
	min, max := m.coordMinMax()
	for i := 0; i < len(m.triangles); i++ {
		for j := 0; j < 3; j++ {
			m.triangles[i].vertices[j].x = normalize(m.triangles[i].vertices[j].x, min, max)
			m.triangles[i].vertices[j].y = normalize(m.triangles[i].vertices[j].y, min, max)
			m.triangles[i].vertices[j].z = normalize(m.triangles[i].vertices[j].z, min, max)
		}
	}
}

// centerCoordinates transforms all coordinates so the center of the model is at
// the origin point.
func (m *Mesh) centerCoordinates() {
	vec3 := m.coordCenter()
	matrix := MatrixTranslation(vec3)
	for i := 0; i < len(m.triangles); i++ {
		for j := 0; j < 3; j++ {
			transformed := matrix.MulVec4(m.triangles[i].vertices[j].Vec4()).Vec3()
			m.triangles[i].vertices[j] = transformed
		}
	}
}

// coordMinMax returns the minimum and maximum value for all vertex coordinates.  This is
// useful for normalization.
func (m *Mesh) coordMinMax() (float64, float64) {
	var min float64 = math.MaxInt16
	var max float64 = math.MinInt16

	for _, t := range m.triangles {
		for i := 0; i < 3; i++ {
			// Min
			min = math.Min(t.vertices[i].x, min)
			min = math.Min(t.vertices[i].y, min)
			min = math.Min(t.vertices[i].z, min)
			// Max
			max = math.Max(t.vertices[i].x, max)
			max = math.Max(t.vertices[i].y, max)
			max = math.Max(t.vertices[i].z, max)
		}
	}
	return min, max
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

	// Not using the Y coord since FFT maps already sit on the floor. Adding
	// the Y translation would put the floor at the models 1/2 height point.
	x := -(maxx + minx) / 2.0
	y := -(maxy + miny) / 2.0
	z := -(maxz + minz) / 2.0

	return Vec3{x, y, z}
}
