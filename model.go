package main

type Model struct {
	mesh              Mesh
	trianglesToRender []Triangle
	worldMatrix       Matrix
}

func NewModel(mesh Mesh) Model {
	m := Model{mesh: mesh}
	m.UpdateMatrix()
	return m
}

func (m *Model) UpdateMatrix() {
	m.worldMatrix = MatrixWorld(m.mesh.scale, m.mesh.rotation, m.mesh.translation)
}

func (m *Model) Matrix() Matrix {
	return m.worldMatrix
}
