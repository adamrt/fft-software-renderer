package main

const (
	windowWidth  = 1024
	windowHeight = 1024
)

func main() {

	window := NewWindow(windowWidth, windowHeight)
	defer window.Close()
	renderer := NewRenderer(window)
	reader := NewReader("/home/adam/tmp/fft.iso")

	engine := NewEngine(window, renderer, reader)
	engine.setMap(54)

	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}
