package main

const (
	windowWidth  = 1024
	windowHeight = 1024
)

func main() {

	window := NewWindow(windowWidth, windowHeight)
	defer window.Close()
	renderer := NewRenderer(window)
	engine := NewEngine(window, renderer)

	engine.loadObj("assets/f22.obj")
	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}
