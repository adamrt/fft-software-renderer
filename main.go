package main

const (
	windowWidth  = 1280
	windowHeight = 720
)

func main() {

	window := NewWindow(windowWidth, windowHeight)
	defer window.Close()
	renderer := NewRenderer(window)
	engine := NewEngine(window, renderer)

	engine.loadObj("cube.obj")
	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}
