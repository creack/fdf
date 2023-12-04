package math3

import (
	"fmt"
	"testing"
)

func TestMatMul(t *testing.T) {
	m1 := Matrix{
		i: Vec{X: 2, Y: 3, Z: 4},
		j: Vec{X: 3, Y: 5, Z: 5},
		k: Vec{X: 4, Y: 6, Z: 3},
	}
	m2 := Matrix{
		i: Vec{X: 1, Y: -1, Z: 3},
		j: Vec{X: 2, Y: 2, Z: 2},
		k: Vec{X: 1, Y: 1, Z: 1},
	}

	fmt.Printf("%v\n", m1.Multiply(m2))
}
