package main

func shapeCube() []Vec3 {
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
