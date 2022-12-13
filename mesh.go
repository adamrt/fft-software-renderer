package main

type Mesh struct {
	faces    []Face
	vertices []Vec3
	rotation Vec3
}

type Face struct {
	indexes [3]int
}

// Triangle represents a projected triangle of Vec2 points.
type Triangle struct {
	points [3]Vec2
}
