package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS        = 60
	MSPerFrame = (1000 / FPS)
)

// Data
var (
	trianglesToRender = []Triangle{}
	mesh              = Mesh{}

	// Timing
	previous uint32 = 0

	// For now the camera position is at 0,0,0 until we get a proper camera with a
	// lookat() function and a view matrix.
	cameraPosition = Vec3{0, 0, 0}
	projMatrix     Matrix
)

type Engine struct {
	isRunning bool
	window    *Window
	renderer  *Renderer
}

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{window: window, renderer: renderer}
}

func (e *Engine) loadObj(file string) {
	mesh = NewMeshFromObj(file)
}
func (e *Engine) setup() {
	e.isRunning = true
	previous = sdl.GetTicks()

	aspect := float64(e.window.height) / float64(e.window.width)
	fov := math.Pi / 3.0 // (180/3 = 60 degrees). Value is in radians.
	projMatrix = MatrixPerspective(fov, aspect, 0.1, 100.0)
	projMatrix = MatrixOrtho(-5, 5, -5, 5, 1, 100)
}

func (e *Engine) processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false
		case *sdl.KeyboardEvent:
			if t.Type != sdl.KEYDOWN {
				continue
			}
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				e.isRunning = false
			}
		}
	}
}

func (e *Engine) update() {
	// Variable timestep
	if wait := MSPerFrame - (sdl.GetTicks() - previous); wait > 0 && wait <= MSPerFrame {
		sdl.Delay(wait)
	}
	previous = sdl.GetTicks()

	// Mesh transformation setup
	mesh.rotation.x += 0.005
	mesh.rotation.y += 0.005
	// Temporary until we have a camera/view matrix
	mesh.translation.z = 5.0

	worldMatrix := MatrixWorld(mesh.scale, mesh.rotation, mesh.translation)

	for _, triangle := range mesh.triangles {
		var vertices [3]Vec3
		for i, vertex := range triangle.vertices {
			// Transform
			vertex = worldMatrix.MulVec3(vertex)

			// Project
			vertex := projMatrix.MulVec3(vertex)

			// Invert the Y asis to compensate for the Y axis of the model and
			// the color buffer being different (+Y up vs +Y down, respectively).
			vertex.y *= -1

			// Scale to the viewport
			vertex.x *= float64(e.window.width / 2)
			vertex.y *= float64(e.window.height / 2)

			// Translate to center of screen
			vertex.x += float64(e.window.width / 2)
			vertex.y += float64(e.window.height / 2)

			vertices[i] = vertex
		}

		// Backface culling
		a, b, c := vertices[0], vertices[1], vertices[2]
		ab := b.Sub(a)
		ac := c.Sub(a)
		if sign := ac.x*ab.y - ab.x*ac.y; sign < 0 {
			continue
		}

		for i, vertex := range vertices {
			triangle.points[i] = Vec2{
				x: vertex.x,
				y: vertex.y,
			}
		}

		trianglesToRender = append(trianglesToRender, triangle)
	}
}

func (e *Engine) render() {
	// Draw
	for _, t := range trianglesToRender {
		// Draw triangles
		a, b, c := t.points[0], t.points[1], t.points[2]
		e.renderer.DrawLine(int(a.x), int(a.y), int(b.x), int(b.y), White)
		e.renderer.DrawLine(int(b.x), int(b.y), int(c.x), int(c.y), White)
		e.renderer.DrawLine(int(c.x), int(c.y), int(a.x), int(a.y), White)

		// Draw vertices
		// e.renderer.DrawRect(int(a.x)-2, int(a.y)-2, 4, 4, Red)
		// e.renderer.DrawRect(int(b.x)-2, int(b.y)-2, 4, 4, Red)
		// e.renderer.DrawRect(int(c.x)-2, int(c.y)-2, 4, 4, Red)
	}

	// Present
	e.window.Present()

	// Clear triangles from last frame
	trianglesToRender = trianglesToRender[:0]
}
