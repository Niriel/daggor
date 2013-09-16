// glm project glm.go
package glm

// This library provides some linear algebra functions.
// It is designed with OpenGl in mind.

// Matrices and vectors are stored in Go arrays.
// None of the methods defined here on matrices and arrays mutate the
// matrix or array.  A new matrix or array is always returned.
// You can bypass this and mutate the arrays yourself if you know what you
// are doing.  Remember however that mutables can be dangerous in a
// concurrent environment.

import (
	"fmt"
	"math"
)

const DEG_TO_RAD = float64(math.Pi) / 180

type Vector2 [2]float64
type Vector3 [3]float64
type Vector4 [4]float64
type Matrix4 [16]float64

var IDENTITY4 = Matrix4{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}
var YUP = IDENTITY4
var ZUP = Matrix4{
	1, 0, 0, 0,
	0, 0, -1, 0,
	0, 1, 0, 0,
	0, 0, 0, 1}

type TaitBryan struct {
	h, p, r int
}

type Attitude struct {
	tb     TaitBryan
	Angles Vector3
}

func (self Vector3) X() float64 {
	return self[0]
}

func (self Vector3) Y() float64 {
	return self[1]
}

func (self Vector3) Z() float64 {
	return self[2]
}

func (self Vector4) X() float64 {
	return self[0]
}

func (self Vector4) Y() float64 {
	return self[1]
}

func (self Vector4) Z() float64 {
	return self[2]
}

func (self Vector4) W() float64 {
	return self[3]
}

func (self Vector3) Xyz() (float64, float64, float64) {
	return self[0], self[1], self[2]
}

func (self Vector4) Xyz() (float64, float64, float64) {
	return self[0], self[1], self[2]
}

func (self Vector4) Xyzw() (float64, float64, float64, float64) {
	return self[0], self[1], self[2], self[3]
}

func (self Vector3) SetX(x float64) Vector3 {
	self[0] = x
	return self
}

func (self Vector3) SetY(y float64) Vector3 {
	self[1] = y
	return self
}

func (self Vector3) SetZ(z float64) Vector3 {
	self[2] = z
	return self
}

func (self Vector3) SetXyz(x, y, z float64) Vector3 {
	self[0] = x
	self[1] = y
	self[2] = z
	return self
}

func (self Vector4) SetX(x float64) Vector4 {
	self[0] = x
	return self
}

func (self Vector4) SetY(y float64) Vector4 {
	self[1] = y
	return self
}

func (self Vector4) SetZ(z float64) Vector4 {
	self[2] = z
	return self
}

func (self Vector4) SetW(w float64) Vector4 {
	self[3] = w
	return self
}

func (self Vector4) SetXyz(x, y, z float64) Vector4 {
	self[0] = x
	self[1] = y
	self[2] = z
	return self
}

func (self Vector4) SetXyzw(x, y, z, w float64) Vector4 {
	self[0] = x
	self[1] = y
	self[2] = z
	self[3] = w
	return self
}

func (self Vector3) To4(w float64) Vector4 {
	return Vector4{self[0], self[1], self[2], w}
}

func (self Vector4) To3() Vector3 {
	return Vector3{self[0], self[1], self[2]}
}

func (self Vector3) Add(other Vector3) Vector3 {
	return Vector3{
		self[0] + other[0],
		self[1] + other[1],
		self[2] + other[2],
	}
}

func (self Vector4) Add(other Vector4) Vector4 {
	return Vector4{
		self[0] + other[0],
		self[1] + other[1],
		self[2] + other[2],
		self[3] + other[3],
	}
}

func (self Vector3) Sub(other Vector3) Vector3 {
	return Vector3{
		self[0] - other[0],
		self[1] - other[1],
		self[2] - other[2],
	}
}

func (self Vector4) Sub(other Vector4) Vector4 {
	return Vector4{
		self[0] - other[0],
		self[1] - other[1],
		self[2] - other[2],
		self[3] - other[3],
	}
}

func (self Vector3) Smult(scalar float64) Vector3 {
	return Vector3{
		self[0] * scalar,
		self[1] * scalar,
		self[2] * scalar,
	}
}

func (self Vector4) Smult(scalar float64) Vector4 {
	return Vector4{
		self[0] * scalar,
		self[1] * scalar,
		self[2] * scalar,
		self[3] * scalar,
	}
}

func (self Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		self[1]*other[2] - self[2]*other[1],
		self[2]*other[0] - self[0]*other[2],
		self[0]*other[1] - self[1]*other[0],
	}
}

func (self Vector3) Dot(other Vector3) float64 {
	return self[0]*other[0] + self[1]*other[1] + self[2]*other[2]
}

func (self Vector4) Norm() float64 {
	return math.Sqrt(
		self[0]*self[0] +
			self[1]*self[1] +
			self[2]*self[2] +
			self[3]*self[3])
}

func (self Vector3) Norm() float64 {
	return float64(math.Sqrt(self[0]*self[0] + self[1]*self[1] + self[2]*self[2]))
}

func (self Vector3) Normed() Vector3 {
	return self.Smult(1 / self.Norm())
}

func (self Vector3) Translation() Matrix4 {
	var t Matrix4
	t[0] = 1
	t[5] = 1
	t[10] = 1
	t[15] = 1
	t[12] = self[0]
	t[13] = self[1]
	t[14] = self[2]
	return t
}

func (self Vector3) TranslationInv() Matrix4 {
	var t Matrix4
	t[0] = 1
	t[5] = 1
	t[10] = 1
	t[15] = 1
	t[12] = -self[0]
	t[13] = -self[1]
	t[14] = -self[2]
	return t
}

func (self Vector4) Translation() Matrix4 {
	var t Matrix4
	t[0] = 1
	t[5] = 1
	t[10] = 1
	t[15] = 1
	t[12] = self[0]
	t[13] = self[1]
	t[14] = self[2]
	return t
}

func (self Vector4) TranslationInv() Matrix4 {
	var t Matrix4
	t[0] = 1
	t[5] = 1
	t[10] = 1
	t[15] = 1
	t[12] = -self[0]
	t[13] = -self[1]
	t[14] = -self[2]
	return t
}

func (self Vector3) Gl() [3]float32 {
	return [3]float32{
		float32(self[0]),
		float32(self[1]),
		float32(self[2]),
	}
}

func (self Vector4) Gl() [4]float32 {
	return [4]float32{
		float32(self[0]),
		float32(self[1]),
		float32(self[2]),
		float32(self[3]),
	}
}

func MakeTaitBryan(h, p, r int) TaitBryan {
	if h < 0 || h > 3 {
		panic(fmt.Sprintln("Heading out of bound: %i not in [0, 1, 2].", h))
	}
	if p < 0 || p > 3 {
		panic(fmt.Sprintln("Pitch out of bound: %i not in [0, 1, 2].", p))
	}
	if r < 0 || r > 3 {
		panic(fmt.Sprintln("Roll out of bound: %i not in [0, 1, 2].", r))
	}
	if h == p {
		panic("Header and pitch must be different.")
	}
	if h == r {
		panic("Header and roll must be different.")
	}
	if p == r {
		panic("Pitch and roll must be different.")
	}
	var tb TaitBryan
	tb.h = h
	tb.p = p
	tb.r = r
	return tb
}

func (self TaitBryan) MakeAttitude() Attitude {
	var att Attitude
	att.tb = self
	return att
}

func (self Attitude) GetTb() TaitBryan {
	return self.tb
}

func (self Attitude) X() float64 {
	return self.Angles[0]
}

func (self Attitude) Y() float64 {
	return self.Angles[1]
}

func (self Attitude) Z() float64 {
	return self.Angles[2]
}

func (self Attitude) H() float64 {
	return self.Angles[self.tb.h]
}

func (self Attitude) P() float64 {
	return self.Angles[self.tb.p]
}

func (self Attitude) R() float64 {
	return self.Angles[self.tb.r]
}

func (self Attitude) Xyz() (float64, float64, float64) {
	return self.X(), self.Y(), self.Z()
}

func (self Attitude) Xyzv() Vector3 {
	return Vector3{self.X(), self.Y(), self.Z()}
}

func (self Attitude) Hpr() (float64, float64, float64) {
	return self.H(), self.P(), self.R()
}

func (self Attitude) Hprv() Vector3 {
	return Vector3{self.H(), self.P(), self.R()}
}

func (self Attitude) SetX(x float64) Attitude {
	self.Angles[0] = x
	return self
}

func (self Attitude) SetY(y float64) Attitude {
	self.Angles[1] = y
	return self
}

func (self Attitude) SetZ(z float64) Attitude {
	self.Angles[2] = z
	return self
}

func (self Attitude) SetH(h float64) Attitude {
	self.Angles[self.tb.h] = h
	return self
}

func (self Attitude) SetP(p float64) Attitude {
	self.Angles[self.tb.p] = p
	return self
}

func (self Attitude) SetR(r float64) Attitude {
	self.Angles[self.tb.r] = r
	return self
}

func (self Attitude) SetXyz(x, y, z float64) Attitude {
	self.Angles[0] = x
	self.Angles[1] = y
	self.Angles[2] = z
	return self
}

func (self Attitude) SetXyzv(xyz Vector3) Attitude {
	self.Angles = xyz
	return self
}

func (self Attitude) SetHpr(h, p, r float64) Attitude {
	self.Angles[self.tb.h] = h
	self.Angles[self.tb.p] = p
	self.Angles[self.tb.r] = r
	return self
}

func (self Attitude) SetHprv(hpr Vector3) Attitude {
	return self.SetHpr(hpr[0], hpr[1], hpr[2])
}

func (a Matrix4) Mult(b Matrix4) Matrix4 {
	return Matrix4{
		a[0]*b[0] + a[4]*b[1] + a[8]*b[2] + a[12]*b[3],
		a[1]*b[0] + a[5]*b[1] + a[9]*b[2] + a[13]*b[3],
		a[2]*b[0] + a[6]*b[1] + a[10]*b[2] + a[14]*b[3],
		a[3]*b[0] + a[7]*b[1] + a[11]*b[2] + a[15]*b[3],

		a[0]*b[4] + a[4]*b[5] + a[8]*b[6] + a[12]*b[7],
		a[1]*b[4] + a[5]*b[5] + a[9]*b[6] + a[13]*b[7],
		a[2]*b[4] + a[6]*b[5] + a[10]*b[6] + a[14]*b[7],
		a[3]*b[4] + a[7]*b[5] + a[11]*b[6] + a[15]*b[7],

		a[0]*b[8] + a[4]*b[9] + a[8]*b[10] + a[12]*b[11],
		a[1]*b[8] + a[5]*b[9] + a[9]*b[10] + a[13]*b[11],
		a[2]*b[8] + a[6]*b[9] + a[10]*b[10] + a[14]*b[11],
		a[3]*b[8] + a[7]*b[9] + a[11]*b[10] + a[15]*b[11],

		a[0]*b[12] + a[4]*b[13] + a[8]*b[14] + a[12]*b[15],
		a[1]*b[12] + a[5]*b[13] + a[9]*b[14] + a[13]*b[15],
		a[2]*b[12] + a[6]*b[13] + a[10]*b[14] + a[14]*b[15],
		a[3]*b[12] + a[7]*b[13] + a[11]*b[14] + a[15]*b[15],
	}
}

func (a Matrix4) MultV(b Vector4) Vector4 {
	return Vector4{
		a[0]*b[0] + a[4]*b[1] + a[8]*b[2] + a[12]*b[3],
		a[1]*b[0] + a[5]*b[1] + a[9]*b[2] + a[13]*b[3],
		a[2]*b[0] + a[6]*b[1] + a[10]*b[2] + a[14]*b[3],
		a[3]*b[0] + a[7]*b[1] + a[11]*b[2] + a[15]*b[3],
	}
}

func (a Matrix4) Gl() [16]float32 {
	var result [16]float32
	for i := 0; i < 16; i++ {
		result[i] = float32(a[i])
	}
	return result
}

func RotX(angle float64) Matrix4 {
	var r Matrix4
	a := angle * DEG_TO_RAD
	c := math.Cos(a)
	s := math.Sin(a)
	r[0] = 1
	r[5] = c
	r[6] = s
	r[9] = -s
	r[10] = c
	r[15] = 1
	return r
}

func RotY(angle float64) Matrix4 {
	var r Matrix4
	a := angle * DEG_TO_RAD
	c := math.Cos(a)
	s := math.Sin(a)
	r[0] = c
	r[2] = -s
	r[5] = 1
	r[8] = s
	r[10] = c
	r[15] = 1
	return r
}

func RotZ(angle float64) Matrix4 {
	var r Matrix4
	a := angle * DEG_TO_RAD
	c := math.Cos(a)
	s := math.Sin(a)
	r[0] = c
	r[1] = s
	r[4] = -s
	r[5] = c
	r[10] = 1
	r[15] = 1
	return r
}

func (self Attitude) Rotation() Matrix4 {
	rs := [3]Matrix4{RotX(self.X()), RotY(self.Y()), RotZ(self.Z())}
	rh := rs[self.tb.h]
	rp := rs[self.tb.p]
	rr := rs[self.tb.r]
	return rh.Mult(rp).Mult(rr)
}

func (self Attitude) RotationInv() Matrix4 {
	rs := [3]Matrix4{RotX(-self.X()), RotY(-self.Y()), RotZ(-self.Z())}
	rh := rs[self.tb.h]
	rp := rs[self.tb.p]
	rr := rs[self.tb.r]
	return rr.Mult(rp).Mult(rh)
}

// Perspective projection.
func PerspectiveProj(fov_h, aspect, z_near, z_far float64) Matrix4 {
	var m Matrix4
	neg_depth := z_near - z_far
	e := 1 / math.Tan(fov_h*.5*DEG_TO_RAD) // Focale length.
	m[0] = e
	m[5] = e * aspect
	m[10] = (z_near + z_far) / neg_depth
	m[14] = 2 * (z_near * z_far) / neg_depth
	m[11] = -1
	return m
}

// Orthographic projection.
func OrthographicProj(fov_h, aspect, z_near, z_far float64) Matrix4 {
	var m Matrix4
	neg_depth := z_near - z_far
	e := 1 / math.Tan(fov_h*.5*DEG_TO_RAD) // Focale length.
	m[0] = e
	m[5] = e * aspect
	m[10] = 2 / neg_depth
	m[14] = (z_near + z_far) / neg_depth
	m[15] = 1
	return m
}

// Compute the normal to a triangle from two of its sides.
func FaceNormal2(d10, d20 Vector3) Vector3 {
	d := d10.Cross(d20)
	return d.Smult(1 / d.Norm())
}

// Compute the normal to a triangle from its three vertices.
func FaceNormal3(v0, v1, v2 Vector3) Vector3 {
	return FaceNormal2(v1.Sub(v0), v2.Sub(v0))
}

func AngleU2(d0, d1 Vector3) float64 {
	n0 := d0.Normed()
	n1 := d1.Normed()
	cos := n0.Dot(n1)
	return math.Acos(cos)
}

func AngleU3(v0, v1, v2 Vector3) float64 {
	return AngleU2(v1.Sub(v0), v2.Sub(v0))
}
