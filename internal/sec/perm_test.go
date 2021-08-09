package sec

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestPermute(t *testing.T) {
	testdata := []struct {
		n        int64
		expected [][]int64
	}{
		{1, [][]int64{
			{0},
		}},
		{2, [][]int64{
			{0, 1}, {1, 0},
		}},
		{3, [][]int64{
			{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0},
		}},
	}
	for _, v := range testdata {
		for k, e := range v.expected {
			p := permute(v.n, big.NewInt(int64(k)))
			if !reflect.DeepEqual(p, e) {
				t.Fatalf("got: %v expected: %v", p, e)
			}
		}
	}
}

func TestFactorial(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
	}

	f := factorial(1024)
	t.Logf("%v %v", f, len(fmt.Sprintf("%v", f)))
}

func TestSchedule(t *testing.T) {
	data := []byte("salt")
	tests := []struct {
		n        int64
		expected []int64
	}{
		{1, []int64{0}},
		{2, []int64{0, 1}},
		{3, []int64{0, 2, 1}},
		{8, []int64{0, 1, 2, 3, 4, 6, 5, 7}},
		{1024, nil},
	}
	for _, tc := range tests {
		p, err := schedule(tc.n, data)
		if err != nil {
			t.Fatal(err)
		}
		if tc.expected != nil && !reflect.DeepEqual(p, tc.expected) {
			t.Fatalf("got: %v expected: %v", p, tc.expected)
		}
	}
}
