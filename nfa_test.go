package main


import (
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
