package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type TextTexture struct {
	texture *sdl.Texture
	rect    *sdl.Rect
}

type Window struct {
	width  int
	height int

	window    *sdl.Window
	renderer  *sdl.Renderer
	fgTexture *sdl.Texture // Texture for colorbuffer
	bgTexture *sdl.Texture // Static texture for background
	font      *ttf.Font

	colorbuffer  []Color
	textTextures []TextTexture // Static texture for background
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

	if err := ttf.Init(); err != nil {
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

	font, err := ttf.OpenFont("assets/arial.ttf", 16)
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
		font:      font,

		colorbuffer: make([]Color, width*height),
	}
	w.SetDefaultBackground()
	return &w
}

func (w *Window) SetText(x, y int, text string, color Color) {
	surface, err := w.font.RenderUTF8Blended(text, color.SDL())
	if err != nil {
		panic(err)
	}
	texture, err := w.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	_, _, width, height, _ := texture.Query()
	rect := &sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height}
	w.textTextures = append(w.textTextures, TextTexture{texture, rect})
}

func (w *Window) TextBackground(width, height int, color Color) {
	texture, err := w.renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(width), int32(height))
	if err != nil {
		panic(err)
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	buffer := make([]Color, width*height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x == 0 || y == 0 || x == width-1 || y == height-1 {
				buffer[y*width+x] = White
			} else {
				buffer[y*width+x] = color
			}
		}
	}

	texture.Update(nil, unsafe.Pointer(&buffer[0]), width*4)
	rect := &sdl.Rect{X: int32(0), Y: int32(0), W: int32(width), H: int32(height)}
	w.textTextures = append(w.textTextures, TextTexture{texture, rect})
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
	for _, tt := range w.textTextures {
		w.renderer.Copy(tt.texture, nil, tt.rect)
	}

	w.renderer.Present()
	w.Clear(Transparent)
	w.textTextures = w.textTextures[:0]
}

func (w *Window) SetDefaultBackground() {
	// Update the texture since it possibly wont change unless a diff bg is used.
	bgBuffer := GenerateCheckerboard(w.width, w.height, LightGray, DarkGray)
	w.bgTexture.Update(nil, unsafe.Pointer(&bgBuffer[0]), w.width*4)
}

func (w *Window) Close() {
	w.fgTexture.Destroy()
	w.bgTexture.Destroy()
	for _, tt := range w.textTextures {
		tt.texture.Destroy()
	}
	w.renderer.Destroy()
	w.window.Destroy()
	w.font.Close()
	ttf.Quit()
	sdl.Quit()
}
