package sec

import (
	"testing"
)

func TestPower2(t *testing.T) {
	tests := []struct {
		n        int64
		expected int64
	}{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 2},
		{4, 4},
		{5, 6},
		{6, 9},
		{7, 12},
		{8, 15},
		{9, 18},
		{10, 21},
		{11, 25},
		{12, 28},
		{13, 32},
		{14, 36},
		{15, 40},
		{16, 44},
		{94, 485},
		{1024, 8769},
		{4096, 43250},
	}
	for i, tc := range tests {
		got := Power2(tc.n)
		if got != tc.expected {
			t.Fatalf("[%v] got: %v want: %v", i, got, tc.expected)
		}
	}
}
