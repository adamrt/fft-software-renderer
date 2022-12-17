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
	ax, ay int, au, av float64,
	bx, by int, bu, bv float64,
	cx, cy int, cu, cv float64,
	texture Texture, palette Palette,
) {
	// We need to sort the vertices by y-coordinate ascending (y0 < y1 < y2)
	if ay > by {
		ay, by = by, ay
		ax, bx = bx, ax
		au, bu = bu, au
		av, bv = bv, av
	}
	if by > cy {
		by, cy = cy, by
		bx, cx = cx, bx
		bu, cu = cu, bu
		bv, cv = cv, bv
	}
	if ay > by {
		ay, by = by, ay
		ax, bx = bx, ax
		au, bu = bu, au
		av, bv = bv, av
	}

	// Create vector points and texture coords after we sort the vertices
	point_a := Vec2{float64(ax), float64(ay)}
	point_b := Vec2{float64(bx), float64(by)}
	point_c := Vec2{float64(cx), float64(cy)}

	//
	// Render the upper part of the triangle (flat-bottom)
	//
	inv_slope_1 := 0.0
	inv_slope_2 := 0.0

	if by-ay != 0 {
		inv_slope_1 = float64(bx-ax) / float64(abs(by-ay))
	}
	if cy-ay != 0 {
		inv_slope_2 = float64(cx-ax) / float64(abs(cy-ay))
	}

	if by-ay != 0 {
		for y := ay; y <= by; y++ {
			var x_start int = bx + int(float64(y-by)*inv_slope_1)
			var x_end int = ax + int(float64(y-ay)*inv_slope_2)

			if x_end < x_start {
				x_start, x_end = x_end, x_start // swap if x_start is to the right of x_end
			}

			for x := x_start; x < x_end; x++ {
				// Draw our pixel with the color that comes from the texture
				r.drawTexel(x, y, texture, palette, point_a, point_b, point_c, au, av, bu, bv, cu, cv)
			}
		}
	}

	//
	// Render the bottom part of the triangle (flat-top)
	//
	inv_slope_1 = 0.0
	inv_slope_2 = 0.0

	if cy-by != 0 {
		inv_slope_1 = float64(cx-bx) / float64(abs(cy-by))
	}
	if cy-ay != 0 {
		inv_slope_2 = float64(cx-ax) / float64(abs(cy-ay))
	}

	if cy-by != 0 {
		for y := by; y <= cy; y++ {
			var x_start int = bx + int(float64(y-by)*inv_slope_1)
			var x_end int = ax + int(float64(y-ay)*inv_slope_2)

			if x_end < x_start {
				x_start, x_end = x_end, x_start // swap if x_start is to the right of x_end
			}

			for x := x_start; x < x_end; x++ {
				// Draw our pixel with the color that comes from the texture
				r.drawTexel(x, y, texture, palette, point_a, point_b, point_c, au, av, bu, bv, cu, cv)
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
	point_a, point_b, point_c Vec2,
	au, av, bu, bv, cu, cv float64,
) {
	point_p := Vec2{float64(x), float64(y)}
	weights := barycentricWeights(point_a, point_b, point_c, point_p)

	alpha := weights.x
	beta := weights.y
	gamma := weights.z

	// Perform the interpolation of all U and V values using barycentric weights
	interpolated_u := (au)*alpha + (bu)*beta + (cu)*gamma
	interpolated_v := (av)*alpha + (bv)*beta + (cv)*gamma

	// Map the UV coordinate to the full texture width and height
	tex_x := abs(int(interpolated_u * float64(texture.width)))
	tex_y := abs(int(interpolated_v * float64(texture.height)))

	idx := (texture.width * tex_y) + tex_x
	if idx >= 0 && idx < textureLen {
		textureColor := texture.data[(tex_y*texture.width)+tex_x]
		// If there is a palette, the current color components will
		// represent the index into the palette.
		if palette != nil {
			textureColor = palette[textureColor.R]
		}

		// Transparent texture
		if textureColor.A == 0 && textureColor.R == 0 && textureColor.G == 0 && textureColor.B == 0 {
			return
		}

		r.window.SetPixel(x, y, textureColor)
	}
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
