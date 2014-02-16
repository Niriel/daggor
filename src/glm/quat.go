package glm

import (
	"math"
)

type Quat [4]float64

func (v Vector3) Quat(angle float64) Quat {
	angle *= DEG_TO_RAD
	angle *= .5
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Quat{c, v[0] * s, v[1] * s, v[2] * s}
}

func (v Vector4) Quat(angle float64) Quat {
	return v.To3().Quat(angle)
}

func (q Quat) Neg() Quat {
	return Quat{-q[0], -q[1], -q[2], -q[3]}
}

func (q Quat) Smult(s float64) Quat {
	return Quat{q[0] * s, q[1] * s, q[2] * s, q[3] * s}
}

func (q0 Quat) Add(q1 Quat) Quat {
	return Quat{
		q0[0] + q1[0],
		q0[1] + q1[1],
		q0[2] + q1[2],
		q0[3] + q1[3],
	}
}

func (q Quat) Matrix() Matrix4 {
	// Source: wikipedia.  Need better source.
	w, x, y, z := q[0], q[1], q[2], q[3]
	x2 := x * x
	y2 := y * y
	z2 := z * z

	// [1 - 2 * y2 - 2 * z2     2 * x * y - 2 * z * w   2 * x * z + 2 * y * w]
	// [2 * x * y + 2 * z * w   1 - 2 * x2 - 2 * z2     2 * y * z - 2 * x * w]
	// [2 * x * z - 2 * y * w   2 * y * z + 2 * x * w   1 - 2 * x2 - 2 * y2  ]
	return Matrix4{
		1 - 2*y2 - 2*z2,
		2*x*y + 2*z*w,
		2*x*z - 2*y*w,
		0,
		2*x*y - 2*z*w,
		1 - 2*x2 - 2*z2,
		2*y*z + 2*x*w,
		0,
		2*x*z + 2*y*w,
		2*y*z - 2*x*w,
		1 - 2*x2 - 2*y2,
		0,
		0, 0, 0, 1,
	}
}

func (q0 Quat) Dot(q1 Quat) float64 {
	return q0[0]*q1[0] + q0[1]*q1[1] + q0[2]*q1[2] + q0[3]*q1[3]
}

func (q0 Quat) Mult(q1 Quat) Quat {
	//While waiting to find Hamilton's book:
	//@article{vicci2001quaternions,
	// title={Quaternions and rotations in 3-space: The algebra and its geometric interpretation},
	// author={Vicci, Leandra},
	// journal={Microelectronic Systems Laboratory, Departement of Computer Science, University of North Carolina at Chapel Hill},
	// year={2001},
	// publisher={Citeseer}
	//}
	return Quat{
		q0[0]*q1[0] - q0[1]*q1[1] - q0[2]*q1[2] - q0[3]*q1[3],
		q0[0]*q1[1] + q0[1]*q1[0] + q0[2]*q1[3] - q0[3]*q1[2],
		q0[0]*q1[2] + q0[2]*q1[0] + q0[3]*q1[1] - q0[1]*q1[3],
		q0[0]*q1[3] + q0[3]*q1[0] + q0[1]*q1[2] - q0[2]*q1[1],
	}
}

func (q Quat) Conj() Quat {
	return Quat{q[0], -q[1], -q[2], -q[3]}
}

func (q Quat) Hamilton(p Quat) Quat {
	return q.Mult(p.Mult(q.Conj()))
}

func (q Quat) Rotate3(v Vector3) Vector3 {
	p := Quat{0, v[0], v[1], v[2]}
	p1 := q.Hamilton(p)
	return Vector3{p1[1], p1[2], p1[3]}
}

func (q Quat) Rotate4(v Vector4) Vector4 {
	p := Quat{0, v[0], v[1], v[2]}
	p1 := q.Hamilton(p)
	return Vector4{p1[1], p1[2], p1[3], v[3]}
}
