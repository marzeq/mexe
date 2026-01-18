package evaluator

import (
	"fmt"
	"math/big"
	"strings"
)

const DefaultPrec = 1024

type Number struct {
	i *big.Int
	f *big.Float
}

func NewInt(v *big.Int) Number {
	return Number{i: new(big.Int).Set(v)}
}

func NewFloat(v *big.Float) Number {
	return Number{f: new(big.Float).Set(v)}
}

func Zero() Number {
	return NewInt(big.NewInt(0))
}

func One() Number {
	return NewInt(big.NewInt(1))
}

func ParseNumber(s string) (Number, error) {
	if strings.ContainsAny(s, ".eE") {
		f, ok := new(big.Float).SetPrec(DefaultPrec).SetString(s)
		if !ok {
			return Number{}, fmt.Errorf("invalid number: %s", s)
		}
		return NewFloat(f), nil
	}

	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return Number{}, fmt.Errorf("invalid number: %s", s)
	}
	return NewInt(i), nil
}

func (n Number) IsInt() bool {
	return n.i != nil
}

func (n Number) Int() *big.Int {
	if n.i != nil {
		return new(big.Int).Set(n.i)
	}
	i, _ := n.f.Int(nil)
	return i
}

func (n Number) Float() *big.Float {
	if n.f != nil {
		return new(big.Float).Set(n.f)
	}
	return new(big.Float).SetPrec(DefaultPrec).SetInt(n.i)
}

func promote(a, b Number) (Number, Number) {
	if a.IsInt() && b.IsInt() {
		return a, b
	}
	return NewFloat(a.Float()), NewFloat(b.Float())
}

func (n Number) Format(prec int) string {
	if n.i != nil {
		return n.i.String()
	}
	return n.f.Text('f', prec)
}
