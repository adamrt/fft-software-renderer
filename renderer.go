package main

type Renderer struct {
	window *Window
}

func NewRenderer(window *Window) *Renderer {
	return &Renderer{window}
}

func (r *Renderer) DrawRect(x, y, w, h int, color Color) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			r.window.SetPixel(i, j, color)
		}
	}
}
