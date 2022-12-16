package main

type Mesh struct {
	triangles []Triangle

	scale       Vec3
	rotation    Vec3
	translation Vec3
}

func NewMesh() Mesh {
	return Mesh{scale: Vec3{1, 1, 1}}
}

type Triangle struct {
	vertices [3]Vec3
	points   [3]Vec2
	color    Color
}
