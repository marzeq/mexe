package mexe

import (
	"fmt"

	"github.com/marzeq/mexe/evaluator"
	"github.com/marzeq/mexe/parser"
	"github.com/marzeq/mexe/tokeniser"
)

func Evaluate(expression string, degreesMode bool) (float64, error) {
	tokens, err := tokeniser.Tokenise(expression)
	if err != nil {
		return 0, fmt.Errorf("Tokeniser error: %v", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return 0, fmt.Errorf("Parser error: %v", err)
	}

	eval, err := evaluator.Evaluate(ast, degreesMode)
	if err != nil {
		return 0, fmt.Errorf("Evaluator error: %v", err)
	}

	return eval, nil
}

func CheckSyntax(expression string) error {
	tokens, err := tokeniser.Tokenise(expression)
	if err != nil {
		return fmt.Errorf("Tokeniser error: %v", err)
	}

	_, err = parser.Parse(tokens)
	if err != nil {
		return fmt.Errorf("Parser error: %v", err)
	}

	return nil
}
