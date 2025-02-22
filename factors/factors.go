package factors

import (
	"fmt"
	"iter"
	"os"
	"strings"
)

var smallPrime = map[uint64]bool{
	1:  true,
	2:  true,
	3:  true,
	5:  true,
	7:  true,
	11: true,
	13: true,
	17: true,
	19: true,
	23: true,
	29: true,
	31: true,
	37: true,
	41: true,
	43: true,
	47: true,
	53: true,
	59: true,
	61: true,
	67: true,
	71: true,
	73: true,
	79: true,
	83: true,
	89: true,
	97: true,
}

func pow2Iter() iter.Seq2[int, uint64] {
	return func(yield func(int, uint64) bool) {
		var x uint64 = 2
		i := 1
		var m uint64 = 1 << 63
		for x < m {
			x <<= 1
			i++
			if !yield(i, x) {
				return
			}
		}
	}
}

func makeIter(a, b uint64) iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		head := a
		for head < b {
			head++
			if !yield(head) {
				return
			}
		}
	}
}

func isPrime(k uint64) bool {
	if k == 0 {
		return false
	}

	if k < 100 {
		_, ok := smallPrime[k]
		return ok
	}

	kr := uint64(k / 2)

	for i := range makeIter(2, kr) {
		if k%i == 0 {
			return false
		}
	}

	return true
}

func findFactors(k uint64) (out string) {
	defer func() {
		out = strings.TrimSpace(out) + ";"
	}()

	if isPrime(k) {
		out += fmt.Sprintf("%d", k)
		return
	}

	k_ := k

outer:
	for k_ > 1 {
		if isPrime(k_) {
			out += fmt.Sprintf(" %d", k_)
			return
		}

		for i, p2 := range pow2Iter() {
			if k_ < p2 {
				break
			}

			if k_ == p2 {
				out += strings.Repeat(" 2", i)
				return
			}
		}

		// dont care about floats
		for i := range makeIter(1, k-1) {
			if (k_ % i) == 0 {
				out += fmt.Sprintf(" %d", i)
				if !isPrime(i) {
					out += "*"
				}
				k_ = k_ / i
				goto outer
			}
		}
		out += fmt.Sprintf(" %d", k_)
		return
	}
	return
}

func toInt64(read []byte) int64 {

	var value int32
	value |= int32(read[3])
	value |= int32(read[2]) << 8
	value |= int32(read[1]) << 16
	value |= int32(read[0]) << 24
	return int64(value)
}

func Read(n, nIdx string, k int64) string {
	f, err := os.Open(n)
	if err != nil {
		panic(err)
	}

	fIdx, err := os.Open(nIdx)
	if err != nil {
		panic(err)
	}

	ref1 := make([]byte, 4)
	ref2 := make([]byte, 4)

	i := 4 * (k - 1)

	fIdx.ReadAt(ref1, i)
	fIdx.ReadAt(ref2, i+4)

	start := toInt64(ref1) / 8
	end := toInt64(ref2) / 8

	leng := end - start

	val := make([]byte, leng-1)
	f.ReadAt(val, start)

	return string(val)
}

func Make(file, table *os.File) {
	var a uint64 = 0
	var n uint64 = 2<<15 - 1
	//a = 32726
	//n = 1

	var offset int

	for k := range makeIter(a, a+n) {
		factors := findFactors(k)
		b32 := make([]byte, 4)
		b32[3] = byte(offset)
		b32[2] = byte(offset >> 8)
		b32[1] = byte(offset >> 16)
		b32[0] = byte(offset >> 24)

		table.Write(b32)
		file.Write([]byte(factors))
		offset += len(factors) * 8
	}
}
