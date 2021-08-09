package sec

import (
	"testing"
)

func TestSortString(t *testing.T) {
	// expected: !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~
	expected := make([]byte, 127-33)
	for i := 33; i < 127; i++ {
		expected[i-33] = byte(i)
	}
	s := sortString(DefaultCodeset)
	if s != string(expected) {
		t.Fatalf("got: %s %v expected: %s %v", s, len(s), expected, len(expected))
	}
	t.Log(s)
}

func keyTestdata() ([]byte, []byte) {
	salt := `
	Although it is not possible to "decrypt" password hashes to obtain the original passwords, it is possible to "crack" the hashes in some circumstances.

	The basic steps are:
	
	Select a password you think the victim has chosen

	Calculate the hash
	Compare the hash you calculated to the hash of the victim. If they match, you have correctly "cracked" the hash and now know the plaintext value of their password.
	This process is repeated for a large number of potential candidate passwords. Different methods can be used to select candidate passwords, including:
	
	- Lists of passwords obtained from other compromised sites
	- Brute force (trying every possible candidate)
	- Dictionaries or wordlists of common passwords
	While the number of permutations can be enormous with high speed hardware (such as GPUs) and cloud services with many servers for rent, the cost to an attacker is relatively small to do successful password cracking especially when best practices for hashing are not followed.
	
	Strong passwords stored with modern hashing algorithms and using hashing best practices should be effectively impossible for an attacker to crack.
	`

	secret := `
	The best way to protect sensitive information is to not store it in the first place.
	`
	return []byte(secret), []byte(salt)
}

func TestKeyGen(t *testing.T) {
	secret, salt := keyTestdata()
	masks := ""

	codebook := MakeCodebook(DefaultCodeset, masks)
	gen := KeyGen(codebook, []byte(secret), []byte(salt), 64, 3)
	if gen == nil {
		t.Fatalf("failed to create gen")
	}

	testdata := []struct {
		domain   string
		user     string
		pepper   string
		expected string
	}{
		{"example.com", "user", "", "6TfUT'^D!glUg~O8+bAT)z_/7G<O!KNL369(%=!?0~;J=`]&K&<{\"Jsawy>lol'N"},
	}

	for _, tc := range testdata {
		pwd, err := gen(tc.domain, tc.user, tc.pepper, 1)
		if err != nil {
			t.Fatal(err)
		}
		if pwd[0] != tc.expected {
			t.Fatalf("got: %s expected: %s", pwd[0], tc.expected)
		}
		t.Log(pwd)
	}
}

func TestMaskCodebook(t *testing.T) {
	testdata := []struct {
		mask     string
		expected string
	}{
		{"!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~", ""},
		{alpha + numeric + symbol, ""},
		{alpha + numeric, symbol},
		{alpha + symbol, numeric},
		{numeric + symbol, alpha},
		{alpha, sortString(numeric + symbol)},
		{"", sortString(DefaultCodeset)},
		{"\"\\", "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"},
	}
	for i, d := range testdata {
		cb := MakeCodebook(DefaultCodeset, d.mask)
		if cb != d.expected {
			t.Fatalf("[%v] got: %s expected: %s", i, cb, d.expected)
		}
		t.Log(cb)
	}
}

func TestSpaceOut(t *testing.T) {
	testdata := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"1", "1"},
		{"12", "12"},
		{"123", "123"},
		{"1234", "1234"},
		{"12345", "12345"},
		{"123456", "12345 6"},
		{"1234567", "12345 67"},
		{"12345678", "12345 678"},
		{"123456789", "12345 6789"},
		{"1234567890", "12345 67890"},
		{"1234567890a", "12345 67890a"},
		{"1234567890ab", "12345 67890a b"},
		{"1234567890abc", "12345 67890a bc"},
		{"1234567890abcd", "12345 67890a bcd"},
		{"1234567890abcde", "12345 67890a bcde"},
		{"1234567890abcdef", "12345 67890a bcdef"},
	}
	for i, v := range testdata {
		s := SpaceOut(v.input, 0)
		if s != v.expected {
			t.Fatalf("[%v] got: %s| expected: %s|", i, s, v.expected)
		}
	}
}

func BenchmarkKeyGen8(b *testing.B) {
	benchmarkKeyGen(8, b)
}

func BenchmarkKeyGen12(b *testing.B) {
	benchmarkKeyGen(12, b)
}

func BenchmarkKeyGen16(b *testing.B) {
	benchmarkKeyGen(16, b)
}

func benchmarkKeyGen(iteration int, b *testing.B) {
	secret, salt := keyTestdata()
	codebook := MakeCodebook(DefaultCodeset, "")
	gen := KeyGen(codebook, []byte(secret), []byte(salt), 64, iteration)
	domain := "example.com"
	user := "user@example.com"
	pepper := "random text"
	for i := 0; i < b.N; i++ {
		gen(domain, user, pepper, 100)
	}
}
