package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Engine struct {
	isRunning bool
	window    *Window
}

func NewEngine(window *Window) *Engine {
	return &Engine{window: window}
}

func (e *Engine) setup() { e.isRunning = true }

func (e *Engine) processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false
		case *sdl.KeyboardEvent:
			if t.Type != sdl.KEYDOWN {
				continue
			}
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				e.isRunning = false
			}
		}
	}
}

func (e *Engine) update() {}

func (e *Engine) render() {
	e.window.SetPixel(100, 200, Red)
	e.window.Present()
	e.window.Clear(Transparent)
}
