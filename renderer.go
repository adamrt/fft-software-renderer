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

func (r *Renderer) DrawTriangle(t Triangle) {
	ax, ay := int(t.points[0].x), int(t.points[0].y)
	bx, by := int(t.points[1].x), int(t.points[1].y)
	cx, cy := int(t.points[2].x), int(t.points[2].y)
	r.DrawLine(ax, ay, bx, by, t.color)
	r.DrawLine(bx, by, cx, cy, t.color)
	r.DrawLine(cx, cy, ax, ay, t.color)

}

func (r *Renderer) DrawFilledTriangle(t Triangle) {
	ax, ay := int(t.points[0].x), int(t.points[0].y)
	bx, by := int(t.points[1].x), int(t.points[1].y)
	cx, cy := int(t.points[2].x), int(t.points[2].y)

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
		r.fillFlatBottomTriangle(ax, ay, bx, by, cx, cy, t)
	} else if ay == by {
		// Draw flat-top triangle
		r.fillFlatTopTriangle(ax, ay, bx, by, cx, cy, t)
	} else {
		// Calculate the new vertex (Mx,My) using triangle similarity
		My := by
		Mx := (((cx - ax) * (by - ay)) / (cy - ay)) + ax

		// Draw flat-bottom triangle
		r.fillFlatBottomTriangle(ax, ay, bx, by, Mx, My, t)

		// Draw flat-top triangle
		r.fillFlatTopTriangle(bx, by, Mx, My, cx, cy, t)
	}
}

func (r *Renderer) DrawTexturedTriangle(t Triangle, texture Texture) {
	ax, ay := int(t.points[0].x), int(t.points[0].y)
	bx, by := int(t.points[1].x), int(t.points[1].y)
	cx, cy := int(t.points[2].x), int(t.points[2].y)
	at, bt, ct := t.texcoords[0], t.texcoords[1], t.texcoords[2]

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
				r.drawTexel(x, y, texture, t, a, b, c, at, bt, ct)
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
				r.drawTexel(x, y, texture, t, a, b, c, at, bt, ct)
			}
		}
	}
}

func (r *Renderer) fillFlatBottomTriangle(ax, ay, bx, by, cx, cy int, t Triangle) {
	// Find the two slopes (two triangle legs)
	invSlope1 := float64(bx-ax) / float64(by-ay)
	invSlope2 := float64(cx-ax) / float64(cy-ay)

	// Start xStart and xEnd from the top vertex (x0,y0)
	xStart := float64(ax)
	xEnd := float64(ax)

	// Loop all the scanlines from top to bottom
	color := t.color.Mul(t.lightColor)
	for y := ay; y <= cy; y++ {
		r.DrawLine(int(xStart), y, int(xEnd), y, color)
		xStart += invSlope1
		xEnd += invSlope2
	}
}

func (r *Renderer) fillFlatTopTriangle(ax, ay, bx, by, cx, cy int, t Triangle) {
	// Find the two slopes (two triangle legs)
	invSlope1 := float64(cx-ax) / float64(cy-ay)
	invSlope2 := float64(cx-bx) / float64(cy-by)

	// Start xStart and xEnd from the bottom vertex (x2,y2)
	xStart := float64(cx)
	xEnd := float64(cx)

	// Loop all the scanlines from bottom to top
	color := t.color.Mul(t.lightColor)
	for y := cy; y >= ay; y-- {
		r.DrawLine(int(xStart), y, int(xEnd), y, color)
		xStart -= invSlope1
		xEnd -= invSlope2
	}
}

func (r *Renderer) drawTexel(x, y int, texture Texture, t Triangle, a, b, c Vec2, at, bt, ct Tex) {
	p := Vec2{float64(x), float64(y)}
	alpha, beta, gamma := barycentricWeights(a, b, c, p)

	// Perform the interpolation of all U and V values using barycentric weights
	interpolatedU := at.u*alpha + bt.u*beta + ct.u*gamma
	interpolatedV := at.v*alpha + bt.v*beta + ct.v*gamma

	// Map the UV coordinate to the full texture width and height
	texX := int(interpolatedU * float64(texture.width))
	texY := int(interpolatedV * float64(texture.height))

	// Validate the index is inside the texture.
	index := (texY * texture.width) + texX
	if index < 0 || index > texture.width*texture.height {
		return
	}

	textureColor := texture.data[index]

	// If there is a palette, the current color components will
	// represent the index into the palette.
	if t.palette != nil {
		textureColor = t.palette[textureColor.R]

		if showLighting {
			textureColor = textureColor.Mul(t.lightColor)
		}
		// Transparent texture
		if textureColor.isTrans() {
			return
		}
	}

	r.window.SetPixel(x, y, textureColor)

}

func barycentricWeights(a, b, c, p Vec2) (float64, float64, float64) {
	ab, ac, ap := b.Sub(a), c.Sub(a), p.Sub(a)
	den := ab.x*ac.y - ac.x*ab.y
	v := (ap.x*ac.y - ac.x*ap.y) / den
	w := (ab.x*ap.y - ap.x*ab.y) / den
	u := 1.0 - v - w
	return u, v, w
}

var origin = Vec3{0, 0, 0}
var xAxis = Vec3{1, 0, 0}
var yAxis = Vec3{0, 1, 0}
var zAxis = Vec3{0, 0, 1}

func (r *Renderer) DrawOriginAxis(camera *Camera) {
	worldMatrix := MatrixWorld(Vec3{1, 1, 1}, Vec3{0, 0, 0}, Vec3{0, 0, 0})
	viewMatrix := camera.ViewMatrix()

	vertices := []Vec3{origin, xAxis, yAxis, zAxis}

	for i, vertex := range vertices {
		vertex = worldMatrix.MulVec3(vertex.Mul(1.0))
		vertex = viewMatrix.MulVec3(vertex)
		vertices[i] = vertex
	}
	for i, vertex := range vertices {
		// Projection
		vertex = camera.ProjectionMatrix().MulVec3(vertex)

		// Invert the Y asis to compensate for the Y axis of the model and
		// the color buffer being different (+Y up vs +Y down, respectively).
		vertex.y *= -1

		// Scale to the viewport
		vertex.x *= float64(r.window.width / 2)
		vertex.y *= float64(r.window.height / 2)

		// Translate to center of screen
		vertex.x += float64(r.window.width / 2)
		vertex.y += float64(r.window.height / 2)

		vertices[i] = vertex
	}
	origin, x, y, z := vertices[0], vertices[1], vertices[2], vertices[3]
	r.DrawLine(int(origin.x), int(origin.y), int(x.x), int(x.y), Magenta)
	r.DrawLine(int(origin.x), int(origin.y), int(y.x), int(y.y), Green)
	r.DrawLine(int(origin.x), int(origin.y), int(z.x), int(z.y), Red)
}
