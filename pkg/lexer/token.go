package lexer

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	INT
	FLOAT
	STRING

	ASSIGN
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULO
	PLUS_ASSIGN
	MINUS_ASSIGN

	EQ
	NOT_EQ
	LT
	GT
	LTE
	GTE

	AND
	OR
	NOT

	COMMA
	COLON
	SEMICOLON
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE

	NEWLINE
	INDENT
	DEDENT

	IF
	ELIF
	ELSE
	WHILE
	FOR
	IN
	DEF
	RETURN
	PRINT
	TRUE
	FALSE
	NONE
	RANGE
	PASS
)

var keywords = map[string]TokenType{
	"if":     IF,
	"elif":   ELIF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"in":     IN,
	"def":    DEF,
	"return": RETURN,
	"print":  PRINT,
	"True":   TRUE,
	"False":  FALSE,
	"None":   NONE,
	"and":    AND,
	"or":     OR,
	"not":    NOT,
	"range":  RANGE,
	"pass":   PASS,
}

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q, %d:%d}", t.Type, t.Lexeme, t.Line, t.Column)
}

func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case MODULO:
		return "MODULO"
	case PLUS_ASSIGN:
		return "PLUS_ASSIGN"
	case MINUS_ASSIGN:
		return "MINUS_ASSIGN"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case SEMICOLON:
		return "SEMICOLON"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case NEWLINE:
		return "NEWLINE"
	case INDENT:
		return "INDENT"
	case DEDENT:
		return "DEDENT"
	case IF:
		return "IF"
	case ELIF:
		return "ELIF"
	case ELSE:
		return "ELSE"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case IN:
		return "IN"
	case DEF:
		return "DEF"
	case RETURN:
		return "RETURN"
	case PRINT:
		return "PRINT"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NONE:
		return "NONE"
	case RANGE:
		return "RANGE"
	case PASS:
		return "PASS"
	default:
		return "UNKNOWN"
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}