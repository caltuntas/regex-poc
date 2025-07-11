package main

import "testing"

func TestNextToken(t *testing.T) {
	l := New("pa.*t")
	tests := []struct {
		expectedType TokenType
		expectedLiteral string
	} {
		{LITERAL, "p"},
		{LITERAL, "a"},
		{DOT, "."},
		{STAR, "*"},
		{LITERAL, "t"},
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

func TestCharacterClasses(t *testing.T) {
	l := New("pa[ab]c")
	tests := []struct {
		expectedType TokenType
		expectedLiteral string
	} {
		{LITERAL, "p"},
		{LITERAL, "a"},
		{LBRACKET, "["},
		{LITERAL, "a"},
		{LITERAL, "b"},
		{RBRACKET, "]"},
		{LITERAL, "c"},
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
