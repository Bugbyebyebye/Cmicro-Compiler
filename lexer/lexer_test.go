package lexer

import (
	"Cmicro-Compiler/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `"foobar"
				"foo bar" `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		l := New(input)
		token := l.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokeType wrong. expected=%q, got=%q",
				i, tt.expectedType, token.Type)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, token.Literal)
		}
	}
}
