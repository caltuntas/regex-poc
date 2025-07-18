package main

import (
	"testing"
	"strings"
)

var nb NodeBuilder

func TestLiteralNodeToNFA(t *testing.T) {
	ast := nb.Lit('p')

	nfa := Compile(ast)
	expected := `
		state1,state2,p
	`
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestSequenceLiteralToNFA(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('p'),
		nb.Lit('a'),
	)
	nfa := Compile(ast)

	expected := `
		state1,state2,p
		state2,state3,a
	`
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestStarLiteralToNFA(t *testing.T) {
	ast := nb.Star(nb.Lit('a'))
	nfa := Compile(ast)

	expected := `
		state1,state2,ε
		state2,state3,a
		state3,state2,ε
		state3,state4,ε
		state1,state4,ε
	`
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestSequenceWithStarToNFA(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('a'),
		nb.Star(nb.Lit('b')),
	)

	nfa := Compile(ast)
	expected := `
		state1,state2,a
		state2,state3,ε
		state3,state4,b
		state4,state5,ε
		state4,state3,ε
		state2,state5,ε
	`

	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestSequenceWithDotStarToNFANew(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('a'),
		nb.Star(nb.Meta(DOT)),
	)

	expected := `
		state1,state2,a
		state2,state3,ε
		state3,state4,.
		state4,state5,ε
		state4,state3,ε
		state2,state5,ε
	`

	nfa := Compile(ast)
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestSequenceWithMetaCharacterToNFA(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('a'),
		nb.Meta(WHITESPACE),
	)

	expected := `
		state1,state2,a
		state2,state3,\s
	`

	nfa := Compile(ast)
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestPatternWithCharListToNFA(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('p'),
		nb.List(nb.Lit('a'), nb.Lit('b')),
	)

	expected := `
		state1,state2,p
		state2,state3,ε
		state3,state4,a
		state4,state8,ε
		state2,state5,ε
		state5,state6,b
		state6,state8,ε
	`

	nfa := Compile(ast)
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestDotLiteralToNFA(t *testing.T) {
	ast := nb.Meta(DOT)
	expected := `
		state1,state2,.
	`
	nfa := Compile(ast)
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
	}
}

func TestStarWithCharsetToNFA(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('a'),
		nb.Star(nb.List(nb.Lit('b'), nb.Lit('c'))),
		nb.Lit('d'),
	)
	expected := `
state1,state2,a
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
state2,state7,ε
	`
	nfa := Compile(ast)
	actualEncoding := nfa.Encode()
	expectedNfa := CreateNfaFromString(expected)
	expectedEncoding := expectedNfa.Encode()

	if actualEncoding != expectedEncoding {
		t.Fatalf("NFA mismatch:\nGot:\n%s\nExpected:\n%s", actualEncoding, expectedEncoding)
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
