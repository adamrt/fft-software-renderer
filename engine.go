package main

import (
	"fmt"
	"sort"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS        = 60
	MSPerFrame = (1000 / FPS)
)

var (
	// Options
	autorotate        = false
	showHelp          = true
	showTexture       = true
	showWireframe     = false
	showMapBackground = false
	showLighting      = true

	// Timing
	previous   uint32
	delta      float64
	frameCount int

	// Controls
	leftButtonDown = false

	// Data
	model      Model
	currentMap int

	modelScale float64 = 15.0
)

type Engine struct {
	window   *Window
	renderer *Renderer
	camera   *Camera

	reader    *Reader
	isRunning bool
}

func NewEngine(window *Window, renderer *Renderer, reader *Reader) *Engine {
	return &Engine{
		window:   window,
		renderer: renderer,
		reader:   reader,
		camera:   NewCamera(Vec3{1, 1, -1}, Vec3{0, 0, 0}, Vec3{0, 1, 0}, window.width, window.height),
	}
}

func (e *Engine) setup() {
	e.isRunning = true
	previous = sdl.GetTicks()
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
			case sdl.K_a:
				autorotate = !autorotate
			case sdl.K_h:
				showHelp = !showHelp
			case sdl.K_t:
				showTexture = !showTexture
			case sdl.K_w:
				showWireframe = !showWireframe
			case sdl.K_l:
				showLighting = !showLighting
			case sdl.K_p:
				e.camera.toggleProjection()
			case sdl.K_b:
				e.toggleBackgorund()
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
				e.camera.ProcessMouseMovement(float64(t.XRel), float64(t.YRel), delta)
			}
		case *sdl.MouseWheelEvent:
			e.camera.AdjustZoom(float64(t.PreciseY))
		}
	}
}

func (e *Engine) update() {
	frameCount++
	// Variable timestep
	if wait := MSPerFrame - (sdl.GetTicks() - previous); wait > 0 && wait <= MSPerFrame {
		sdl.Delay(wait)
	}

	delta = float64(sdl.GetTicks() - previous)
	if frameCount > 10 {
		e.window.SetTitle(fmt.Sprintf("FPS: %.2f", 1000.0/delta))
		frameCount = 0
	}
	delta = delta / 1000.0

	previous = sdl.GetTicks()

	if autorotate {
		model.mesh.rotation.y += 0.5 * delta
	}

	matrix := MatrixScale(Vec3{modelScale, modelScale, modelScale})
	lights := make([]DirectionalLight, len(model.mesh.directionalLights))
	for i, light := range model.mesh.directionalLights {
		light.position = matrix.MulVec3(light.position)
		lights[i] = light
	}

	model.UpdateMatrix()
	for _, triangle := range model.mesh.triangles {
		var vertices [3]Vec3

		// Transform vertices with World Matrix
		for i, vertex := range triangle.vertices {
			vertex = model.Matrix().MulVec3(vertex)
			vertices[i] = vertex
		}

		normal := verticesNormal(vertices)
		for _, light := range lights {
			intensity := -normal.Dot(vertices[0].Sub(light.position).Normalize())
			color := light.color.Scale(intensity)
			triangle.lightColor = triangle.lightColor.Add(color)
		}
		triangle.lightColor = triangle.lightColor.Add(model.mesh.ambientLight.color).Scale(2.0)

		for i, vertex := range vertices {
			vertex = e.camera.ViewMatrix().MulVec3(vertex)
			vertices[i] = vertex
		}

		// Calculate average depth after vertex is in view space.
		a, b, c := vertices[0], vertices[1], vertices[2]
		triangle.avgDepth = (a.z + b.z + c.z) / 3.0

		// Project vertices with Projection Matrix.
		for i, vertex := range vertices {
			// Projection
			vertex := e.camera.ProjectionMatrix().MulVec4(vertex.Vec4())

			// Perspective divide is using perspective projection.
			if e.camera.projection == Perspective {
				if vertex.w != 0 {
					vertex.x /= vertex.w
					vertex.y /= vertex.w
					vertex.z /= vertex.w
				}
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
		if showTexture {
			e.renderer.DrawTexturedTriangle(t, model.mesh.texture)
		} else {
			e.renderer.DrawFilledTriangle(t)
		}

		if showWireframe {
			t.color = Magenta
			e.renderer.DrawTriangle(t)
		}

		// e.renderer.DrawOriginAxis(e.camera)
	}

	var (
		textHelp       = "[H]elp: Show"
		textProj       = "[P]rojection: "
		textBackground = "[B]ackground: "
		textLighting   = "[L]ighting: "
		textTexture    = "[T]exture: "
		textWireframe  = "[W]ireframe: "
		textAutorotate = "[A]uto rotate: "
	)
	if showHelp {
		if e.camera.projection == Orthographic {
			textProj += "Orthographic"
		} else {
			textProj += "Perspective"
		}
		if showMapBackground {
			textBackground += "Map"
		} else {
			textBackground += "Default"
		}
		if showLighting {
			textLighting += "Enabled"
		} else {
			textLighting += "Disabled"
		}
		if showTexture {
			textTexture += "Show"
		} else {
			textTexture += "Hide"
		}
		if showWireframe {
			textWireframe += "Show"
		} else {
			textWireframe += "Hide"
		}
		if autorotate {
			textAutorotate += "On"
		} else {
			textAutorotate += "Off"
		}
		e.window.TextBackground(200, 250, Color{255, 255, 255, 30})
		e.window.SetText(10, 10, textHelp, White)
		e.window.SetText(10, 40, textProj, White)
		e.window.SetText(10, 70, textTexture, White)
		e.window.SetText(10, 100, textLighting, White)
		e.window.SetText(10, 130, textWireframe, White)
		e.window.SetText(10, 160, textBackground, White)
		e.window.SetText(10, 190, textAutorotate, White)
		e.window.SetText(10, 220, "[J/K] Next/Previous", White)
	}
	// Present
	e.window.Present()

	// Clear triangles from last frame
	model.trianglesToRender = model.trianglesToRender[:0]
}

func (e *Engine) loadObj(file string) {
	model.mesh = NewMeshFromObj(file)
}

func (e *Engine) toggleBackgorund() {
	showMapBackground = !showMapBackground
	if showMapBackground {
		e.updateBackgroundTexture()
	} else {
		e.window.SetDefaultBackground()
	}
}

func (e *Engine) setMap(n int) {
	currentMap = n
	model.mesh = e.reader.ReadMesh(n)

	// Center camera on center of obj
	center := model.mesh.coordCenter().Mul(modelScale)
	e.camera.front = center
	e.camera.updateViewMatrix()

	if showMapBackground {
		e.updateBackgroundTexture()
	}
}

func (e *Engine) updateBackgroundTexture() {
	bg := model.mesh.background
	bgBuffer := make([]Color, e.window.width*e.window.height)
	for y := 0; y < e.window.height; y++ {
		color := bg.At(y, e.window.height)
		for x := 0; x < e.window.width; x++ {
			bgBuffer[(e.window.width*(e.window.height-y-1))+x] = color
		}
	}
	e.window.bgTexture.Update(nil, unsafe.Pointer(&bgBuffer[0]), e.window.width*4)
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

func verticesNormal(vv [3]Vec3) Vec3 {
	a := vv[0]
	b := vv[1]
	c := vv[2]
	vectorAB := b.Sub(a).Normalize()
	vectorAC := c.Sub(a).Normalize()
	normal := vectorAB.Cross(vectorAC).Normalize()
	return normal
}
