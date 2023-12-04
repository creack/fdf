package math3

import (
	"testing"
)

func TestMatMul(t *testing.T) {
	m1 := Matrix{
		{2, 3, 4},
		{3, 5, 6},
		{4, 5, 3},
	}
	m2 := Matrix{
		{1, 2, 1},
		{-1, 2, 1},
		{3, 2, 1},
	}
	expected := Matrix{
		{11, 18, 9},
		{16, 28, 14},
		{8, 24, 12},
	}
	if got := m1.Multiply(m2); got != expected {
		t.Errorf("Unexpected matrix multiplication result.\nGot:      %v\nExpected: %v", got, expected)
	}
}
