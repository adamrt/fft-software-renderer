package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Engine struct {
	isRunning bool
	window    *Window
	renderer  *Renderer
}

func NewEngine(window *Window, renderer *Renderer) *Engine {
	return &Engine{window: window, renderer: renderer}
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
	e.renderer.DrawRect(100, 200, 300, 450, Red)

	e.window.Present()
}
