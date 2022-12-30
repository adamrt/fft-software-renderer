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

func (r *Renderer) DrawTexturedTriangle(
	ax, ay int, at Tex,
	bx, by int, bt Tex,
	cx, cy int, ct Tex,
	texture Texture, palette Palette,
) {
	// We need to sort the vertices by y-coordinate ascending (y0 < y1 < y2)
	if ay > by {
		ay, by = by, ay
		ax, bx = bx, ax
		at, bt = bt, at
	}
	if by > cy {
		by, cy = cy, by
		bx, cx = cx, bx
		bt, ct = ct, bt
	}
	if ay > by {
		ay, by = by, ay
		ax, bx = bx, ax
		at, bt = bt, at
	}

	// Create vector points and texture coords after we sort the vertices
	a := Vec2{float64(ax), float64(ay)}
	b := Vec2{float64(bx), float64(by)}
	c := Vec2{float64(cx), float64(cy)}

	//
	// Render the upper part of the triangle (flat-bottom)
	//
	invSlope1 := 0.0
	invSlope2 := 0.0

	if by-ay != 0 {
		invSlope1 = float64(bx-ax) / float64(abs(by-ay))
	}
	if cy-ay != 0 {
		invSlope2 = float64(cx-ax) / float64(abs(cy-ay))
	}

	if by-ay != 0 {
		for y := ay; y <= by; y++ {
			var xStart int = int(float64(bx) + (float64(y-by) * invSlope1))
			var xEnd int = int(float64(ax) + (float64(y-ay) * invSlope2))

			if xEnd < xStart {
				xStart, xEnd = xEnd, xStart // swap if xStart is to the right of xEnd
			}

			for x := xStart; x < xEnd; x++ {
				// Draw our pixel with the color that comes from the texture
				r.drawTexel(x, y, texture, palette, a, b, c, at, bt, ct)
			}
		}
	}

	//
	// Render the bottom part of the triangle (flat-top)
	//
	invSlope1 = 0.0
	invSlope2 = 0.0

	if cy-by != 0 {
		invSlope1 = float64(cx-bx) / float64(abs(cy-by))
	}
	if cy-ay != 0 {
		invSlope2 = float64(cx-ax) / float64(abs(cy-ay))
	}

	if cy-by != 0 {
		for y := by; y <= cy; y++ {
			var xStart int = int(float64(bx) + (float64(y-by) * invSlope1))
			var xEnd int = int(float64(ax) + (float64(y-ay) * invSlope2))

			if xEnd < xStart {
				xStart, xEnd = xEnd, xStart // swap if xStart is to the right of xEnd
			}

			for x := xStart; x < xEnd; x++ {
				// Draw our pixel with the color that comes from the texture
				r.drawTexel(x, y, texture, palette, a, b, c, at, bt, ct)
			}
		}
	}
}

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

func (r *Renderer) drawTexel(
	x, y int,
	texture Texture, palette Palette,
	a, b, c Vec2,
	at, bt, ct Tex,
) {
	p := Vec2{float64(x), float64(y)}
	weights := barycentricWeights(a, b, c, p)

	alpha := weights.x
	beta := weights.y
	gamma := weights.z

	// Perform the interpolation of all U and V values using barycentric weights
	interpolatedU := (at.u)*alpha + (bt.u)*beta + (ct.u)*gamma
	interpolatedV := (at.v)*alpha + (bt.v)*beta + (ct.v)*gamma

	// Map the UV coordinate to the full texture width and height
	texX := abs(int(interpolatedU * float64(texture.width)))
	texY := abs(int(interpolatedV * float64(texture.height)))

	// Validate the index is inside the texture.
	index := (texY * texture.width) + texX
	if index < 0 || index > texture.width*texture.height {
		return
	}

	textureColor := texture.data[index]
	// If there is a palette, the current color components will
	// represent the index into the palette.
	if palette != nil {
		textureColor = palette[textureColor.R]
	}

	// Transparent texture
	if textureColor.isTrans() {
		return
	}

	r.window.SetPixel(x, y, textureColor)

}

func barycentricWeights(a, b, c, p Vec2) Vec3 {
	ab := b.Sub(a)
	bc := c.Sub(b)
	ac := c.Sub(a)
	ap := p.Sub(a)
	bp := p.Sub(b)

	// Calcualte the area of the full triangle ABC using cross product (area of parallelogram)
	triangleArea := (ab.x*ac.y - ab.y*ac.x)

	// Weight alpha is the area of subtriangle BCP divided by the area of the full triangle ABC
	alpha := (bc.x*bp.y - bp.x*bc.y) / triangleArea

	// Weight beta is the area of subtriangle ACP divided by the area of the full triangle ABC
	beta := (ap.x*ac.y - ac.x*ap.y) / triangleArea

	// Weight gamma is easily found since barycentric cooordinates always add up to 1
	gamma := 1 - alpha - beta

	weights := Vec3{alpha, beta, gamma}
	return weights
}
