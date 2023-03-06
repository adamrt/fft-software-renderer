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

	path, err := getPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Usage: go run *.go <path-to-bin>")
		os.Exit(1)
	}

	reader := NewReader(path)
	defer reader.Close()

	engine := NewEngine(window, renderer, reader)
	engine.setMap(49)

	engine.setup()
	for engine.isRunning {
		engine.processInput()
		engine.update()
		engine.render()
	}
}

func getPath() (string, error) {
	var path string
	if len(os.Args) > 1 {
		// User specified path
		path = os.Args[1]
	} else {
		// Use default path
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, "media", "emu", "fft.bin")
	}

	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	return path, nil
}
