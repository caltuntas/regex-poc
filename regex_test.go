package main


import (
	"testing"
)

func TestMatchSequenceAB(t *testing.T) {
	l := New("ab")
	parser := NewParser(l)
	ast := parser.Ast()
	nfa := Compile(ast)

	cases := []struct {
		input string
		match bool
	}{
		{"ab", true},  
		{"a", false},  
		{"b", false},  
		{"", false},   
		{"abc", false},
		{"xab", false},
		{"abx", false},
	}

	for _, c := range cases {
		if got := Match(nfa, c.input); got != c.match {
			t.Errorf("Match(%q) = %v, want %v", c.input, got, c.match)
		}
	}
}
