package main

import (
	"math"
)

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

// DrawLine draws a solid line using the DDA algorithm.
func (r *Renderer) DrawLine(x0, y0, x1, y1 int, color Color) {
	deltaX := x1 - x0
	deltaY := y1 - y0

	var longestSideLength int
	if abs(deltaX) >= abs(deltaY) {
		longestSideLength = abs(deltaX)
	} else {
		longestSideLength = abs(deltaY)
	}

	incX := float64(deltaX) / float64(longestSideLength)
	incY := float64(deltaY) / float64(longestSideLength)

	currentX := float64(x0)
	currentY := float64(y0)

	for i := 0; i <= longestSideLength; i++ {
		r.window.SetPixel(int(math.Round(currentX)), int(math.Round(currentY)), color)
		currentX += incX
		currentY += incY
	}
}

func (r *Renderer) DrawTriangle(ax, ay, bx, by, cx, cy int, color Color) {
	r.DrawLine(int(ax), int(ay), int(bx), int(by), color)
	r.DrawLine(int(bx), int(by), int(cx), int(cy), color)
	r.DrawLine(int(cx), int(cy), int(ax), int(ay), color)

}

// Draw a filled triangle with the flat-top/flat-bottom method
func (r *Renderer) DrawFilledTriangle(ax, ay, bx, by, cx, cy int, color Color) {
	// We need to sort the vertices by y-coordinate ascending (y0 < y1 < y2)
	if ay > by {
		ay, by = by, ay
		ax, bx = bx, ax
	}
	if by > cy {
		by, cy = cy, by
		bx, cx = cx, bx
	}
	if ay > by {
		ay, by = by, ay
		ax, bx = bx, ax
	}

	if by == cy {
		// Draw flat-bottom triangle
		r.fillFlatBottomTriangle(ax, ay, bx, by, cx, cy, color)
	} else if ay == by {
		// Draw flat-top triangle
		r.fillFlatTopTriangle(ax, ay, bx, by, cx, cy, color)
	} else {
		// Calculate the new vertex (Mx,My) using triangle similarity
		My := by
		Mx := (((cx - ax) * (by - ay)) / (cy - ay)) + ax

		// Draw flat-bottom triangle
		r.fillFlatBottomTriangle(ax, ay, bx, by, Mx, My, color)

		// Draw flat-top triangle
		r.fillFlatTopTriangle(bx, by, Mx, My, cx, cy, color)
	}
}

// Draw a filled a triangle with a flat bottom
func (r *Renderer) fillFlatBottomTriangle(ax, ay, bx, by, cx, cy int, color Color) {
	// Find the two slopes (two triangle legs)
	invSlope1 := float64(bx-ax) / float64(by-ay)
	invSlope2 := float64(cx-ax) / float64(cy-ay)

	// Start xStart and xEnd from the top vertex (x0,y0)
	xStart := float64(ax)
	xEnd := float64(ax)

	// Loop all the scanlines from top to bottom
	for y := ay; y <= cy; y++ {
		r.DrawLine(int(xStart), y, int(xEnd), y, color)
		xStart += invSlope1
		xEnd += invSlope2
	}
}

// Draw a filled a triangle with a flat top
func (r *Renderer) fillFlatTopTriangle(ax, ay, bx, by, cx, cy int, color Color) {
	// Find the two slopes (two triangle legs)
	invSlope1 := float64(cx-ax) / float64(cy-ay)
	invSlope2 := float64(cx-bx) / float64(cy-by)

	// Start xStart and xEnd from the bottom vertex (x2,y2)
	xStart := float64(cx)
	xEnd := float64(cx)

	// Loop all the scanlines from bottom to top
	for y := cy; y >= ay; y-- {
		r.DrawLine(int(xStart), y, int(xEnd), y, color)
		xStart -= invSlope1
		xEnd -= invSlope2
	}
}
