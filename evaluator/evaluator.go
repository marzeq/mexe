package evaluator

import (
	"fmt"
	"math"
	"math/big"

	"github.com/marzeq/mexe/parser"
)

func Evaluate(expr parser.Node, degreesMode bool) (Number, error) {
	return evaluate(expr, degreesMode, &scope{constants: Constants, functions: map[string]parser.Definition{}}, 0)
}

type scope struct {
	constants map[string]Number
	functions map[string]parser.Definition
	locals    map[string]Number
}

func EvaluateProgram(program *parser.Program, degreesMode bool) (Number, error) {
	s := &scope{constants: make(map[string]Number), functions: make(map[string]parser.Definition)}
	for name, value := range Constants {
		s.constants[name] = value
	}
	for _, def := range program.Definitions {
		if def.Params != nil {
			delete(s.constants, def.Name)
			s.functions[def.Name] = def
			continue
		}
		value, err := evaluate(def.Value, degreesMode, s, 0)
		if err != nil {
			return Zero(), fmt.Errorf("constant %s: %w", def.Name, err)
		}
		delete(s.functions, def.Name)
		s.constants[def.Name] = value
	}
	return evaluate(program.Expression, degreesMode, s, 0)
}

func evaluate(expr parser.Node, degreesMode bool, s *scope, depth int) (Number, error) {
	if depth > 1000 {
		return Zero(), fmt.Errorf("function call depth exceeded")
	}
	switch n := expr.(type) {

	case *parser.ConstantNode:
		if v, err := ParseNumber(n.Name); err == nil {
			return v, nil
		}
		if c, ok := s.locals[n.Name]; ok {
			return c, nil
		}
		if c, ok := s.constants[n.Name]; ok {
			return c, nil
		}
		return Zero(), fmt.Errorf("unknown constant: %s", n.Name)

	case *parser.UnaryOpNode:
		val, err := evaluate(n.Child, degreesMode, s, depth)
		if err != nil {
			return Zero(), err
		}
		switch n.OpType {
		case parser.UnaryOpTypeNegate:
			if val.IsInt() {
				return NewInt(new(big.Int).Neg(val.Int())), nil
			}
			return NewFloat(new(big.Float).Neg(val.Float())), nil

		case parser.UnaryOpTypeAbs:
			if val.IsInt() {
				return NewInt(new(big.Int).Abs(val.Int())), nil
			}
			return NewFloat(new(big.Float).Abs(val.Float())), nil

		case parser.UnaryOpTypeFact:
			return fact(val)

		default:
			return Zero(), fmt.Errorf("unknown unary operator")
		}

	case *parser.BinaryOpNode:
		lhs, err := evaluate(n.Left, degreesMode, s, depth)
		if err != nil {
			return Zero(), err
		}
		rhs, err := evaluate(n.Right, degreesMode, s, depth)
		if err != nil {
			return Zero(), err
		}

		switch n.OpType {
		case parser.BinaryOpTypeAdd:
			a, b := promote(lhs, rhs)
			if a.IsInt() {
				return NewInt(new(big.Int).Add(a.Int(), b.Int())), nil
			}
			return NewFloat(new(big.Float).Add(a.Float(), b.Float())), nil

		case parser.BinaryOpTypeSubtract:
			a, b := promote(lhs, rhs)
			if a.IsInt() {
				return NewInt(new(big.Int).Sub(a.Int(), b.Int())), nil
			}
			return NewFloat(new(big.Float).Sub(a.Float(), b.Float())), nil

		case parser.BinaryOpTypeMultiply:
			a, b := promote(lhs, rhs)
			if a.IsInt() {
				return NewInt(new(big.Int).Mul(a.Int(), b.Int())), nil
			}
			return NewFloat(new(big.Float).Mul(a.Float(), b.Float())), nil

		case parser.BinaryOpTypeDivide:
			if rhs.IsInt() && rhs.Int().Sign() == 0 {
				return Zero(), fmt.Errorf("division by zero")
			}
			return NewFloat(new(big.Float).Quo(lhs.Float(), rhs.Float())), nil

		case parser.BinaryOpTypeExponent:
			if lhs.IsInt() && rhs.IsInt() && rhs.Int().Sign() >= 0 {
				return NewInt(new(big.Int).Exp(lhs.Int(), rhs.Int(), nil)), nil
			}
			lhsf, _ := lhs.Float().Float64()
			rhsf, _ := rhs.Float().Float64()
			f := math.Pow(
				lhsf,
				rhsf,
			)
			return NewFloat(big.NewFloat(f).SetPrec(DefaultPrec)), nil

		case parser.BinaryOpTypeMod:
			if !lhs.IsInt() || !rhs.IsInt() {
				return Zero(), ErrIntegerExpected
			}
			return NewInt(new(big.Int).Mod(lhs.Int(), rhs.Int())), nil

		default:
			return Zero(), fmt.Errorf("unknown binary operator")
		}

	case *parser.FuncApplicationNode:
		args := make([]Number, len(n.Args))
		for i, a := range n.Args {
			v, err := evaluate(a, degreesMode, s, depth)
			if err != nil {
				return Zero(), err
			}
			args[i] = v
		}

		if def, ok := s.functions[n.FuncName]; ok {
			if len(args) != len(def.Params) {
				return Zero(), fmt.Errorf("function %s: %w", n.FuncName, ErrInvalidArgCount)
			}
			locals := make(map[string]Number, len(args))
			for i, param := range def.Params {
				locals[param] = args[i]
			}
			return evaluate(def.Value, degreesMode, &scope{constants: s.constants, functions: s.functions, locals: locals}, depth+1)
		}
		fn, ok := Functions[n.FuncName]
		if !ok {
			return Zero(), fmt.Errorf("unknown function: %s", n.FuncName)
		}
		return fn(args, degreesMode)

	default:
		return Zero(), fmt.Errorf("unknown node type")
	}
}
