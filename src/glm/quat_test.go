package glm

import (
	"testing"
)

func TestSlerp(test *testing.T) {
	rest := Vector4{1, 0, 0, 1}
	axis := Vector3{0, 0, 1}
	q0 := axis.Quat(0)
	q1 := axis.Quat(180)
	for i := 0; i <= 10; i++ {
		t := float64(i) / 10.
		q := Slerp(q0, q1, t)
		M := q.Matrix()
		dest_m := M.MultV(rest)
		dest_q := q.Rotate4(rest)
		// Converting the quaternion to a matrix or using an Hamiltonian product
		// should yield the same result.
		distance := (dest_m.Sub(dest_q)).Norm()
		if distance >= 1.e-15 {
			test.Errorf("Matrix and Quaternion rotation differ by %v: %v ≠ %v.",
				distance, dest_m, dest_q)
		}
	}
}
