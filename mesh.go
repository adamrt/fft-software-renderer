package main

type Mesh struct {
	triangles []Triangle
	rotation  Vec3
}

type Triangle struct {
	points [3]Vec3
}
