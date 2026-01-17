package parser

import (
	"fmt"

	"github.com/marzeq/mexe/tokeniser"
)

type Parser struct {
	Tokens []tokeniser.Token
	Pos    int
}

func (p *Parser) Peek() tokeniser.Token {
	if p.Pos >= len(p.Tokens) {
		return tokeniser.EofToken
	}
	return p.Tokens[p.Pos]
}

func (p *Parser) Consume() tokeniser.Token {
	if p.Pos >= len(p.Tokens) {
		return tokeniser.EofToken
	}
	tok := p.Tokens[p.Pos]
	p.Pos++
	return tok
}

func (p *Parser) Skip() {
	if p.Pos < len(p.Tokens) {
		p.Pos++
	}
}

func (p *Parser) Match(expectedType tokeniser.TokenType) bool {
	tok := p.Peek()
	if tok.Type == expectedType {
		p.Skip()
		return true
	}
	return false
}

func (p *Parser) Expect(expectedType tokeniser.TokenType) (tokeniser.Token, error) {
	tok := p.Consume()
	if tok.Type != expectedType {
		return tokeniser.EofToken, fmt.Errorf("Unexpected token: %s", tok.String())
	}
	return tok, nil
}

func (p *Parser) parseExpression() (Node, error) {
	return p.parseAddSub()
}

func (p *Parser) parseAddSub() (Node, error) {
	left, err := p.parseMulDivMod()
	if err != nil {
		return nil, err
	}

	for {
		switch p.Peek().Type {
		case tokeniser.TOKEN_PLUS, tokeniser.TOKEN_MINUS:
			opTok := p.Consume()
			right, err := p.parseMulDivMod()
			if err != nil {
				return nil, err
			}

			opType := BinaryOpTypeAdd
			if opTok.Type == tokeniser.TOKEN_MINUS {
				opType = BinaryOpTypeSubtract
			}

			left = &BinaryOpNode{
				OpType: opType,
				Left:   left,
				Right:  right,
			}
		default:
			return left, nil
		}
	}
}

func isImplicitMulStart(t tokeniser.TokenType) bool {
	switch t {
	case tokeniser.TOKEN_NUMBER,
		tokeniser.TOKEN_IDENTIFIER,
		tokeniser.TOKEN_LPAREN:
		return true
	}
	return false
}

func (p *Parser) parseMulDivMod() (Node, error) {
	left, err := p.parseExponent()
	if err != nil {
		return nil, err
	}

	for {
		switch p.Peek().Type {
		case tokeniser.TOKEN_ASTERISK, tokeniser.TOKEN_SLASH, tokeniser.TOKEN_MOD, tokeniser.TOKEN_PERCENT:
			opTok := p.Consume()
			right, err := p.parseExponent()
			if err != nil {
				return nil, err
			}

			var opType BinaryOpType
			switch opTok.Type {
			case tokeniser.TOKEN_ASTERISK:
				opType = BinaryOpTypeMultiply
			case tokeniser.TOKEN_SLASH:
				opType = BinaryOpTypeDivide
			case tokeniser.TOKEN_MOD, tokeniser.TOKEN_PERCENT:
				opType = BinaryOpTypeMod
			}

			left = &BinaryOpNode{
				OpType: opType,
				Left:   left,
				Right:  right,
			}

		default:
			if isImplicitMulStart(p.Peek().Type) {
				right, err := p.parseExponent()
				if err != nil {
					return nil, err
				}

				left = &BinaryOpNode{
					OpType: BinaryOpTypeMultiply,
					Left:   left,
					Right:  right,
				}
				continue
			}

			return left, nil
		}
	}
}

func (p *Parser) parseExponent() (Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	if p.Match(tokeniser.TOKEN_CARET) {
		right, err := p.parseExponent()
		if err != nil {
			return nil, err
		}

		return &BinaryOpNode{
			OpType: BinaryOpTypeExponent,
			Left:   left,
			Right:  right,
		}, nil
	}

	return left, nil
}

func (p *Parser) parseUnary() (Node, error) {
	var node Node
	var err error

	if p.Match(tokeniser.TOKEN_MINUS) {
		child, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		node = &UnaryOpNode{
			OpType: UnaryOpTypeNegate,
			Child:  child,
		}
	} else {
		node, err = p.parsePrimary()
		if err != nil {
			return nil, err
		}
	}

	for p.Match(tokeniser.TOKEN_EXCLAM) {
		node = &UnaryOpNode{
			OpType: UnaryOpTypeFact,
			Child:  node,
		}
	}

	return node, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	tok := p.Peek()

	switch tok.Type {
	case tokeniser.TOKEN_NUMBER:
		p.Consume()
		return &ConstantNode{Name: tok.Value}, nil

	case tokeniser.TOKEN_IDENTIFIER, tokeniser.TOKEN_MOD:
		p.Consume()

		name := tok.Value
		if tok.Type == tokeniser.TOKEN_MOD {
			name = "mod"
		}

		if p.Match(tokeniser.TOKEN_LPAREN) {
			var args []Node

			if !p.Match(tokeniser.TOKEN_RPAREN) {
				for {
					arg, err := p.parseExpression()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)

					if p.Match(tokeniser.TOKEN_COMMA) {
						if p.Match(tokeniser.TOKEN_RPAREN) {
							break
						}
						continue
					}

					if _, err := p.Expect(tokeniser.TOKEN_RPAREN); err != nil {
						return nil, err
					}
					break
				}
			}

			return &FuncApplicationNode{
				FuncName: name,
				Args:     args,
			}, nil
		}

		return &ConstantNode{Name: name}, nil

	case tokeniser.TOKEN_LPAREN:
		p.Consume()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.Expect(tokeniser.TOKEN_RPAREN); err != nil {
			return nil, err
		}
		return expr, nil

	case tokeniser.TOKEN_PIPE:
		p.Consume()

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if _, err := p.Expect(tokeniser.TOKEN_PIPE); err != nil {
			return nil, err
		}

		return &UnaryOpNode{
			OpType: UnaryOpTypeAbs,
			Child:  expr,
		}, nil
	}

	return nil, fmt.Errorf("Unexpected token: %s", tok.String())
}

func Parse(tokens []tokeniser.Token) (Node, error) {
	parser := Parser{
		Tokens: tokens,
		Pos:    0,
	}

	parsed, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}

	if parser.Peek().Type != tokeniser.TOKEN_EOF {
		return nil, fmt.Errorf("Unexpected token: %s", parser.Peek().String())
	}

	return parsed, nil
}
