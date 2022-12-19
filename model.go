package main

type Model struct {
	mesh              Mesh
	trianglesToRender []Triangle
	worldMatrix       Matrix
}

func NewModel(mesh Mesh) Model {
	m := Model{mesh: mesh}
	m.UpdateWorldMatrix()
	return m
}

func (m *Model) UpdateWorldMatrix() {
	m.worldMatrix = MatrixWorld(m.mesh.scale, m.mesh.rotation, m.mesh.translation)
}

func (m *Model) WorldMatrix() Matrix {
	return m.worldMatrix
}
