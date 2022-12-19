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
	currentMap     int

	autorotate    = false
	perspective   = false
	showTexture   = true
	showWireframe = false

	// Timing
	previous uint32
	delta    float64

	model      = Model{}
	projMatrix Matrix
	camera     = NewCamera(Vec3{-1, 1, -1}, Vec3{0, 0, 0}, Vec3{0, 1, 0})

	light DirectionalLight
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
			case sdl.K_t:
				showTexture = !showTexture
			case sdl.K_w:
				showWireframe = !showWireframe
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
		model.mesh.rotation.y += 0.5 * delta
	}

	viewMatrix := LookAt(camera.eye, camera.front, camera.up)
	e.updateModel(&model, viewMatrix)
}

func (e *Engine) updateModel(model *Model, viewMatrix Matrix) {
	model.UpdateWorldMatrix()

	for _, triangle := range model.mesh.triangles {
		var vertices [3]Vec3

		// Transform vertices with World Matrix
		for i, vertex := range triangle.vertices {
			vertex = model.WorldMatrix().MulVec3(vertex)
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

		model.trianglesToRender = append(model.trianglesToRender, triangle)
	}

	// Painters algorithm. Sort the projected triangles so the ones further away are
	// rendered first. This is based on the average of a triangles vertices so there
	// are visual issues. A depth buffer will solve this issue.
	sort.Slice(model.trianglesToRender, func(i, j int) bool {
		return model.trianglesToRender[i].avgDepth > model.trianglesToRender[j].avgDepth
	})
}

func (e *Engine) render() {
	// Draw
	for _, t := range model.trianglesToRender {
		// Draw triangles
		a, b, c := t.points[0], t.points[1], t.points[2]

		if showTexture {
			at, bt, ct := t.texcoords[0], t.texcoords[1], t.texcoords[2]
			e.renderer.DrawTexturedTriangle(int(a.x), int(a.y), at.U, at.V, int(b.x), int(b.y), bt.U, bt.V, int(c.x), int(c.y), ct.U, ct.V, model.mesh.texture, t.palette)
		} else {
			e.renderer.DrawFilledTriangle(int(a.x), int(a.y), int(b.x), int(b.y), int(c.x), int(c.y), t.color)
		}

		if showWireframe {
			e.renderer.DrawTriangle(int(a.x), int(a.y), int(b.x), int(b.y), int(c.x), int(c.y), Magenta)
		}
	}

	// Present
	e.window.Present()

	// Clear triangles from last frame
	model.trianglesToRender = model.trianglesToRender[:0]
}

func (e *Engine) renderModel(model *Model, viewMatrix Matrix) {

}

func (e *Engine) loadObj(file string) { model.mesh = NewMeshFromObj(file) }
func (e *Engine) setMap(n int) {
	currentMap = n
	model.mesh = e.reader.ReadMesh(n)
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
		model.mesh.scale = Vec3{2, 2, 2}
	} else {
		model.mesh.scale = Vec3{5, 5, 5}
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
