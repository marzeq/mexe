package parser

import (
	"fmt"

	"github.com/marzeq/mexe/tokeniser"
)

type Definition struct {
	Name   string
	Params []string // nil for a constant; non-nil for a function
	Value  Node
}

type Program struct {
	Definitions []Definition
	Expression  Node
}

// ParseProgram parses nonblank statements separated by semicolons or newlines.
// Every statement before the last must be a definition.
func ParseProgram(tokens []tokeniser.Token) (*Program, error) {
	var statements [][]tokeniser.Token
	var current []tokeniser.Token
	for _, tok := range tokens {
		if tok.Type == tokeniser.TOKEN_SEPARATOR || tok.Type == tokeniser.TOKEN_EOF {
			if len(current) != 0 {
				statements = append(statements, current)
				current = nil
			}
			continue
		}
		current = append(current, tok)
	}
	if len(statements) == 0 {
		return nil, fmt.Errorf("expected an expression")
	}

	program := &Program{}
	for i, statement := range statements {
		statement = append(statement, tokeniser.EofToken)
		if i == len(statements)-1 {
			expr, err := Parse(statement)
			if err != nil {
				return nil, err
			}
			program.Expression = expr
			break
		}
		def, err := parseDefinition(statement)
		if err != nil {
			return nil, fmt.Errorf("definition %d: %w", i+1, err)
		}
		program.Definitions = append(program.Definitions, def)
	}
	return program, nil
}

func parseDefinition(tokens []tokeniser.Token) (Definition, error) {
	p := Parser{Tokens: tokens}
	name, err := p.Expect(tokeniser.TOKEN_IDENTIFIER)
	if err != nil {
		return Definition{}, err
	}
	def := Definition{Name: name.Value}
	if p.Match(tokeniser.TOKEN_LPAREN) {
		def.Params = []string{}
		seen := map[string]bool{}
		if !p.Match(tokeniser.TOKEN_RPAREN) {
			for {
				param, err := p.Expect(tokeniser.TOKEN_IDENTIFIER)
				if err != nil {
					return Definition{}, err
				}
				if seen[param.Value] {
					return Definition{}, fmt.Errorf("duplicate parameter: %s", param.Value)
				}
				seen[param.Value] = true
				def.Params = append(def.Params, param.Value)
				if !p.Match(tokeniser.TOKEN_COMMA) {
					break
				}
			}
			if _, err := p.Expect(tokeniser.TOKEN_RPAREN); err != nil {
				return Definition{}, err
			}
		}
	}
	if _, err := p.Expect(tokeniser.TOKEN_EQUALS); err != nil {
		return Definition{}, err
	}
	def.Value, err = p.parseExpression()
	if err != nil {
		return Definition{}, err
	}
	if p.Peek().Type != tokeniser.TOKEN_EOF {
		return Definition{}, fmt.Errorf("Unexpected token: %s", p.Peek())
	}
	return def, nil
}
