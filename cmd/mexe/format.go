package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/marzeq/mexe/evaluator"
)

func formatNumber(n evaluator.Number, prec int) string {
	if n.IsInt() {
		return formatInt(n.Int())
	}
	return formatFloat(n.Float(), prec)
}

func formatInt(i *big.Int) string {
	if i.Sign() >= 0 {
		n := new(big.Int).Mul(i, i)
		if n.IsInt64() {
			k, m := splitSquare(int(n.Int64()))
			if m == 1 {
				return strconv.Itoa(k)
			}
			if k == 1 {
				return fmt.Sprintf("sqrt(%d)", m)
			}
			return fmt.Sprintf("%dsqrt(%d)", k, m)
		}
	}
	return i.String()
}

func formatFloat(f *big.Float, prec int) string {
	x, _ := f.Float64()

	eps := math.Pow(10, -float64(prec))
	if math.Abs(x) < eps {
		return "0"
	}

	if math.Abs(x) < 1e6 {
		if k := x / math.Pi; isNiceRel(k, prec) {
			return formatMul(k, "pi", prec)
		}
		if k := x / math.E; isNiceRel(k, prec) {
			return formatMul(k, "e", prec)
		}
	}

	s := f.Text('f', prec)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func splitSquare(n int) (k int, m int) {
	k = 1
	m = n
	for i := 2; i*i <= m; i++ {
		for m%(i*i) == 0 {
			m /= i * i
			k *= i
		}
	}
	return
}

func isNiceRel(x float64, prec int) bool {
	if x == 0 {
		return true
	}
	rel := math.Pow(10, -float64(prec))
	r := math.Round(x)
	return math.Abs(x-r)/math.Abs(r) < rel ||
		math.Abs(2*x-math.Round(2*x))/math.Abs(2*x) < rel
}

func formatMul(k float64, sym string, prec int) string {
	r := math.Round(k)
	if math.Abs(k-r) < math.Pow(10, -float64(prec)) {
		if r == 1 {
			return sym
		}
		if r == -1 {
			return "-" + sym
		}
		return fmt.Sprintf("%.0f%s", r, sym)
	}
	return fmt.Sprintf("%.*g%s", prec, k, sym)
}

