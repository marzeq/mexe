package tokeniser

import (
	"fmt"
	"strings"
)

type Tokeniser struct {
	Source []rune
	Pos    int
	Tokens []Token
}

func (t *Tokeniser) Peek() rune {
	if t.Pos >= len(t.Source) {
		return 0
	}
	return t.Source[t.Pos]
}

func (t *Tokeniser) Consume() rune {
	if t.Pos >= len(t.Source) {
		return 0
	}
	ch := t.Source[t.Pos]
	t.Pos++
	return ch
}

func (t *Tokeniser) Skip() {
	t.Pos++
}

func (t *Tokeniser) ConsumeWhile(predicate func(rune) bool) string {
	sb := strings.Builder{}

	for {
		ch := t.Peek()
		if ch == 0 || !predicate(ch) {
			break
		}
		sb.WriteRune(ch)
		t.Skip()
	}
	return sb.String()
}

func (t *Tokeniser) SkipWhile(predicate func(rune) bool) {
	for {
		ch := t.Peek()
		if ch == 0 || !predicate(ch) {
			break
		}
		t.Skip()
	}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (t *Tokeniser) SkipWhitespace() {
	t.SkipWhile(isWhitespace)
}

func (t *Tokeniser) Tokenise() error {
	for {
		t.SkipWhitespace()

		startPos := t.Pos
		ch := t.Peek()
		if ch == 0 {
			t.Tokens = append(t.Tokens, EofToken)
			break
		}

		if isDigit(ch) {
			numberStr := t.ConsumeWhile(isDigit)

			if t.Peek() == '.' {
				t.Skip()

				decimalStr := t.ConsumeWhile(isDigit)
				if decimalStr == "" {
					decimalStr = "0"
				}

				t.Tokens = append(t.Tokens, Token{
					Type:  TOKEN_NUMBER,
					Pos:   startPos,
					Value: numberStr + "." + decimalStr,
				})
			} else {
				t.Tokens = append(t.Tokens, Token{
					Type:  TOKEN_NUMBER,
					Pos:   startPos,
					Value: numberStr,
				})
			}
		} else if isLetter(ch) {
			identStr := t.ConsumeWhile(func(c rune) bool {
				return isLetter(c) || isDigit(c) || c == '_'
			})

			if identStr == "mod" {
				t.Tokens = append(t.Tokens, Token{
					Type:  TOKEN_MOD,
					Pos:   startPos,
					Value: identStr,
				})
			} else {
				t.Tokens = append(t.Tokens, Token{
					Type:  TOKEN_IDENTIFIER,
					Pos:   startPos,
					Value: identStr,
				})
			}
		} else if ch == 'π' {
			t.Tokens = append(t.Tokens, Token{
				Type:  TOKEN_IDENTIFIER,
				Pos:   startPos,
				Value: "pi",
			})
			t.Skip()
		} else if ch == '.' {
			t.Skip()

			decimalStr := t.ConsumeWhile(isDigit)
			if decimalStr == "" {
				return fmt.Errorf("Expected digits after decimal point")
			}
			t.Tokens = append(t.Tokens, Token{
				Type:  TOKEN_NUMBER,
				Pos:   startPos,
				Value: "0." + decimalStr,
			})
		} else if ch == '%' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_PERCENT, Pos: startPos})
			t.Skip()
		} else if ch == '+' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_PLUS, Pos: startPos})
			t.Skip()
		} else if ch == '-' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_MINUS, Pos: startPos})
			t.Skip()
		} else if ch == '*' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_ASTERISK, Pos: startPos})
			t.Skip()
		} else if ch == '/' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_SLASH, Pos: startPos})
			t.Skip()
		} else if ch == '^' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_CARET, Pos: startPos})
			t.Skip()
		} else if ch == '!' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_EXCLAM, Pos: startPos})
			t.Skip()
		} else if ch == '(' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_LPAREN, Pos: startPos})
			t.Skip()
		} else if ch == ')' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_RPAREN, Pos: startPos})
			t.Skip()
		} else if ch == ',' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_COMMA, Pos: startPos})
			t.Skip()
		} else if ch == '|' {
			t.Tokens = append(t.Tokens, Token{Type: TOKEN_PIPE, Pos: startPos})
			t.Skip()
		} else {
			return fmt.Errorf("Unknown character: %s", string(ch))
		}
	}

	return nil
}

func Tokenise(source string) ([]Token, error) {
	tokeniser := Tokeniser{
		Source: []rune(source),
		Pos:    0,
	}
	err := tokeniser.Tokenise()
	return tokeniser.Tokens, err
}
