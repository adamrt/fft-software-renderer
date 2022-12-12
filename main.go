package main

const (
	windowWidth  = 1280
	windowHeight = 720
)

func main() {

	window := NewWindow(windowWidth, windowHeight, false)
	defer window.Close()

	engine := NewEngine(window)

	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}
