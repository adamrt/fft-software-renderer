package main

import (
	"os"
	"os/user"
	"path/filepath"
)

const (
	windowWidth  = 1024
	windowHeight = 768
)

func main() {

	window := NewWindow(windowWidth, windowHeight)
	defer window.Close()
	renderer := NewRenderer(window)

	path := getPath()
	reader := NewReader(path)

	engine := NewEngine(window, renderer, reader)
	engine.setMap(49)

	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}

func getPath() string {
	// Use user specified path
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	// Use default path
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, "tmp", "emu", "fft.bin")
}
