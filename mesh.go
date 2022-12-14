package main

type Mesh struct {
	triangles []Triangle
	rotation  Vec3
}

type Triangle struct {
	vertices [3]Vec3
	points   [3]Vec2
}
