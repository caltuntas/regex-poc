package main

import (
	"fmt"
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




func TestMatchRepetitionABStar(t *testing.T) {
	l := New("ab*")
	parser := NewParser(l)
	ast := parser.Ast()
	nfa := Compile(ast)

	str := nfa.ToString()
	fmt.Println("Ast = " +ast.String())
	fmt.Println("Nfa")
	fmt.Println(str)

	cases := []struct {
		input string
		match bool
	}{
		{"a", true},       
		{"ab", true},      
		{"abb", true},     
		{"abbb", true},    
		{"b", false},      
		{"ba", false},     
		{"", false},       
		{"abbbbcd", false},
		{"xab", false},    
	}

	for _, c := range cases {
		if got := Match(nfa, c.input); got != c.match {
			t.Errorf("Match(%q) = %v, want %v", c.input, got, c.match)
		}
	}
}

func TestMatchComplexPatternABStarCDStar(t *testing.T) {
	l := New("ab*cd*")
	parser := NewParser(l)
	ast := parser.Ast()
	nfa := Compile(ast)

	cases := []struct {
		input string
		match bool
	}{
		{"ac", true},           
		{"abc", true},          
		{"abbbbbc", true},      
		{"acd", true},          
		{"abbbbbcdddd", true},  
		{"a", false},           
		{"ab", false},          
		{"abbbbb", false},      
		{"xabc", false},        
		{"abccd", false},       
	}

	for _, c := range cases {
		if got := Match(nfa, c.input); got != c.match {
			t.Errorf("Match(%q) = %v, want %v", c.input, got, c.match)
		}
	}
}

func TestMatchSimpleDotStar(t *testing.T) {
	l := New("aa.*")
	parser := NewParser(l)
	ast := parser.Ast()
	nfa := Compile(ast)

	cases := []struct {
		input string
		match bool
	}{
		{"aac", true},           
		{"aab", true},          
		{"aabbcc", true},      
		{"aa", true},          
		{"aaabbb", true},  
		{"a", false},           
		{"ab", false},          
		{"abbbbb", false},      
		{"xabc", false},        
		{"abccd", false},       
	}

	for _, c := range cases {
		if got := Match(nfa, c.input); got != c.match {
			t.Errorf("Match(%q) = %v, want %v", c.input, got, c.match)
		}
	}
}
