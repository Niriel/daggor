package glm

import (
	"fmt"
	"math"
)

func CheckInterpolatorDomain(t float64) {
	if t < 0 || t > 1 {
		panic(fmt.Sprintf("Interpolator parameter out of [0,1]: %v.", t))
	}
}

// Function with first and second derivatives =0 at its extremities.
func InterpSmooth(t float64) float64 {
	CheckInterpolatorDomain(t)
	// y0(t) =   at5 +   bt4 +  ct3 + dt2 + et + f
	// y1(t) =  5at4 +  4bt3 + 3ct2 + 2dt + e
	// y2(t) = 20at3 + 12bt2 + 6ct  + 2d

	// Constraints at t=0 imply that d, e and f are zero.

	// y0(t) =   at5 +   bt4 +  ct3
	// y1(t) =  5at4 +  4bt3 + 3ct2
	// y2(t) = 20at3 + 12bt2 + 6ct

	// y0(1) =   a +   b +  c = 1
	// y1(1) =  5a +  4b + 3c = 0
	// y2(1) = 20a + 12b + 6c = 0

	//         [1   1   1][a]   [1]        [a] [  6]
	// Solving [5   4   3][b] = [0] yields [b]=[-15]
	//         [20 12   6][c]   [0]        [c] [ 10]

	// y0(.5) = .5.
	// Better still:
	// y0(t) + y0(1-t) = 1, which means that the function is symetrical.

	t3 := t * t * t
	t4 := t3 * t
	return 6*t4*t - 15*t4 + 10*t3
}

func Slerp(q0, q1 Quat, t float64) Quat {
	CheckInterpolatorDomain(t)
	cos_angle := q0.Dot(q1)
	if cos_angle < 0 {
		q0 = q0.Neg() // Ensure shortest path.
	}
	var k0, k1 float64
	if cos_angle > .9999 {
		// Linear interpolation to avoid divisions by tiny numbers.
		// Note: does not guarantee normalized quaternions.
		k0 = 1 - t
		k1 = t
	} else {
		sin_angle := math.Sqrt(1 - cos_angle*cos_angle)
		angle := math.Atan2(sin_angle, cos_angle)
		one_over_sin_angle := 1 / sin_angle
		k0 = math.Sin(angle*(1-t)) * one_over_sin_angle
		k1 = math.Sin(angle*t) * one_over_sin_angle
	}
	return Quat{
		k0*q0[0] + k1*q1[0],
		k0*q0[1] + k1*q1[1],
		k0*q0[2] + k1*q1[2],
		k0*q0[3] + k1*q1[3],
	}
}

func InterpolateTranslation(v0, v1 Vector3, t float64) Vector3 {
	switch {
	case t <= 0:
		return v0
	case t >= 1:
		return v1
	case v0 == v1:
		return v0
	}
	r := 1 - t
	return Vector3{
		v0[0]*r + v1[0]*t,
		v0[1]*r + v1[1]*t,
		v0[2]*r + v1[2]*t,
	}
}

func InterpolateRotation(q0, q1 Quat, t float64) Quat {
	switch {
	case t <= 0:
		return q0
	case t >= 1:
		return q1
	case q0 == q1:
		return q0
	}
	return Slerp(q0, q1, t)
}
