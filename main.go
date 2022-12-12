package main

func main() {

	window := NewWindow(800, 800)
	defer window.Close()

	engine := NewEngine(window)

	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}
