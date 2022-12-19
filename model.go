package main

type Model struct {
	mesh              Mesh
	trianglesToRender []Triangle
}

func NewModel(mesh Mesh) Model {
	return Model{
		mesh: mesh,
	}
}

func (m Model) Matrix() Matrix {
	return MatrixWorld(m.mesh.scale, m.mesh.rotation, m.mesh.translation)
}
