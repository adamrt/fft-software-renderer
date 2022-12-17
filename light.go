package main

type DirectionalLight struct {
	direction Vec3
	position  Vec3
	color     Color
}

func NewDirectionLight(direction Vec3) DirectionalLight {
	return DirectionalLight{direction: direction}
}

type AmbientLight struct {
	color Color
}
