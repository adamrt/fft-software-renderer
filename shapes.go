package main

func cubeMesh() Mesh {
	return Mesh{
		vertices: cubeVertices,
		faces:    cubeFaces,
		rotation: Vec3{0, 0, 0},
	}
}

var cubeVertices = []Vec3{
	{x: -1, y: -1, z: -1}, // 1
	{x: -1, y: 1, z: -1},  // 2
	{x: 1, y: 1, z: -1},   // 3
	{x: 1, y: -1, z: -1},  // 4
	{x: 1, y: 1, z: 1},    // 5
	{x: 1, y: -1, z: 1},   // 6
	{x: -1, y: 1, z: 1},   // 7
	{x: -1, y: -1, z: 1},  // 8
}

var cubeFaces = []Face{
	// front
	{[3]int{1, 2, 3}},
	{[3]int{1, 3, 4}},
	// right
	{[3]int{4, 3, 5}},
	{[3]int{4, 5, 6}},
	// back
	{[3]int{6, 5, 7}},
	{[3]int{6, 7, 8}},
	// left
	{[3]int{8, 7, 2}},
	{[3]int{8, 2, 1}},
	// top
	{[3]int{2, 7, 5}},
	{[3]int{2, 5, 3}},
	// bottom
	{[3]int{6, 8, 1}},
	{[3]int{6, 1, 4}},
}

func cubeDots() []Vec3 {
	var x, y, z float64
	var cube []Vec3
	for x = -1; x <= 1; x += 0.25 {
		for y = -1; y <= 1; y += 0.25 {
			for z = -1; z <= 1; z += 0.25 {
				cube = append(cube, Vec3{x, y, z})
			}
		}
	}
	return cube
}
