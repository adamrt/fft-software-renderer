package main

import (
	"math"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS        = 60
	MSPerFrame = (1000 / FPS)
)

// Data
var (
	leftButtonDown bool = false

	autorotate  = false
	perspective = false

	trianglesToRender = []Triangle{}
	mesh              = Mesh{}

	// Timing
	previous uint32
	delta    float64

	projMatrix Matrix
	camera     = NewCamera(Vec3{-1, 1, -1}, Vec3{0, 0, 0}, Vec3{0, 1, 0})

	light DirectionalLight

	currentMap int
)

type Engine struct {
	window    *Window
	renderer  *Renderer
	reader    *Reader
	isRunning bool
}

func NewEngine(window *Window, renderer *Renderer, reader *Reader) *Engine {
	return &Engine{window: window, renderer: renderer, reader: reader}
}

func (e *Engine) setup() {
	e.isRunning = true
	previous = sdl.GetTicks()

	aspect := float64(e.window.height) / float64(e.window.width)
	fov := math.Pi / 3.0 // (180/3 = 60 degrees). Value is in radians.
	projMatrix = MatrixOrtho(-3, 3, -3, 3, 1, 100)
	if perspective {
		projMatrix = MatrixPerspective(fov, aspect, 1.0, 100.0)
	}

	light = NewDirectionLight(Vec3{0, 0, 1})
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
			case sdl.K_SPACE:
				autorotate = !autorotate
			case sdl.K_p:
				e.changePerspective()
			case sdl.K_j:
				e.prevMap()
			case sdl.K_k:
				e.nextMap()
			}
		case *sdl.MouseButtonEvent:
			if t.Button == sdl.BUTTON_LEFT {
				leftButtonDown = t.Type == sdl.MOUSEBUTTONDOWN
			}
		case *sdl.MouseMotionEvent:
			if leftButtonDown {
				camera.ProcessMouseMovement(float64(t.XRel), float64(t.YRel), delta)
			}
		}
	}
}

func (e *Engine) update() {
	// Variable timestep
	if wait := MSPerFrame - (sdl.GetTicks() - previous); wait > 0 && wait <= MSPerFrame {
		sdl.Delay(wait)
	}
	delta = float64(sdl.GetTicks()-previous) / 1000.0
	previous = sdl.GetTicks()

	if autorotate {
		mesh.rotation.y += 0.5 * delta
	}

	worldMatrix := MatrixWorld(mesh.scale, mesh.rotation, mesh.translation)
	viewMatrix := LookAt(camera.eye, camera.front, camera.up)

	for _, triangle := range mesh.triangles {
		var vertices [3]Vec3

		// Transform vertices with World Matrix
		for i, vertex := range triangle.vertices {
			vertex = worldMatrix.MulVec3(vertex)
			vertex = viewMatrix.MulVec3(vertex)
			vertices[i] = vertex
		}

		// Calculate light before projection
		a, b, c := vertices[0], vertices[1], vertices[2]
		ab, ac := b.Sub(a), c.Sub(a)
		normal := ab.Cross(ac).Normalize()
		lightIntensity := -normal.Dot(light.direction)

		triangle.color = triangle.color.Mul(lightIntensity)
		triangle.avgDepth = (a.z + b.z + c.z) / 3.0

		for i, vertex := range vertices {
			// Projection
			vertex := projMatrix.MulVec4(vertex.Vec4())

			// Perspective divide is using perspective projection.
			if perspective && vertex.w != 0 {
				vertex.x /= vertex.w
				vertex.y /= vertex.w
				vertex.z /= vertex.w
			}

			// Invert the Y asis to compensate for the Y axis of the model and
			// the color buffer being different (+Y up vs +Y down, respectively).
			vertex.y *= -1

			// Scale to the viewport
			vertex.x *= float64(e.window.width / 2)
			vertex.y *= float64(e.window.height / 2)

			// Translate to center of screen
			vertex.x += float64(e.window.width / 2)
			vertex.y += float64(e.window.height / 2)

			vertices[i] = vertex.Vec3()
		}

		if shouldCull(vertices) {
			continue
		}

		for i, vertex := range vertices {
			triangle.points[i] = Vec2{vertex.x, vertex.y}
		}

		trianglesToRender = append(trianglesToRender, triangle)
	}

	// Painters algorithm. Sort the projected triangles so the ones further away are
	// rendered first. This is based on the average of a triangles vertices so there
	// are visual issues. A depth buffer will solve this issue.
	sort.Slice(trianglesToRender, func(i, j int) bool {
		return trianglesToRender[i].avgDepth > trianglesToRender[j].avgDepth
	})
}

func (e *Engine) render() {
	// Draw
	for _, t := range trianglesToRender {
		// Draw triangles
		a, b, c := t.points[0], t.points[1], t.points[2]
		at, bt, ct := t.texcoords[0], t.texcoords[1], t.texcoords[2]
		// e.renderer.DrawFilledTriangle(int(a.x), int(a.y), int(b.x), int(b.y), int(c.x), int(c.y), t.color)

		e.renderer.DrawTexturedTriangle(
			int(a.x), int(a.y), at.U, at.V,
			int(b.x), int(b.y), bt.U, bt.V,
			int(c.x), int(c.y), ct.U, ct.V,
			mesh.texture, t.palette)
		// e.renderer.DrawTriangle(int(a.X), int(a.Y), int(b.X), int(b.Y), int(c.X), int(c.Y), White)

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

func (e *Engine) loadObj(file string) { mesh = NewMeshFromObj(file) }
func (e *Engine) setMap(n int) {
	currentMap = n
	mesh = e.reader.ReadMesh(n)
	e.setup()
}

func (e *Engine) prevMap() {
	if currentMap > 1 {
		e.setMap(currentMap - 1)
	}
}

func (e *Engine) nextMap() {
	if currentMap < 125 {
		e.setMap(currentMap + 1)
	}
}

func (e *Engine) changePerspective() {
	perspective = !perspective
	if perspective {
		mesh.scale = Vec3{1, 1, 1}
	} else {
		mesh.scale = Vec3{3, 3, 3}
	}
	e.setup()

}

// Backface culling
//
// This works by checking the winding order of the projected triangle.
// If its CCW then you can ignore this triangle since it would be back-facing.
//
// NOTE: This method must be done after projection vertices.
func shouldCull(vertices [3]Vec3) bool {
	a, b, c := vertices[0], vertices[1], vertices[2]
	ab, ac := b.Sub(a), c.Sub(a)

	if sign := ab.x*ac.y - ac.x*ab.y; sign < 0 {
		return true
	}

	return false
}
