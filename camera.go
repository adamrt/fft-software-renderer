package main

import (
	"math"
)

const (
	MinZoom = 2.0
	MaxZoom = 0.5
)

type Projection int

const (
	Orthographic Projection = iota
	Perspective
)

type Camera struct {
	eye   Vec3
	front Vec3
	up    Vec3

	projection       Projection
	projectionMatrix Matrix
	viewMatrix       Matrix

	width  int
	height int
	zoom   float64
}

func NewCamera(eye, front, up Vec3, width, height int) *Camera {
	c := Camera{
		eye:        eye,
		front:      front,
		up:         up,
		width:      width,
		height:     height,
		projection: Orthographic,
		zoom:       1.0,
	}
	c.updateProjectionMatrix()
	c.updateViewMatrix()
	return &c
}

func (c *Camera) aspectRatio() float64     { return float64(c.height) / float64(c.width) }
func (c *Camera) ViewMatrix() Matrix       { return c.viewMatrix }
func (c *Camera) ProjectionMatrix() Matrix { return c.projectionMatrix }

func (c *Camera) updateViewMatrix() {
	c.viewMatrix = LookAt(c.eye, c.front, c.up)
}

func (c *Camera) AdjustZoom(f float64) {
	f *= 0.1
	c.zoom -= f

	if c.zoom < MaxZoom {
		c.zoom = MaxZoom
	}
	if c.zoom > MinZoom {
		c.zoom = MinZoom
	}

	c.updateProjectionMatrix()
}

func (c *Camera) toggleProjection() {
	if c.projection == Orthographic {
		c.projection = Perspective
	} else {
		c.projection = Orthographic
	}
	c.updateProjectionMatrix()
}

func (c *Camera) updateProjectionMatrix() {
	aspect := c.aspectRatio()

	if c.projection == Orthographic {
		w := 1.0 * c.zoom
		h := 1.0 * aspect * c.zoom
		c.projectionMatrix = MatrixOrtho(-w, w, -h, h, 1.0, 100.0)
	} else {
		// (180/3 = 60 degrees). Value is in radians.
		fov := (math.Pi / 3.0) * c.zoom
		c.projectionMatrix = MatrixPerspective(fov, aspect, 1.0, 100.0)
	}
}

func (c *Camera) ProcessMouseMovement(xrel, yrel, delta float64) {
	const EPS = 0.0001

	minPolarAngle := 0.0
	maxPolarAngle := math.Pi // 180 degrees as radians
	minAzimuthAngle := math.Inf(-1)
	maxAzimuthAngle := math.Inf(1)

	// Compute direction vector from target to camera
	tcam := c.eye.Sub(c.front)

	// Calculate angles based on current camera position plus deltas
	radius := tcam.Length()
	theta := math.Atan2(tcam.x, tcam.z) + (xrel * delta / 4)
	phi := math.Acos(tcam.y/radius) + (-yrel * delta / 4)

	// Restrict phi and theta to be between desired limits
	phi = clamp(phi, minPolarAngle, maxPolarAngle)
	phi = clamp(phi, EPS, math.Pi-EPS)
	theta = clamp(theta, minAzimuthAngle, maxAzimuthAngle)

	// Calculate new cartesian coordinates
	tcam.x = radius * math.Sin(phi) * math.Sin(theta)
	tcam.y = radius * math.Cos(phi)
	tcam.z = radius * math.Sin(phi) * math.Cos(theta)

	// Don't allow camera to go below bottom of map
	tcam.y = math.Max(tcam.y, 0.0)

	// Update camera position and orientation
	c.eye = c.front.Add(tcam)

	c.updateViewMatrix()
}
