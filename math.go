package main

import "math"

func abs(i int) int {
	return int(math.Abs(float64(i)))
}

type Vec2 struct {
	x, y float64
}

type Vec3 struct {
	x, y, z float64
}

func rotate_x(v Vec3, angle float64) Vec3 {
	return Vec3{
		x: v.x,
		y: v.y*math.Cos(angle) - v.z*math.Sin(angle),
		z: v.y*math.Sin(angle) + v.z*math.Cos(angle),
	}
}

func rotate_y(v Vec3, angle float64) Vec3 {
	return Vec3{
		x: v.x*math.Cos(angle) + v.z*math.Sin(angle),
		y: v.y,
		z: v.x*-math.Sin(angle) + v.z*math.Cos(angle),
	}
}

func rotate_z(v Vec3, angle float64) Vec3 {
	return Vec3{
		x: v.x*math.Cos(angle) - v.y*math.Sin(angle),
		y: v.x*math.Sin(angle) + v.y*math.Cos(angle),
		z: v.z,
	}
}
