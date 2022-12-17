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

	reader := NewReader("/home/adam/tmp/fft.iso")
	mesh := reader.ReadMesh(12)
	engine.setMesh(mesh)

	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}
