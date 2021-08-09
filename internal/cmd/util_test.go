package cmd

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	testdata := "This\r is\n a\ttest\r\n"
	expected := "Thisisatest"
	s := normalize(testdata)
	if s != expected {
		t.Fatalf("got: %s expected: %s", s, expected)
	}
}

func TestSplit2(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{"", ""}},
		{"1", []string{"1", "1"}},
		{"12", []string{"12", "12"}},
		{"123", []string{"123", "123"}},
		{"1234", []string{"123", "234"}},
		{"1234567890", []string{"1234567", "4567890"}},
		{"1234567890abcdef", []string{"1234567890ab", "567890abcdef"}},
		{"1234567890abcdefghij", []string{"1234567890abcd", "7890abcdefghij"}},
	}
	for i, tc := range tests {
		s := split2([]byte(tc.input))
		if string(s[0]) != tc.expected[0] || string(s[1]) != tc.expected[1] {
			t.Fatalf("got: %s %s want: %v", string(s[0]), string(s[1]), tc.expected)
		}
		t.Logf("[%v] a: %s b: %s", i, string(s[0]), string(s[1]))
	}
}
