package world

import (
	"glm"
	"math"
)

type Dynamic struct {
	T0 uint64 // Date of birth
}

func (self Dynamic) ModelMat(t uint64) glm.Matrix4 {
	const TAU = 1 * 1000000000 // Nanosecond.
	dt := float64(t - self.T0)
	z := .5 + .2*math.Sin(2*math.Pi*dt/TAU)
	pos := glm.Vector3{1, 0, z}
	return pos.Translation()
}

func (self Dynamic) Mesh(t uint64) []float64 {
	const p = .3                 // Plus sign.
	const m = -p                 // Minus sign.
	const TAUO = .6 * 1000000000 // Nanosecond.
	const TAUZ = .7 * 1000000000 // Nanosecond.
	dt := float64(t - self.T0)
	o := .08 * math.Sin(2*math.Pi*dt/TAUO)
	z := .10 * math.Sin(2*math.Pi*dt/TAUZ+.2)
	return []float64{
		m - o,
		m - o,
		0,

		m - o,
		p + o,
		0,

		p + o,
		m - o,
		0,

		p + o,
		p + o,
		0,

		0,
		0,
		1 + z,
	}
}
