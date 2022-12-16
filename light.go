package main

type Light struct {
	direction Vec3
}

func NewLight(direction Vec3) Light {
	return Light{direction: direction}
}
