package main

const (
	windowWidth  = 800
	windowHeight = 800
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
