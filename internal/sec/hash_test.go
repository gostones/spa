package sec

import (
	"strings"
	"testing"
)

func TestKDF(t *testing.T) {
	testdata := []struct {
		secret   string
		expected string
	}{
		{"", "9Fxj7+8/zuw8dZ0UbCOepPjUyC/B1XRcaaN5X2uEJUcfd4Ge6MPXAOUtsugY02f1AcElaSyAORzsGU/YhC3WrA=="},
		{"a", "nhPwkV8ZLliinFophOwaafCP0js+uABtJICIfe4rCkASSXkqP7vEoB+g2SJBpNum9a3+ODu/TwIhh9Y3NzE7Pg=="},
		{"b", "h3r4xJJCKYUkwW2QFuhwy0jJp/7a2s/qAp6JaxL6vrDfy8EHtLGAm7YfwCM9bNh1IE3ubq+LPMEuXfoeDbKXwQ=="},
		{"c", "V3OB9nPvmF8ZGHBNDaSbIhx4lQIgILoQyICXOnqf9z8JpOHbN9i59r011m9moEuq07HVV91p5LbPcbjhlRysxQ=="},
	}

	for i, tc := range testdata {
		k, err := KDF([]byte(tc.secret), 255, 64)
		if err != nil {
			t.Fatal(err)
		}
		s := Base64(k[0])
		if s != tc.expected {
			t.Fatalf("[%v] got: %s want: %s", i, s, tc.expected)
		}
	}
}

func TestSPA(t *testing.T) {
	secret := `The best way to protect sensitive information is to not store it in the first place.`
	salt := "himalayan pink salt"

	testdata := []struct {
		n        int
		domain   string
		user     string
		expected string
	}{
		{1, "example.com", "user", "8bDQQwLWTdDzVMoF1tepJ0iQEXK3pYJYGpzIqKrvW3UWfkbj7DQs1xmJKlykIbYTrnoOZLQ2gGCWo5v+XsFWVQ=="},
		{2, "example.com", "user", "YsYHan9+H/OkU/P0jV3bW5n+t6eXrhT29STrC6oVy6LAwKn+P153arecJMafg6u2zLXgRyIzG9eDUbEWItTReg=="},
		{3, "example.com", "user", "QEny7UrNCagsQgAiEb007kbmEdLr/mUk0OovWD0sHlA2EzC5XbiErPlUcP9t1LgQ/ttvaw7SOwqq6dNtjjJVZQ=="},
		{4, "example.com", "user", "GGwsc1xkX8TxVuIGoHQjJoSWsMqmy12SJ2FOAnmJupL4bEhB12jzL/DkXBNwn+Jkr2Rjz7SErpzim3mZKj+5iA=="},
	}

	info := func(domain, user, pepper string) []byte {
		s := strings.Join([]string{pepper, user, domain}, ":")
		return []byte(s)
	}

	for i, v := range testdata {
		k, err := SPA([]byte(secret), info(salt, v.domain, v.user), 64, v.n)
		if err != nil {
			t.Fatal(err)
		}
		s := Base64(k)
		if s != v.expected {
			t.Fatalf("[%v] got: %s expected: %s", i, s, v.expected)
		}
	}
}

func TestHMAC(t *testing.T) {
	tests := []struct {
		domain   string
		user     string
		pepper   string
		pin      string
		expected string
	}{
		{"example.com", "user", "", "1234", "TpARyw+6Ot0ZeuJtFFpWTl8JHRUrIKZxAJ99cQKVv1N1fzPR66i8MCpo0jKgXmSIDaayxlmW0qNiTqW5F0PeCw=="},
		{"example.com", "user", "a", "1234", "NsuDIdGRLFGIjfJaHHTxxvGgxPXdnLYrIkq00WC076bOBar/qoeTdbJ6j7HwZA23iBK8XeX8V5e1O8yu9CuhJw=="},
	}

	for _, tc := range tests {
		du := strings.Join([]string{tc.domain, tc.user}, ":")
		z := HMAC([]byte(tc.pepper), []byte(du))
		s := Base64(z)
		if s != tc.expected {
			t.Fatalf("got: %s wanted: %s", s, tc.expected)
		}
	}
}

// data for benchmarking hash functions
func hashTestdata() ([]byte, []byte, []byte) {
	secret := `The best way to protect sensitive information is to not store it in the first place.`
	salt := `himalayan pink salt`
	pepper := `variant per site`
	return []byte(secret), []byte(salt), []byte(pepper)
}

func BenchmarkKDF(b *testing.B) {
	secret, _, _ := hashTestdata()

	for i := 0; i < b.N; i++ {
		KDF(secret, 1, 64)
	}
}

func BenchmarkPbkdf2Key(b *testing.B) {
	secret, salt, _ := hashTestdata()

	for i := 0; i < b.N; i++ {
		pbkdf2Key(secret, salt, 64)
	}
}

func BenchmarkAargon2idKey(b *testing.B) {
	secret, salt, _ := hashTestdata()

	for i := 0; i < b.N; i++ {
		argon2idKey(secret, salt, 64)
	}
}

func BenchmarkScryptKey(b *testing.B) {
	secret, salt, _ := hashTestdata()

	for i := 0; i < b.N; i++ {
		scryptKey(secret, salt, 64)
	}
}

func BenchmarkSPA(b *testing.B) {
	secret, salt, _ := hashTestdata()

	for i := 0; i < b.N; i++ {
		SPA(secret, salt, 64, 1)
	}
}

func BenchmarkSPAKDF8(b *testing.B) {
	benchmarkSPAKDF(8, b)
}

func BenchmarkSPAKDF12(b *testing.B) {
	benchmarkSPAKDF(12, b)
}

func BenchmarkSPAKDF16(b *testing.B) {
	benchmarkSPAKDF(16, b)
}

func BenchmarkSPAKDF1024(b *testing.B) {
	benchmarkSPAKDF(1024, b)
}

func benchmarkSPAKDF(iteration int, b *testing.B) {
	secret, salt, _ := hashTestdata()
	raw := []byte(secret)
	for i := 0; i < b.N; i++ {
		SPAKDF(raw, salt, 1024, 64, iteration)
	}
}

func BenchmarkHMAC(b *testing.B) {
	secret, salt, _ := hashTestdata()

	for i := 0; i < b.N; i++ {
		HMAC(secret, salt)
	}
}

func BenchmarkFNV(b *testing.B) {
	data, _, _ := hashTestdata()

	for i := 0; i < b.N; i++ {
		FNV(data, 1024)
	}
}
