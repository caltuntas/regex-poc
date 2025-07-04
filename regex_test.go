package main

import (
	"testing"
)

func TestRegexMatch(t *testing.T) {
	cases := map[string][]struct {
		input string
		match bool
	} {
	"aa.*": {
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
	},
	"ab*cd*": {
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
	},
	"ab*": {
		{"a", true},       
		{"ab", true},      
		{"abb", true},     
		{"abbb", true},    
		{"b", false},      
		{"ba", false},     
		{"", false},       
		{"abbbbcd", false},
		{"xab", false},    
	},
	"ab": {
		{"ab", true},  
		{"a", false},  
		{"b", false},  
		{"", false},   
		{"abc", false},
		{"xab", false},
		{"abx", false},
	},
}

	for key, val := range cases {
		l := New(key)
		parser := NewParser(l)
		ast := parser.Ast()
		nfa := Compile(ast)
		for _, c := range val {
			if got := Match(nfa, c.input); got != c.match {
				t.Errorf("Match(%q) = %v, want %v", c.input, got, c.match)
			}
		}
	}
}
