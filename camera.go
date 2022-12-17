package main

type Camera struct {
	eye   Vec3
	front Vec3
	up    Vec3
}

func NewCamera(eye, front, up Vec3) *Camera {
	return &Camera{eye, front, up}
}
