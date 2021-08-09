package sec

import (
	"sort"
	"strings"
)

const alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const numeric = "0123456789"

//https://owasp.org/www-community/password-special-characters
const symbol = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

const AlphaNumeric = alpha + numeric
const AlphaNumericSymbol = alpha + numeric + symbol

const DefaultCodeset = alpha + numeric + symbol
const DefaultDivisor = 4

// https://docs.oracle.com/cd/E11223_01/doc.910/e11197/app_special_char.htm#MCMAD416
const oracleSymbol = "@%+\\/'!#$^?:,(){}[]~`-_."
const OracleCodeset = alpha + numeric + oracleSymbol

const EncloseEscape = "\"'()<>[]`{}\\"

// sortString sorts string.
func sortString(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i int, j int) bool { return b[i] < b[j] })
	return string(b)
}

func MakeCodebook(codeset, mask string) string {
	if mask == "" {
		return sortString(codeset)
	}
	lookup := make(map[byte]bool)
	for _, v := range []byte(mask) {
		lookup[v] = true
	}
	cnt := 0
	for _, v := range []byte(codeset) {
		if _, ok := lookup[v]; ok {
			cnt++
		}
	}
	ncs := make([]byte, len(codeset)-cnt)
	for i, j := 0, 0; j < len(codeset); j++ {
		if _, ok := lookup[codeset[j]]; ok {
			continue
		}
		ncs[i] = codeset[j]
		i++
	}
	return sortString(string(ncs))
}

func KeyGen(codebook string, x, y []byte, keyLen, iteration int) func(string, string, string, int) ([]string, error) {
	pud := func(domain, user, pepper string) []byte {
		s := strings.Join([]string{pepper, user, domain}, ":")
		return []byte(s)
	}

	hi := func(x, y []byte, max int) int {
		b := HMAC(x, y)
		return int(FNV(b, uint32(max)))
	}

	pick := func(z []byte) ([]byte, []byte) {
		i := hi(x[0:keyLen], z, len(x)-keyLen)
		secret := x[i : i+keyLen]
		j := hi(secret, z, len(y)-keyLen)
		salt := y[j : j+keyLen]
		return secret, salt
	}

	encode := func(b []byte, shift int) string {
		enc := make([]byte, len(b))
		for i, v := range b {
			idx := (shift + int(v)) % len(codebook)
			enc[i] = codebook[idx]
		}
		return string(enc)
	}

	return func(domain, user, pepper string, keyCount int) ([]string, error) {
		z := pud(domain, user, pepper)
		secret, salt := pick(z)

		keys, err := SPAKDF(HMAC(secret, z), salt, keyCount, keyLen, iteration)
		if err != nil {
			return nil, err
		}

		var enc []string
		for _, v := range keys {
			enc = append(enc, encode(v, 14))
		}
		return enc, nil
	}
}

// SpaceOut inserts spaces.
func SpaceOut(s string, divisor int) string {
	if divisor <= 0 {
		divisor = DefaultDivisor
	}
	if len(s) <= divisor {
		return s
	}

	const space = 0x20

	b := []byte(s)
	// count the number of spaces to insert.
	// the next location of the space is computed using:
	// (byte value % divisor + divisor)
	var cnt int
	for i := 0; ; {
		n := (int(b[i])%divisor + divisor)
		i += n
		if i >= len(b) {
			break
		}
		cnt++
	}
	if cnt == 0 {
		return s
	}

	sb := make([]byte, len(b)+cnt)
	copy := func(i, j, n int) {
		for k := 0; k < n; k++ {
			if i+k >= len(b) {
				return
			}
			sb[j+k] = b[i+k]
		}
		if j+n >= len(sb) {
			return
		}
		sb[j+n] = space
	}

	for i, j := 0, 0; i < len(b); {
		n := (int(b[i])%divisor + divisor)
		copy(i, j, n)
		i += n
		j += n + 1
	}
	return string(sb)
}
