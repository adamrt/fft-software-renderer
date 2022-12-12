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
	texture     *sdl.Texture
	colorbuffer []Color
}

func NewWindow(width, height int) *Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Heretic", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	// Our Color struct is in RGBA8888 format but the SDL texture is set to ABGR8888.
	// SDL reads the strict in big Indian and we are currently on little endian.
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(width), int32(height))
	if err != nil {
		panic(err)
	}

	return &Window{
		width:  width,
		height: height,

		window:   window,
		renderer: renderer,
		texture:  texture,

		colorbuffer: make([]Color, width*height),
	}
}

func (w *Window) SetPixel(x, y int, color Color) {
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
	// w.width*4 is the pitch (size of each row). Width * 32bit color.
	w.texture.Update(nil, unsafe.Pointer(&w.colorbuffer[0]), w.width*4)
	w.renderer.Copy(w.texture, nil, nil)
	w.renderer.Present()
}

func (w *Window) Close() {
	w.texture.Destroy()
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}
