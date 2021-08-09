package sec

import (
	"math"
	"math/big"
)

// https://en.wikipedia.org/wiki/Factorial_number_system
func permute(n int64, k *big.Int) []int64 {
	sets := make([]int64, n)
	sets[0] = 0
	for i := int64(1); i < n; i++ {
		sets[i] = i
	}

	pick := func(digit int64) int64 {
		j := int64(0)
		for i := 0; i < len(sets); i++ {
			if sets[i] != -1 {
				if j == digit {
					d := sets[i]
					sets[i] = -1
					return d
				}
				j++
			}
		}
		// unreachable
		return -1
	}

	var perm []int64
	digits := factoradic(n, k)

	for _, v := range digits {
		p := pick(v)
		perm = append(perm, p)
	}
	return perm
}

func factoradic(n int64, k *big.Int) []int64 {
	factorials := make([]*big.Int, n)
	factorials[0] = big.NewInt(1)
	for i := int64(1); i < n; i++ {
		factorials[i] = big.NewInt(0).Mul(factorials[i-1], big.NewInt(i))
	}

	var digits []int64
	kk := big.NewInt(0).Set(k)

	for i := n - 1; i >= 0; i-- {
		d := big.NewInt(0).Div(kk, factorials[i])
		digits = append(digits, d.Int64())
		kk.Sub(kk, d.Mul(d, factorials[i]))
	}
	return digits
}

func factorial(n int64) *big.Int {
	if n == 0 {
		return big.NewInt(1)
	}
	f := big.NewInt(1)
	for i := int64(1); i <= n; i++ {
		f.Mul(f, big.NewInt(i))
	}
	return f
}

// schedule retuns the permutation of n based on computation from data
func schedule(n int64, data []byte) ([]int64, error) {
	kdf := func(size int) ([]byte, error) {
		c := size/64 + 1
		ba, err := KDF(data, c, 64)
		if err != nil {
			return nil, err
		}
		var h []byte
		for _, k := range ba {
			h = append(h, k...)
		}
		return h[:size], nil
	}

	f := big.NewInt(0)

	x := Power2(n)
	if x > 0 {
		k := FNV(data, uint32(x+1))
		// derive bytes if bigger than 2^32
		if k < 32 {
			v := math.Pow(2, float64(k))
			f.SetInt64(int64(v))
		} else {
			b, err := kdf(int(k / 8))
			if err != nil {
				return nil, err
			}
			f.SetBytes(b)
		}
	}

	p := permute(n, f)

	return p, nil
}
