package main

import "math"

func abs(i int) int {
	return int(math.Abs(float64(i)))
}

//
// Vec2
//

type Vec2 struct {
	x, y float64
}

//
// Vec3
//

type Vec3 struct {
	x, y, z float64
}

func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{
		x: v.x + u.x,
		y: v.y + u.y,
		z: v.z + u.z,
	}
}

func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{
		x: v.x - u.x,
		y: v.y - u.y,
		z: v.z - u.z,
	}
}
func (v Vec3) Mul(f float64) Vec3 {
	return Vec3{
		x: v.x * f,
		y: v.y * f,
		z: v.z * f,
	}
}

func (v Vec3) Div(f float64) Vec3 {
	return Vec3{
		x: v.x / f,
		y: v.y / f,
		z: v.z / f,
	}
}

func (v Vec3) Dot(u Vec3) float64 {
	return v.x*u.x + v.y*u.y + v.z*u.z
}

func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{
		x: v.y*u.z - v.z*u.y,
		y: v.z*u.x - v.x*u.z,
		z: v.x*u.y - v.y*u.x,
	}
}

func (v Vec3) Vec4() Vec4 {
	return Vec4{v.x, v.y, v.z, 1.0}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	return Vec3{
		x: v.x / l,
		y: v.y / l,
		z: v.z / l,
	}
}

//
// Vec4
//

type Vec4 struct {
	x, y, z, w float64
}

func (v Vec4) Vec3() Vec3 {
	return Vec3{v.x, v.y, v.z}
}

//
// Matrix 4x4
//

// Matrix is a 4x4 row-major matrix
type Matrix [4][4]float64

// Mul multiplies a matrix a by matrix b.
func (a Matrix) Mul(b Matrix) Matrix {
	var m Matrix
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			m[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j] + a[i][3]*b[3][j]
		}
	}
	return m
}

// MulVec4 multiples a matrix by a Vec4 and returns a Vec4.
func (m Matrix) MulVec3(v Vec3) Vec3 {
	return m.MulVec4(v.Vec4()).Vec3()
}

// MulVec4 multiples a matrix by a Vec4 and returns a Vec4.
func (m Matrix) MulVec4(v Vec4) Vec4 {
	return Vec4{
		m[0][0]*v.x + m[0][1]*v.y + m[0][2]*v.z + m[0][3]*v.w,
		m[1][0]*v.x + m[1][1]*v.y + m[1][2]*v.z + m[1][3]*v.w,
		m[2][0]*v.x + m[2][1]*v.y + m[2][2]*v.z + m[2][3]*v.w,
		m[3][0]*v.x + m[3][1]*v.y + m[3][2]*v.z + m[3][3]*v.w,
	}
}

// Return an Identity Matrix
// | 1  0  0  0 |
// | 0  1  0  0 |
// | 0  0  1  0 |
// | 0  0  0  0 |
func MatrixIdentity() Matrix {
	return Matrix{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Return the combination of Scale, Translation and Rotation matricies.  This is just a
// helper function to clean up construction in the main Engine.update() method.
//
// Always scale, rotate then translate.
func MatrixWorld(scale, rotation, translation Vec3) Matrix {
	matScale := MatrixScale(scale)
	matRX := MatrixRotationX(rotation.x)
	matRY := MatrixRotationY(rotation.y)
	matRZ := MatrixRotationZ(rotation.z)
	matTrans := MatrixTranslation(translation)

	world := MatrixIdentity()
	world = matScale.Mul(world)
	world = matRX.Mul(world)
	world = matRY.Mul(world)
	world = matRZ.Mul(world)
	world = matTrans.Mul(world)
	return world
}

// Return a Scale Matrix
// | sx  0  0  0 |
// |  0 sy  0  0 |
// |  0  0 sx  0 |
// |  0  0  0  1 |
func MatrixScale(v Vec3) Matrix {
	return Matrix{
		{v.x, 0, 0, 0},
		{0, v.y, 0, 0},
		{0, 0, v.z, 0},
		{0, 0, 0, 1},
	}
}

// Return a Translation Matrix
// | 1  0  0  tx |
// | 0  1  0  ty |
// | 0  0  1  tz |
// | 0  0  0   1 |
func MatrixTranslation(v Vec3) Matrix {
	return Matrix{
		{1, 0, 0, v.x},
		{0, 1, 0, v.y},
		{0, 0, 1, v.z},
		{0, 0, 0, 1},
	}
}

// Return a Rotation Matrix for x axis
// | 1  0    0    0 |
// | 0  cos -sin  0 |
// | 0  sin  cos  0 |
// | 0  0    0    1 |
func MatrixRotationX(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)

	return Matrix{
		{1, 0, 0, 0},
		{0, c, -s, 0},
		{0, s, c, 0},
		{0, 0, 0, 1},
	}
}

// Return a Rotation Matrix for y axis
// | cos  0  sin  0 |
// | 0    1  0    0 |
// |-sin  0  cos  0 |
// | 0    0  0    1 |
func MatrixRotationY(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)

	return Matrix{
		{c, 0, s, 0},
		{0, 1, 0, 0},
		{-s, 0, c, 0},
		{0, 0, 0, 1},
	}
}

// Return a Rotation Matrix for z axis
// | cos -sin  0  0 |
// | sin  cos  0  0 |
// | 0    0    1  0 |
// | 0    0    0  1 |
func MatrixRotationZ(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)

	return Matrix{
		{c, -s, s, 0},
		{s, c, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Return a Perspective Projection Matrix
//
// The 3/2==1 stores the original z value for use in MulProjection so we can do
// perspective divide in MulVec4Proj().
func MatrixPerspective(fov, aspect, znear, zfar float64) Matrix {
	m := Matrix{}
	m[0][0] = aspect * (1 / math.Tan(fov/2))
	m[1][1] = 1 / math.Tan(fov/2)
	m[2][2] = zfar / (zfar - znear)
	m[2][3] = (-zfar * znear) / (zfar - znear)
	m[3][2] = 1.0
	return m
}

func MatrixOrtho(left, right, bottom, top, near, far float64) Matrix {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)
	m := MatrixIdentity()
	m[0][0] = 2.0 / rml
	m[1][1] = 2.0 / tmb
	m[2][2] = -2.0 / fmn
	m[3][0] = -(right + left) / rml
	m[3][1] = -(top + bottom) / tmb
	m[3][2] = -(far + near) / fmn
	return m
}
