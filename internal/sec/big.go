package sec

import (
	"math"
	"math/big"
)

// Power2 returns the max exponent x
// where 2^x < n!
func Power2(n int64) int64 {
	if n < 3 {
		return 0
	}
	f := factorial(n)

	pow := func(x int64) int64 {
		return int64(math.Pow(2, float64(x)))
	}

	base2 := big.NewInt(2)

	// lower/upper bound
	a := new(big.Int)
	b := new(big.Int)

	var i int64 = 0
	for {
		i++
		e := big.NewInt(pow(i))
		b.Exp(base2, e, nil)

		if b.Cmp(f) > 0 {
			break
		}
	}

	a.Exp(base2, big.NewInt(pow(i-1)), nil)
	for lo, hi := pow(i-1), pow(i); ; {
		m := (lo + hi) / 2

		if m == lo {
			return m
		}
		c := new(big.Int).Exp(base2, big.NewInt(m), nil)
		if c.Cmp(f) <= 0 {
			a.Set(c)
			lo = m
		} else {
			b.Set(c)
			hi = m
		}
	}
}
