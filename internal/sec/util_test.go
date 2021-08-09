package sec

import (
	"encoding/hex"
	"unicode/utf8"

	"testing"
)

func TestRandomBytes(t *testing.T) {
	size := 6
	b, err := RandomBytes(size)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != size {
		t.Fatalf("got: %v expected: %v", len(b), size)
	}
	t.Logf("%s", hex.EncodeToString(b))
}

func TestS2b(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"ABC€", 6},
		{"Hello,世界", 12},
	}
	for i, tc := range tests {
		b := S2b(tc.input)
		if len(b) != tc.expected {
			t.Fatalf("[%v] got: %v %v %v want: %v", i, len(b), utf8.RuneCountInString(tc.input), b, tc.expected)
		}
	}
}

func TestB2s(t *testing.T) {
	tests := []struct {
		input    []byte
		expected int
	}{
		{[]byte{65, 66, 67, 226, 130, 172}, 6},
		{[]byte{72, 101, 108, 108, 111, 44, 228, 184, 150, 231, 149, 140}, 12},
	}
	for i, tc := range tests {
		s := B2s(tc.input)
		if len(s) != tc.expected {
			t.Fatalf("[%v] got: %v %v %v want: %v", i, len(s), utf8.RuneCountInString(s), s, tc.expected)
		}
	}
}
