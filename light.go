package main

type DirectionalLight struct {
	Direction Vec3
	Position  Vec3
	Color     Color
}

func NewDirectionLight(direction Vec3) DirectionalLight {
	return DirectionalLight{Direction: direction}
}

type AmbientLight struct {
	Color Color
}
