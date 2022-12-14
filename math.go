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

func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{
		x: v.x + u.x,
		y: v.y + u.y,
		z: v.z + u.z,
	}
}

func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{
		x: v.x - u.x,
		y: v.y - u.y,
		z: v.z - u.z,
	}
}
func (v Vec3) Mul(f float64) Vec3 {
	return Vec3{
		x: v.x * f,
		y: v.y * f,
		z: v.z * f,
	}
}

func (v Vec3) Div(f float64) Vec3 {
	return Vec3{
		x: v.x / f,
		y: v.y / f,
		z: v.z / f,
	}
}

func (v Vec3) Dot(u Vec3) float64 {
	return v.x*u.x + v.y*u.y + v.z*u.z
}

func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{
		x: v.y*u.z - v.z*u.y,
		y: v.z*u.x - v.x*u.z,
		z: v.x*u.y - v.y*u.x,
	}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	return Vec3{
		x: v.x / l,
		y: v.y / l,
		z: v.z / l,
	}
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
