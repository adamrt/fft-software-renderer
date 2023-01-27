package main

import (
	"fmt"
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
	var path string
	if len(os.Args) > 1 {
		// User specified path
		path = os.Args[1]
	} else {
		// Use default path
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(usr.HomeDir, "tmp", "emu", "fft.bin")
	}

	if _, err := os.Stat(path); err != nil {
		fmt.Printf("Usage: go run *.go <path-to-bin>\n\n")
		os.Exit(0)
	}

	return path
}
