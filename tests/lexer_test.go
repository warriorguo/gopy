package tests

import (
	"testing"

	"github.com/warriorguo/gopy/pkg/lexer"
)

func TestLexerBasic(t *testing.T) {
	tests := []struct {
		input    string
		expected []lexer.TokenType
	}{
		{
			"x = 42",
			[]lexer.TokenType{lexer.IDENT, lexer.ASSIGN, lexer.INT, lexer.EOF},
		},
		{
			"print x + y",
			[]lexer.TokenType{lexer.PRINT, lexer.IDENT, lexer.PLUS, lexer.IDENT, lexer.EOF},
		},
		{
			"if x > 0:",
			[]lexer.TokenType{lexer.IF, lexer.IDENT, lexer.GT, lexer.INT, lexer.COLON, lexer.EOF},
		},
		{
			"for i in range(10):",
			[]lexer.TokenType{lexer.FOR, lexer.IDENT, lexer.IN, lexer.RANGE, lexer.LPAREN, lexer.INT, lexer.RPAREN, lexer.COLON, lexer.EOF},
		},
		{
			"for x in [1, 2, 3]:",
			[]lexer.TokenType{lexer.FOR, lexer.IDENT, lexer.IN, lexer.LBRACKET, lexer.INT, lexer.COMMA, lexer.INT, lexer.COMMA, lexer.INT, lexer.RBRACKET, lexer.COLON, lexer.EOF},
		},
		{
			"for item in items:",
			[]lexer.TokenType{lexer.FOR, lexer.IDENT, lexer.IN, lexer.IDENT, lexer.COLON, lexer.EOF},
		},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		tokens := l.AllTokens()

		if len(tokens) != len(test.expected) {
			t.Errorf("Expected %d tokens, got %d for input %q", len(test.expected), len(tokens), test.input)
			continue
		}

		for i, expected := range test.expected {
			if tokens[i].Type != expected {
				t.Errorf("Token %d: expected %s, got %s for input %q", i, expected, tokens[i].Type, test.input)
			}
		}
	}
}

func TestLexerIndentation(t *testing.T) {
	input := `if x > 0:
    print "positive"
    if x > 10:
        print "big"
    print "done"
print "end"`

	expected := []lexer.TokenType{
		lexer.IF, lexer.IDENT, lexer.GT, lexer.INT, lexer.COLON, lexer.NEWLINE,
		lexer.INDENT, lexer.PRINT, lexer.STRING, lexer.NEWLINE,
		lexer.IF, lexer.IDENT, lexer.GT, lexer.INT, lexer.COLON, lexer.NEWLINE,
		lexer.INDENT, lexer.PRINT, lexer.STRING, lexer.NEWLINE,
		lexer.DEDENT, lexer.PRINT, lexer.STRING, lexer.NEWLINE,
		lexer.DEDENT, lexer.PRINT, lexer.STRING, lexer.EOF,
	}

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()

	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
		for i, tok := range tokens {
			t.Logf("Token %d: %s", i, tok.Type)
		}
		return
	}

	for i, expected := range expected {
		if tokens[i].Type != expected {
			t.Errorf("Token %d: expected %s, got %s", i, expected, tokens[i].Type)
		}
	}
}

func TestLexerStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`"hello\nworld"`, "hello\nworld"},
		{`'it\'s'`, "it's"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		token := l.NextToken()

		if token.Type != lexer.STRING {
			t.Errorf("Expected STRING token, got %s", token.Type)
			continue
		}

		if token.Lexeme != test.expected {
			t.Errorf("Expected string %q, got %q", test.expected, token.Lexeme)
		}
	}
}

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		input        string
		expectedType lexer.TokenType
		expectedLex  string
	}{
		{"42", lexer.INT, "42"},
		{"3.14", lexer.FLOAT, "3.14"},
		{"0", lexer.INT, "0"},
		{"123.456", lexer.FLOAT, "123.456"},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		token := l.NextToken()

		if token.Type != test.expectedType {
			t.Errorf("Expected %s token, got %s", test.expectedType, token.Type)
			continue
		}

		if token.Lexeme != test.expectedLex {
			t.Errorf("Expected lexeme %q, got %q", test.expectedLex, token.Lexeme)
		}
	}
}

func TestLexerForKeyword(t *testing.T) {
	tests := []struct {
		input    string
		expected []struct {
			tokenType lexer.TokenType
			lexeme    string
		}
	}{
		{
			"for",
			[]struct {
				tokenType lexer.TokenType
				lexeme    string
			}{
				{lexer.FOR, "for"},
				{lexer.EOF, ""},
			},
		},
		{
			"for i in range(5):",
			[]struct {
				tokenType lexer.TokenType
				lexeme    string
			}{
				{lexer.FOR, "for"},
				{lexer.IDENT, "i"},
				{lexer.IN, "in"},
				{lexer.RANGE, "range"},
				{lexer.LPAREN, "("},
				{lexer.INT, "5"},
				{lexer.RPAREN, ")"},
				{lexer.COLON, ":"},
				{lexer.EOF, ""},
			},
		},
		{
			"for item in items:\n    print item",
			[]struct {
				tokenType lexer.TokenType
				lexeme    string
			}{
				{lexer.FOR, "for"},
				{lexer.IDENT, "item"},
				{lexer.IN, "in"},
				{lexer.IDENT, "items"},
				{lexer.COLON, ":"},
				{lexer.NEWLINE, "\n"},
				{lexer.INDENT, ""},
				{lexer.PRINT, "print"},
				{lexer.IDENT, "item"},
				{lexer.DEDENT, ""}, // EOF triggers DEDENT for unclosed indentation
				{lexer.EOF, ""},
			},
		},
		{
			"for x in [1, 2, 3]:",
			[]struct {
				tokenType lexer.TokenType
				lexeme    string
			}{
				{lexer.FOR, "for"},
				{lexer.IDENT, "x"},
				{lexer.IN, "in"},
				{lexer.LBRACKET, "["},
				{lexer.INT, "1"},
				{lexer.COMMA, ","},
				{lexer.INT, "2"},
				{lexer.COMMA, ","},
				{lexer.INT, "3"},
				{lexer.RBRACKET, "]"},
				{lexer.COLON, ":"},
				{lexer.EOF, ""},
			},
		},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		tokens := l.AllTokens()

		if len(tokens) != len(test.expected) {
			t.Errorf("For input %q: expected %d tokens, got %d", test.input, len(test.expected), len(tokens))
			for i, tok := range tokens {
				t.Logf("Token %d: %s (%q)", i, tok.Type, tok.Lexeme)
			}
			continue
		}

		for i, expected := range test.expected {
			if i >= len(tokens) {
				break
			}
			
			if tokens[i].Type != expected.tokenType {
				t.Errorf("For input %q, token %d: expected type %s, got %s", 
					test.input, i, expected.tokenType, tokens[i].Type)
			}
			
			if expected.lexeme != "" && tokens[i].Lexeme != expected.lexeme {
				t.Errorf("For input %q, token %d: expected lexeme %q, got %q", 
					test.input, i, expected.lexeme, tokens[i].Lexeme)
			}
		}
	}
}