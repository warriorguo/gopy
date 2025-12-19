package lexer

import (
	"unicode"
)

type Lexer struct {
	src      []rune
	position int
	line     int
	column   int

	indentStack    []int
	atLineStart    bool
	pendingDedents int
}

func NewLexer(src string) *Lexer {
	runes := []rune(src)
	return &Lexer{
		src:         runes,
		position:    0,
		line:        1,
		column:      1,
		indentStack: []int{0},
		atLineStart: true,
	}
}

func (l *Lexer) peekChar() rune {
	if l.position >= len(l.src) {
		return 0
	}
	return l.src[l.position]
}

func (l *Lexer) readChar() rune {
	if l.position >= len(l.src) {
		return 0
	}
	ch := l.src[l.position]
	l.position++
	if ch == '\n' {
		l.line++
		l.column = 1
		l.atLineStart = true
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for {
		ch := l.peekChar()
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.readChar()
		} else {
			break
		}
	}
}

func (l *Lexer) readString(quote rune) string {
	var result []rune

	for {
		ch := l.peekChar()
		if ch == 0 || ch == quote {
			break
		}
		if ch == '\\' {
			l.readChar()
			next := l.readChar()
			switch next {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case '\'':
				result = append(result, '\'')
			case '"':
				result = append(result, '"')
			default:
				result = append(result, next)
			}
		} else {
			result = append(result, l.readChar())
		}
	}

	if l.peekChar() == quote {
		l.readChar()
	}

	return string(result)
}

func (l *Lexer) readNumber() (string, TokenType) {
	start := l.position
	hasDecimal := false

	for {
		ch := l.peekChar()
		if unicode.IsDigit(ch) {
			l.readChar()
		} else if ch == '.' && !hasDecimal {
			hasDecimal = true
			l.readChar()
		} else {
			break
		}
	}

	tokenType := INT
	if hasDecimal {
		tokenType = FLOAT
	}

	return string(l.src[start:l.position]), tokenType
}

func (l *Lexer) readIdentifier() string {
	start := l.position

	ch := l.peekChar()
	if !unicode.IsLetter(ch) && ch != '_' {
		return ""
	}

	for {
		ch := l.peekChar()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			l.readChar()
		} else {
			break
		}
	}

	return string(l.src[start:l.position])
}

func (l *Lexer) skipComment() {
	for {
		ch := l.peekChar()
		if ch == 0 || ch == '\n' {
			break
		}
		l.readChar()
	}
}

func (l *Lexer) handleIndentation() *Token {
	if l.pendingDedents > 0 {
		l.pendingDedents--
		return &Token{
			Type:   DEDENT,
			Lexeme: "",
			Line:   l.line,
			Column: l.column,
		}
	}

	if !l.atLineStart {
		return nil
	}

	l.atLineStart = false
	start := l.column
	indentLevel := 0

	for {
		ch := l.peekChar()
		if ch == ' ' {
			indentLevel++
			l.readChar()
		} else if ch == '\t' {
			indentLevel += 8
			l.readChar()
		} else {
			break
		}
	}

	if l.peekChar() == '\n' || l.peekChar() == '#' || l.peekChar() == 0 {
		return nil
	}

	currentIndent := l.indentStack[len(l.indentStack)-1]

	if indentLevel > currentIndent {
		l.indentStack = append(l.indentStack, indentLevel)
		return &Token{
			Type:   INDENT,
			Lexeme: "",
			Line:   l.line,
			Column: start,
		}
	} else if indentLevel < currentIndent {
		dedentCount := 0
		for len(l.indentStack) > 1 && l.indentStack[len(l.indentStack)-1] > indentLevel {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			dedentCount++
		}

		if len(l.indentStack) == 0 || l.indentStack[len(l.indentStack)-1] != indentLevel {
			return &Token{
				Type:   ILLEGAL,
				Lexeme: "indentation error",
				Line:   l.line,
				Column: start,
			}
		}

		l.pendingDedents = dedentCount - 1
		return &Token{
			Type:   DEDENT,
			Lexeme: "",
			Line:   l.line,
			Column: start,
		}
	}

	return nil
}

func (l *Lexer) NextToken() Token {
	if token := l.handleIndentation(); token != nil {
		return *token
	}

	l.skipWhitespace()

	line, column := l.line, l.column
	ch := l.readChar()

	switch ch {
	case 0:
		// Generate all necessary DEDENTs before EOF
		if len(l.indentStack) > 1 {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return Token{
				Type:   DEDENT,
				Lexeme: "",
				Line:   line,
				Column: column,
			}
		}
		return Token{Type: EOF, Lexeme: "", Line: line, Column: column}
	case '\n':
		return Token{Type: NEWLINE, Lexeme: "\n", Line: line, Column: column}
	case '#':
		l.skipComment()
		return l.NextToken()
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: EQ, Lexeme: "==", Line: line, Column: column}
		}
		return Token{Type: ASSIGN, Lexeme: "=", Line: line, Column: column}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: NOT_EQ, Lexeme: "!=", Line: line, Column: column}
		}
		return Token{Type: ILLEGAL, Lexeme: "!", Line: line, Column: column}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: LTE, Lexeme: "<=", Line: line, Column: column}
		}
		return Token{Type: LT, Lexeme: "<", Line: line, Column: column}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: GTE, Lexeme: ">=", Line: line, Column: column}
		}
		return Token{Type: GT, Lexeme: ">", Line: line, Column: column}
	case '+':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: PLUS_ASSIGN, Lexeme: "+=", Line: line, Column: column}
		}
		return Token{Type: PLUS, Lexeme: "+", Line: line, Column: column}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: MINUS_ASSIGN, Lexeme: "-=", Line: line, Column: column}
		}
		return Token{Type: MINUS, Lexeme: "-", Line: line, Column: column}
	case '*':
		return Token{Type: MULTIPLY, Lexeme: "*", Line: line, Column: column}
	case '/':
		return Token{Type: DIVIDE, Lexeme: "/", Line: line, Column: column}
	case '%':
		return Token{Type: MODULO, Lexeme: "%", Line: line, Column: column}
	case ',':
		return Token{Type: COMMA, Lexeme: ",", Line: line, Column: column}
	case ':':
		return Token{Type: COLON, Lexeme: ":", Line: line, Column: column}
	case ';':
		return Token{Type: SEMICOLON, Lexeme: ";", Line: line, Column: column}
	case '(':
		return Token{Type: LPAREN, Lexeme: "(", Line: line, Column: column}
	case ')':
		return Token{Type: RPAREN, Lexeme: ")", Line: line, Column: column}
	case '[':
		return Token{Type: LBRACKET, Lexeme: "[", Line: line, Column: column}
	case ']':
		return Token{Type: RBRACKET, Lexeme: "]", Line: line, Column: column}
	case '{':
		return Token{Type: LBRACE, Lexeme: "{", Line: line, Column: column}
	case '}':
		return Token{Type: RBRACE, Lexeme: "}", Line: line, Column: column}
	case '"', '\'':
		lexeme := l.readString(ch)
		return Token{Type: STRING, Lexeme: lexeme, Line: line, Column: column}
	default:
		if unicode.IsDigit(ch) {
			l.position--
			l.column--
			lexeme, tokenType := l.readNumber()
			return Token{Type: tokenType, Lexeme: lexeme, Line: line, Column: column}
		} else if unicode.IsLetter(ch) || ch == '_' {
			l.position--
			l.column--
			lexeme := l.readIdentifier()
			tokenType := LookupIdent(lexeme)
			return Token{Type: tokenType, Lexeme: lexeme, Line: line, Column: column}
		}
		return Token{Type: ILLEGAL, Lexeme: string(ch), Line: line, Column: column}
	}
}

func (l *Lexer) AllTokens() []Token {
	var tokens []Token
	for {
		token := l.NextToken()
		tokens = append(tokens, token)
		if token.Type == EOF {
			break
		}
	}
	return tokens
}
