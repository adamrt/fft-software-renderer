package main

const (
	windowWidth  = 1024
	windowHeight = 768
)

func main() {

	window := NewWindow(windowWidth, windowHeight)
	defer window.Close()
	renderer := NewRenderer(window)
	reader := NewReader("/home/adam/tmp/fft.iso")

	engine := NewEngine(window, renderer, reader)
	engine.setMap(49)

	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}
