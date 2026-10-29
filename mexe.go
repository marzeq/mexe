package mexe

import (
	"fmt"

	"github.com/marzeq/mexe/evaluator"
	"github.com/marzeq/mexe/parser"
	"github.com/marzeq/mexe/tokeniser"
)

func Evaluate(expression string, degreesMode bool) (evaluator.Number, error) {
	tokens, err := tokeniser.Tokenise(expression)
	if err != nil {
		return evaluator.Number{}, fmt.Errorf("Tokeniser error: %v", err)
	}

	program, err := parser.ParseProgram(tokens)
	if err != nil {
		return evaluator.Number{}, fmt.Errorf("Parser error: %v", err)
	}

	eval, err := evaluator.EvaluateProgram(program, degreesMode)
	if err != nil {
		return evaluator.Number{}, fmt.Errorf("Evaluator error: %v", err)
	}

	return eval, nil
}

func CheckSyntax(expression string) error {
	tokens, err := tokeniser.Tokenise(expression)
	if err != nil {
		return fmt.Errorf("Tokeniser error: %v", err)
	}

	_, err = parser.ParseProgram(tokens)
	if err != nil {
		return fmt.Errorf("Parser error: %v", err)
	}

	return nil
}
