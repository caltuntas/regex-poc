package main

import (
	"fmt"
	"testing"
)

func TestRegexPerformanceBacktrack(t *testing.T) {
	t.Skip()
	pattern := `parent {[\s\S]*type.*[\s\S]*subtype.*[\s\S]*element.*[\s\S]*attributes.*[\s\S]*value testvalue.*[\s\S]*owner person1`
	l := New(pattern)
	parser := NewParser(l)
	ast := parser.Ast()
	nfa := Compile(ast)
	str := nfa.ToDigraph()
	fmt.Println(str)
	if got := MatchBacktrack(ast, config); got != true {
		t.Errorf("Match(%q) = %v, want %v", config, got, false)
	}
}

func TestRegexPerformancePartialBacktrack(t *testing.T) {
	t.Skip()
	pattern := `parent {[\s\S]*type.*[\s\S]*subtype.*[\s\S]*element.*[\s\S]*attributes.*[\s\S]*value testvalue.*[\s\S]*owner person1`
	l := New(pattern)
	parser := NewParser(l)
	ast := parser.Ast()
	nfa := Compile(ast)
	str := nfa.ToDigraph()
	fmt.Println(str)
	if got := MatchBacktrackPartial(ast, allConfig); got != true {
		t.Errorf("Match(%q) = %v, want %v", config, got, false)
	}
}

func TestRegexMatchBacktrack(t *testing.T) {
	cases := map[string][]struct {
		input string
		match bool
	}{
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
			{"abcd", true},
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
		"a.*b*": {
			{"a", true},
			{"ab", true},
			{"a123", true},
			{"axxxbbbb", true},
			{"aBBBB", true},
			{"a.b", true},
			{"abbbbbb", true},
			{"b", false},
			{"x", false},
			{"", false},
		},
		"a[bc]d": {
			{"abd", true},
			{"acd", true},
			{"aad", false},
			{"abcd", false},
			{"abbd", false},
			{"ad", false},
			{"abd ", false},
			{" abd", false},
			{"Abd", false},
			{"aCd", false},
			{"", false},
		},
		"pa\\sb": {
			{"pa b", true},
			{"pa\tb", true},
			{"pa\nb", true},
			{"pa\r\nb", false},
			{"pa\v b", false},
			{"pa\vb", true},
			{"pa\fb", true},
			{"pab", false},
			{"pa  b", false},
			{" pa b", false},
			{"pa b ", false},
			{"p a b", false},
			{"", false},
		},
		"a[bc]*d": {
			{"ad", true},
			{"abd", true},
			{"acd", true},
			{"abcd", true},
			{"abcbcd", true},
			{"abccbd", true},
			{"a", false},
			{"d", false},
			{"abxd", false},
			{"axcd", false},
			{"ab cd", false},
			{"abcbcbcbcd", true},
			{"abcbdx", false},
			{"abccd", true},
			{"", false},
		},
		"a[\\sb]*d": {
			{"ad", true},
			{"abd", true},
			{"abbd", true},
			{"a d", true},
			{"a\tbd", true},
			{"a \t\nbd", true},
			{"a \tb b\t\n\rbd", true},
			{"abxd", false},
			{"abcd", false},
			{"axd", false},
			{"a\n\n\n\n\nd", true},
			{"a\rbd", true},
			{"a\vd", true},
			{"a\fd", true},
			{"", false},
			{"a    d", true},
			{"abd ", false},
			{" ab d", false},
		},
		"pa[\\s\\S]*b": {
			{"pab", true},
			{"pa123b", true},
			{"pa b", true},
			{"pa\tb", true},
			{"pa\nb", true},
			{"pa something b", true},
			{"pa---b", true},
			{"pa\nmulti\nline\nb", true},
			{"pabbbbb", true},
			{"pa", false},
			{"p", false},
			{"pb", false},
			{"ab", false},
			{"", false},
			{"pa middle x", false},
		},
		"pa {": {
			{"pa {", true},
			{"pa  {", false},
			{"pa{", false},
			{"pa  {", false},
			{"p a {", false},
			{"pa\t{", false},
			{"pa", false},
			{"", false},
			{"{ pa", false},
		},
		"pa.*b": {
			{"pab", true},
			{"paxb", true},
			{"paxyzb", true},
			{"pabbb", true},
			{"paX", false},
			{"pabbbbbX", false},
		},
		"a.*b": {
			{"ab", true},
			{"acb", true},
			{"acccb", true},
			{"abbbb", true},
			{"acccx", false},
			{"a", false},
		},
		"a.*b.*c": {
			{"abc", true},
			{"axybzzc", true},
			{"abbbc", true},
			{"abbbx", false},
			{"ac", false},
			{"ab", false},
		},
		"a*b*": {
			{"", true},
			{"a", true},
			{"b", true},
			{"ab", true},
			{"aab", true},
			{"abb", true},
			{"aaabb", true},
			{"c", false},
		},
		".*b": {
			{"b", true},
			{"bb", true},
			{"ab", true},
			{"aab", true},
			{"abc", false},
			{"", false},
		},
		"a.*z": {
			{"az", true},
			{"axyz", true},
			{"axxxxxz", true},
			{"axxxxx", false},
		},
		"a.*a.*a.*b": {
			{"abbbbb", false},
			{"aaab", true},
			{"aab", false},
			{"ab", false},
			{"aaaa", false},
			{"aaaab", true},
			{"aaaaaaab", true},
		},
		"a*a*a*a*b": {
			{"aaaaaaaaaaaaaa", false},
		},
	}

	for key, val := range cases {
		l := New(key)
		parser := NewParser(l)
		ast := parser.Ast()
		for _, c := range val {
			fmt.Printf("match pattern=%s, input=%s\n\n", key, c.input)
			got := MatchBacktrack(ast, c.input)
			fmt.Printf("result pattern=%s, input=%s, result=%t\n", key, c.input, got)
			if got != c.match {
				t.Errorf("Pattern = %s, Match(%q) = %v, want %v", key, c.input, got, c.match)
			}
		}
	}
}

func TestBacktrackingRegexMatchPartial(t *testing.T) {
	cases := map[string][]struct {
		input string
		match bool
	}{
		"aa.*": {
			{"aa", true},
			{"aab", true},
			{"aa123", true},
			{"xxxy aab ttt", true},
			{"baaxy", true},
			{"123aa", true},
			{"aa bb", true},
			{"AAb", false},
			{"xaax", true},
			{"xaxa", false},
			{"aaa", true},
			{"baa", true},
			{"", false},
		},
	}

	for key, val := range cases {
		l := New(key)
		parser := NewParser(l)
		ast := parser.Ast()
		for _, c := range val {
			if got := MatchBacktrackPartial(ast, c.input); got != c.match {
				t.Errorf("Pattern = %s, Match(%q) = %v, want %v", key, c.input, got, c.match)
			}
		}
	}
}
