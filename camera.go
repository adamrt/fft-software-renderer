package main

import "math"

type Camera struct {
	eye   Vec3
	front Vec3
	up    Vec3
}

func NewCamera(eye, front, up Vec3) *Camera {
	return &Camera{eye, front, up}
}

func (c *Camera) ProcessMouseMovement(xrel, yrel, delta float64) {
	const EPS = 0.0001

	minPolarAngle := 0.0
	maxPolarAngle := math.Pi // 180 degrees as radians
	minAzimuthAngle := math.Inf(-1)
	maxAzimuthAngle := math.Inf(1)

	// Compute direction vector from target to camera
	tcam := c.eye
	tcam = tcam.Sub(c.front)

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

	// Update camera position and orientation
	c.eye = c.front.Add(tcam)
}
