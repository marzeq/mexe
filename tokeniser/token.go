package tokeniser

type TokenType int

const (
	TOKEN_EOF TokenType = iota

	TOKEN_NUMBER

	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_ASTERISK // *
	TOKEN_SLASH    // /
	TOKEN_CARET    // ^
	TOKEN_MOD      // "mod"
	TOKEN_PERCENT  // %
	TOKEN_EXCLAM   // !

	TOKEN_LPAREN // (
	TOKEN_RPAREN // (
	TOKEN_COMMA  // ,

	TOKEN_PIPE // |

	TOKEN_IDENTIFIER
)

type Token struct {
	Type  TokenType
	Pos   int
	Value string
}

func (t Token) String() string {
	switch t.Type {
	case TOKEN_EOF:
		return "(EOF)"
	case TOKEN_NUMBER:
		return t.Value
	case TOKEN_MINUS:
		return "-"
	case TOKEN_PLUS:
		return "+"
	case TOKEN_ASTERISK:
		return "*"
	case TOKEN_SLASH:
		return "/"
	case TOKEN_CARET:
		return "^"
	case TOKEN_MOD:
		return "mod"
	case TOKEN_PERCENT:
		return "%"
	case TOKEN_EXCLAM:
		return "!"
	case TOKEN_LPAREN:
		return "("
	case TOKEN_RPAREN:
		return ")"
	case TOKEN_COMMA:
		return ","
	case TOKEN_PIPE:
		return "|"
	case TOKEN_IDENTIFIER:
		return t.Value
	default:
		panic("Unknown token type")
	}
}

var EofToken = Token{Type: TOKEN_EOF}
