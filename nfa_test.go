package main

import (
	"strings"
	"testing"
)

var nb NodeBuilder

func CreateNfaFromString(str string) Nfa {
	var stateMap map[string]*State
	stateMap =make(map[string]*State)
	nfa := Nfa{}
	lines := strings.Split(strings.TrimSpace(str), "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ",")
		fromState := parts[0]
		toState := parts[1]
		transition := parts[2]
		if nfa.Start == nil {
			start := nfa.NewStart()
			stateMap[fromState] = start
			to := State{}
			stateMap[toState] = &to
			if transition == "ε" {
				start.AddEpsilonTo(&to)
			} else {
				start.AddTransition(getTransitionType(transition), transition, &to)
			}
		} else {
			from := stateMap[fromState]
			if from == nil {
				state := nfa.NewState()
				stateMap[fromState] = state
				from = state
			}
			to := stateMap[toState]
			if to == nil {
				state := nfa.NewState()
			 stateMap[toState] = state
				to = state
			}
			if transition == "ε" {
				from.AddEpsilonTo(to)
			} else {
				from.AddTransition(getTransitionType(transition), transition, to)
			}
		}
	}
	return nfa
}

func TestNFAStates(t *testing.T) {
	cases := []struct {
		node     Node
		expected string
	}{
		{
			nb.Lit('p'),
			`state1,state2,p`,
		},
		{
			nb.Seq(
				nb.Lit('p'),
				nb.Lit('a'),
			),
			`state1,state2,p
			 state2,state3,a`,
		},
		{
			nb.Star(nb.Lit('a')),
			`state1,state2,ε
			 state2,state3,a
			 state3,state2,ε
			 state3,state4,ε
			 state1,state4,ε`,
		},
		{
			nb.Seq(
				nb.Lit('a'),
				nb.Star(nb.Lit('b')),
			),
			`state1,state2,a
		   state2,state3,ε
			 state3,state4,b
			 state4,state5,ε
			 state4,state3,ε
			 state2,state5,ε`,
		},
		{
			nb.Seq(
				nb.Lit('a'),
				nb.Star(nb.Meta(DOT)),
			),
			`state1,state2,a
			 state2,state3,ε
			 state3,state4,.
			 state4,state5,ε
			 state4,state3,ε
			 state2,state5,ε`,
		},
		{
			nb.Seq(
				nb.Lit('a'),
				nb.Meta(WHITESPACE),
			),
			`state1,state2,a
			 state2,state3,\s`,
		},
		{
			nb.Seq(
				nb.Lit('p'),
				nb.List(nb.Lit('a'), nb.Lit('b')),
			),

			`state1,state2,p
			 state2,state3,ε
			 state3,state4,a
			 state4,state8,ε
			 state2,state5,ε
			 state5,state6,b
			 state6,state8,ε`,
		 },
		{
			nb.Meta(DOT),
			`state1,state2,.`,
		},
		{
			nb.Seq(
				nb.Lit('a'),
				nb.Star(nb.List(nb.Lit('b'), nb.Lit('c'))),
				nb.Lit('d'),
			),
			`state1,state2,a
			 state2,state3,ε
			 state3,state4,ε
			 state4,state5,b
			 state5,state6,ε
			 state6,state7,ε
			 state7,state8,d
			 state6,state3,ε
			 state3,state9,ε
			 state9,state10,c
			 state10,state6,ε
			 state2,state7,ε`,
		 },
	}

	for _, tc := range cases {
		nfa := Compile(tc.node)
		actualEncoding := nfa.Encode()
		expectedNfa := CreateNfaFromString(tc.expected)
		expectedEncoding := expectedNfa.Encode()
		if actualEncoding != expectedEncoding {
			t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
		}
	}
}

func TestEncode(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('p'),
		nb.List(nb.Lit('a'), nb.Lit('b')),
	)
	nfa := Compile(ast)
	actualEncoding := nfa.Encode()
	expectedEncoding := "(s-[Literal:p]->(s-[ε]->(s-[Literal:a]->(s-[ε]->)))(s-[ε]->(s-[Literal:b]->(s-[ε]-><back>))))"
	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestToDigraph(t *testing.T) {
	expected := `
digraph {
s1->s2 [label=ε]
s2->s4 [label=a]
s4->s6 [label=ε]
s4->s2 [label=ε]
s1->s6 [label=ε]
}
`
	ast := nb.Star(nb.Lit('a'))
	nfa := Compile(ast)
	actual := nfa.ToDigraph()
	if strings.TrimSpace(actual) != strings.TrimSpace(expected) {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
