package main

import (
	"fmt"
	"testing"
)

var nb NodeBuilder

func TestLiteralNodeToNFA(t *testing.T) {
	ast := nb.Lit('p')

	nfa := Compile(ast)

	if nfa.Start == nil {
		t.Fatal("Start state is nil")
	}
	if nfa.Accept == nil {
		t.Fatal("Accept state is nil")
	}

	transitions := nfa.Start.Transitions['p']
	if len(transitions) != 1 {
		t.Fatalf("Expected one transition on 'p', got %d", len(transitions))
	}

	if transitions[0] != nfa.Accept {
		t.Errorf("Expected transition to accept state, got %v", transitions[0])
	}

	if len(nfa.Start.Epsilon) != 0 {
		t.Errorf("Start state should not have epsilon transitions")
	}
	if len(nfa.Accept.Transitions) != 0 && len(nfa.Accept.Epsilon) != 0 {
		t.Errorf("Accept state should have no outgoing transitions")
	}
}

func TestSequenceLiteralToNFA(t *testing.T) {
	ast := nb.Seq(
		nb.Lit('p'),
		nb.Lit('a'),
	)

	nfa := Compile(ast)

	if nfa.Start == nil || nfa.Accept == nil {
		t.Fatal("Start or Accept state is nil")
	}

	pTargets := nfa.Start.Transitions['p']
	if len(pTargets) != 1 {
		t.Fatalf("Expected 1 transition on 'p', got %d", len(pTargets))
	}
	intermediate := pTargets[0]

	aTargets := intermediate.Transitions['a']
	if len(aTargets) != 1 {
		t.Fatalf("Expected 1 transition on 'a', got %d", len(aTargets))
	}
	if aTargets[0] != nfa.Accept {
		t.Errorf("Expected 'a' to go to accept state, got %v", aTargets[0])
	}

	if len(nfa.Start.Epsilon) != 0 {
		t.Errorf("Start should not have epsilon transitions")
	}
	if len(intermediate.Epsilon) != 0 {
		t.Errorf("Intermediate should not have epsilon transitions")
	}
	if len(nfa.Accept.Transitions) != 0 && len(nfa.Accept.Epsilon) != 0 {
		t.Errorf("Accept state should not have outgoing transitions")
	}
}

func TestStarLiteralToNFA(t *testing.T) {
	ast := nb.Star(nb.Lit('a'))
	nfa := Compile(ast)

	if nfa.Start == nil || nfa.Accept == nil {
		t.Fatal("Start or Accept state is nil")
	}

	start := nfa.Start
	accept := nfa.Accept

	if len(start.Epsilon) != 2 {
		t.Fatalf("Start state should have 2 ε-transitions, got %d", len(start.Epsilon))
	}

	var entryA *State
	if start.Epsilon[0] == accept {
		entryA = start.Epsilon[1]
	} else if start.Epsilon[1] == accept {
		entryA = start.Epsilon[0]
	} else {
		entryA = start.Epsilon[0]
	}

	aTargets := entryA.Transitions['a']
	if len(aTargets) != 1 {
		t.Fatalf("Expected one 'a' transition from entryA, got %d", len(aTargets))
	}
	mid := aTargets[0]

	foundBack := false
	foundExit := false
	for _, eps := range mid.Epsilon {
		if eps == entryA {
			foundBack = true
		}
		if eps == accept {
			foundExit = true
		}
	}
	if !foundBack {
		t.Errorf("Expected ε-transition from mid back to entryA")
	}
	if !foundExit {
		t.Errorf("Expected ε-transition from mid to accept")
	}
}

func TestSequenceWithStarToNFA(t *testing.T) {
	ast := nb.Seq(
			nb.Lit('a'),
			nb.Star(nb.Lit('b')),
	)

	nfa := Compile(ast)
	str := nfa.ToString()
	fmt.Println(str)

	if nfa.Start == nil || nfa.Accept == nil {
		t.Fatal("Start or Accept state is nil")
	}

	start := nfa.Start
	accept := nfa.Accept

	aTargets := start.Transitions['a']
	if len(aTargets) != 1 {
		t.Fatalf("Expected 1 'a' transition from start, got %d", len(aTargets))
	}
	s2 := aTargets[0]


	if len(s2.Epsilon) != 2 {
		t.Fatalf("Expected 2 ε-transitions from s2 (entry of b*), got %d", len(s2.Epsilon))
	}

	var entryB *State
	if s2.Epsilon[0] == accept {
		entryB = s2.Epsilon[1]
	} else if s2.Epsilon[1] == accept {
		entryB = s2.Epsilon[0]
	} else {
		t.Fatalf("Expected one ε-transition from s2 to accept")
	}

	bTargets := entryB.Transitions['b']
	if len(bTargets) != 1 {
		t.Fatalf("Expected 1 'b' transition from entryB, got %d", len(bTargets))
	}
	mid := bTargets[0]

	foundBack := false
	foundExit := false
	for _, eps := range mid.Epsilon {
		if eps == entryB {
			foundBack = true
		}
		if eps == accept {
			foundExit = true
		}
	}
	if !foundBack {
		t.Errorf("Expected ε-transition from mid to entryB for looping")
	}
	if !foundExit {
		t.Errorf("Expected ε-transition from mid to accept")
	}
}

func TestSequenceWithDotStarToNFA(t *testing.T) {
	ast := nb.Seq(
			nb.Lit('a'),
			nb.Star(nb.Dot()),
	)

	nfa := Compile(ast)
	str := nfa.ToString()
	fmt.Println(str)

	if nfa.Start == nil || nfa.Accept == nil {
		t.Fatal("Start or Accept state is nil")
	}

	start := nfa.Start
	accept := nfa.Accept

	aTargets := start.Transitions['a']
	if len(aTargets) != 1 {
		t.Fatalf("Expected 1 'a' transition from start, got %d", len(aTargets))
	}
	s2 := aTargets[0]


	if len(s2.Epsilon) != 2 {
		t.Fatalf("Expected 2 ε-transitions from s2 (entry of .*), got %d", len(s2.Epsilon))
	}

	var entryB *State
	if s2.Epsilon[0] == accept {
		entryB = s2.Epsilon[1]
	} else if s2.Epsilon[1] == accept {
		entryB = s2.Epsilon[0]
	} else {
		t.Fatalf("Expected one ε-transition from s2 to accept")
	}

	bTargets := entryB.Transitions['.']
	if len(bTargets) != 1 {
		t.Fatalf("Expected 1 '.' transition from entryB, got %d", len(bTargets))
	}
	mid := bTargets[0]

	foundBack := false
	foundExit := false
	for _, eps := range mid.Epsilon {
		if eps == entryB {
			foundBack = true
		}
		if eps == accept {
			foundExit = true
		}
	}
	if !foundBack {
		t.Errorf("Expected ε-transition from mid to entryB for looping")
	}
	if !foundExit {
		t.Errorf("Expected ε-transition from mid to accept")
	}
}
