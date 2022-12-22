package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	width       int
	height      int
	window      *sdl.Window
	renderer    *sdl.Renderer
	fgTexture   *sdl.Texture // Texture for colorbuffer
	bgTexture   *sdl.Texture // Static texture for background
	colorbuffer []Color
}

// NewWindowFullscreen returns a fullscreen window with
// half the resolution of the display. The resolution is
// pulled from the device itself.
func NewWindowFullscreen() *Window {
	return newWindow(0, 0, true)
}

func NewWindow(width, height int) *Window {
	return newWindow(width, height, false)
}

func newWindow(width, height int, fullscreen bool) *Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	var windowFlags uint32 = sdl.WINDOW_SHOWN

	if fullscreen {
		windowFlags = sdl.WINDOW_FULLSCREEN

		mode, err := sdl.GetCurrentDisplayMode(0)
		if err != nil {
			panic(err)
		}
		width = int(mode.W) / 2
		height = int(mode.H) / 2
	}

	window, err := sdl.CreateWindow(
		"Heretic",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(width),
		int32(height),
		windowFlags,
	)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	fgTexture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(width), int32(height))
	if err != nil {
		panic(err)
	}

	// This is required for the FG texture to use transparent instead of opaque and
	// let the bg texture be seen.
	fgTexture.SetBlendMode(sdl.BLENDMODE_BLEND)

	// Background Texture. Static and pre-drawn.
	bgTexture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STATIC, int32(width), int32(height))
	if err != nil {
		panic(err)
	}

	w := Window{
		width:  width,
		height: height,

		window:    window,
		renderer:  renderer,
		fgTexture: fgTexture,
		bgTexture: bgTexture,

		colorbuffer: make([]Color, width*height),
	}
	w.SetDefaultBackground()
	return &w
}

func (w *Window) SetPixel(x, y int, color Color) {
	if x < 0 || x >= w.width || y < 0 || y >= w.height {
		return
	}
	w.colorbuffer[(w.width*y)+x] = color
}

func (w *Window) Clear(color Color) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			w.SetPixel(x, y, color)
		}
	}
}

func (w *Window) Present() {
	w.fgTexture.Update(nil, unsafe.Pointer(&w.colorbuffer[0]), w.width*4)

	w.renderer.Copy(w.bgTexture, nil, nil)
	w.renderer.Copy(w.fgTexture, nil, nil)

	w.renderer.Present()
	w.Clear(Transparent)
}

func (w *Window) SetDefaultBackground() {
	// Update the texture since it possibly wont change unless a diff bg is used.
	bgBuffer := GenerateCheckerboard(w.width, w.height, LightGray, DarkGray)
	w.bgTexture.Update(nil, unsafe.Pointer(&bgBuffer[0]), w.width*4)
}

func (w *Window) Close() {
	w.fgTexture.Destroy()
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}
