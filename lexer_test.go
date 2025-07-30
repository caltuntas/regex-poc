package main

import "testing"

func TestNextToken(t *testing.T) {
	l := New("pa.*tpa[ab]cpa\\s")
	tests := []struct {
		expectedType TokenType
		expectedLiteral string
	} {
		{LITERAL, "p"},
		{LITERAL, "a"},
		{DOT, "."},
		{STAR, "*"},
		{LITERAL, "t"},
		{LITERAL, "p"},
		{LITERAL, "a"},
		{LBRACKET, "["},
		{LITERAL, "a"},
		{LITERAL, "b"},
		{RBRACKET, "]"},
		{LITERAL, "c"},
		{LITERAL, "p"},
		{LITERAL, "a"},
		{ESCAPE, "\\"},
		{LITERAL, "s"},
	}

	for i, test := range tests {
		token := l.NextToken()
		if token.Type != test.expectedType {
			t.Fatalf("test[%d], expected type doesn't match. expected = %q, got = %q",i, token.Type, test.expectedType)
		}
		if token.Value != test.expectedLiteral {
			t.Fatalf("test[%d], expected type doesn't match. expected = %q, got = %q",i, token.Value, test.expectedLiteral)
		}
	}
}
