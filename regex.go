package main

import (
	"fmt"
	"slices"
	"unicode"
)

func Compile(n Node) Nfa {
	initMatchers()
	return compileNode(n)
}

type matcherFunc func(t Transition, char rune) bool

var matchers map[TransitionType]matcherFunc

func initMatchers() {
	matchers = map[TransitionType]matcherFunc{
		Literal: matchLiteral,
		Meta:    matchMeta,
	}
}

func compileNode(n Node) Nfa {
	switch n := n.(type) {
	case *LiteralNode:
		return compileLiteral(n)
	case *SequenceNode:
		return compileSequence(n)
	case *StarNode:
		return compileStar(n)
	case *CharList:
		return compileCharList(n)
	case *MetaCharacterNode:
		return compileMetaCharacter(n)
	default:
		panic(fmt.Sprintf("Unknown node type %T", n))
	}
}

func compileLiteral(n *LiteralNode) Nfa {
	nfa := NewNfa()
	nfa.Start.AddTransition(Literal, string(n.Value), nfa.Accept)
	return nfa
}

func compileMetaCharacter(n *MetaCharacterNode) Nfa {
	nfa := NewNfa()
	nfa.Start.AddTransition(Meta, n.Value, nfa.Accept)
	return nfa
}

func compileSequence(n *SequenceNode) Nfa {
	nfa := compileNode(n.Children[0])
	for i := 1; i < len(n.Children); i++ {
		childNfa := compileNode(n.Children[i])
		nfa = concat(nfa, childNfa)
	}
	return nfa
}

func compileStar(n *StarNode) Nfa {
	nfa := NewNfa()
	childNfa := compileNode(n.Child)
	// TODO: changing the order of following epsilon transitions breaks the encoding
	nfa.Start.AddEpsilonTo(childNfa.Start)
	nfa.Start.AddEpsilonTo(nfa.Accept)
	childNfa.Accept.AddEpsilonTo(nfa.Accept)
	childNfa.Accept.AddEpsilonTo(childNfa.Start)
	return nfa
}

func compileCharList(n *CharList) Nfa {
	var charListNfa Nfa
	for _, node := range n.Chars {
		if charListNfa == (Nfa{}) {
			charListNfa = compileNode(node)
		} else {
			nfa := compileNode(node)
			charListNfa = union(charListNfa, nfa)
		}
	}
	return charListNfa
}

func union(n1 Nfa, n2 Nfa) Nfa {
	nfa := NewNfa()
	nfa.Start.AddEpsilonTo(n1.Start)
	nfa.Start.AddEpsilonTo(n2.Start)
	n1.Accept.AddEpsilonTo(nfa.Accept)
	n2.Accept.AddEpsilonTo(nfa.Accept)
	return nfa
}

func concat(n1 Nfa, n2 Nfa) Nfa {
	nfa := NewNfa()
	nfa.Start = n1.Start
	nfa.Accept = n2.Accept
	n1.Accept.Transitions = n2.Start.Transitions
	n1.Accept.Epsilon = n2.Start.Epsilon
	n1.Accept = n2.Start
	return nfa
}

func matchLiteral(t Transition, char rune) bool {
	return t.Condition == string(char)
}

func matchMeta(t Transition, char rune) bool {
	if t.Condition == DOT {
		return true
	} else if t.Condition == WHITESPACE {
		return unicode.IsSpace(char)
	} else if t.Condition == NONWHITESPACE {
		switch char {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			return false
		default:
			return true
		}
	}
	return false
}

func Match(n Nfa, input string) bool {
	states := closures(n.Start)
	for i, char := range input {
		fmt.Printf("checking character %d=%s\n", i, string(char))
		var nextStates []*State
		visited := make(map[*State]bool)
		for _, s := range states {
			var targetStates []*State
			for _, t := range s.Transitions {
				if matchers[t.Type](t, char) {
					targetStates = append(targetStates, t.State)
				}
			}
			for _, ts := range targetStates {
				closureStates := closures(ts)
				for _, c := range closureStates {
					if visited[c] == false {
						visited[c] = true
						nextStates = append(nextStates, c)
					} else {
						fmt.Println("Visited before")
					}
				}
			}
		}
		fmt.Printf("Count of next states=%d", len(nextStates))
		states = nextStates
	}
	return slices.Contains(states, n.Accept)
}

func closures(n *State) []*State {
	var states []*State

	var findClosures func(childState *State)
	findClosures = func(childState *State) {
		if slices.Contains(states, childState) == false {
			states = append(states, childState)
		}
		for _, epsilonState := range childState.Epsilon {
			findClosures(epsilonState)
		}
	}
	findClosures(n)
	return states
}
