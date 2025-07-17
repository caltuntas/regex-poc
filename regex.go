package main

import (
	"fmt"
	"slices"
	"unicode"
)

func Compile(n Node) Nfa {
	nfa := NewNfa("s")
	initMatchers()
	return compileNode(&nfa, n)
}

type matcherFunc func(t Transition, char rune) bool
var matchers map[string]matcherFunc

func initMatchers() {
	matchers = map[string]matcherFunc {
		LITERAL: matchLiteral,
		META: matchMeta,
	}
}

func compileNode(nfa *Nfa, n Node) Nfa {
	switch n:=n.(type) {
	case *LiteralNode:
		return compileLiteralNode(nfa, n)
	case *SequenceNode:
		return compileSequenceNode(nfa, n)
	case *StarNode:
		return compileStarNode(nfa, n)
	case *CharList:
		return compileCharList(nfa, n)
	case *MetaCharacterNode:
		return compileMetaCharacter(nfa, n)
	default:
		panic(fmt.Sprintf("Unknown node type %T", n))
	}
}

func compileLiteralNode(parentNfa *Nfa, n *LiteralNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
	accept := nfa.NewAccept()
	start.AddTransition(string(n.Value),accept)
	nfa.Start = start
	nfa.Accept = accept
	parentNfa.StateCount = nfa.StateCount
	return nfa
}

func compileMetaCharacter(parentNfa *Nfa, n *MetaCharacterNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
	accept := nfa.NewAccept()
	start.AddTransition(string(n.Value),accept)
	nfa.Start = start
	nfa.Accept = accept
	parentNfa.StateCount = nfa.StateCount
	return nfa
}

func compileSequenceNode(parentNfa *Nfa, n *SequenceNode) Nfa {
	nfa := compileNode(parentNfa, n.Children[0])
	for i := 1; i < len(n.Children); i++ {
		childNfa := compileNode(parentNfa, n.Children[i])
		nfa = concat(parentNfa, nfa, childNfa)
	}
	return nfa
}

func compileStarNode(parentNfa *Nfa, n *StarNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
	accept := nfa.NewAccept()
	parentNfa.StateCount = nfa.StateCount
	childNfa := compileNode(parentNfa, n.Child)
	childStart := childNfa.Start
	childAccept := childNfa.Accept
	// TODO: changing the order of following epsilon transitions breaks the encoding
	start.AddEpsilonTo(childStart)
	start.AddEpsilonTo(accept)
	childAccept.AddEpsilonTo(accept)
	childAccept.AddEpsilonTo(childStart)
	return nfa
}

func compileCharList(parentNfa *Nfa, n *CharList) Nfa {
	var charListNfa Nfa
	for _, node := range n.Chars {
		if charListNfa == (Nfa{}) {
			charListNfa = compileNode(parentNfa,node)
		} else {
			nfa := compileNode(parentNfa, node)
			charListNfa = union(parentNfa, charListNfa, nfa)
		}
	}
	return charListNfa
}

func union(parentNfa *Nfa, n1 Nfa, n2 Nfa) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
  accept := nfa.NewAccept()
	start.AddEpsilonTo(n1.Start)
	start.AddEpsilonTo(n2.Start)
	n1.Accept.AddEpsilonTo(accept)
	n2.Accept.AddEpsilonTo(accept)
	return nfa
}

func concat(parentNfa *Nfa, n1 Nfa, n2 Nfa) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	nfa.AddStart(n1.Start)
	nfa.AddAccept(n2.Accept)
	n1.Accept.Transitions = n2.Start.Transitions
	n1.Accept.Epsilon = n2.Start.Epsilon
	n1.Accept = n2.Start
	return nfa
}

func matchLiteral(t Transition, char rune) bool {
	return t.Value == string(char)
}

func matchMeta(t Transition, char rune) bool {
	if t.Value == DOT {
		return true
	} else if t.Value == WHITESPACE {
		return unicode.IsSpace(char)
	}
	return false
}

func Match(n Nfa, input string) bool {
	states := closures(n.Start)
	for _,char := range input {
		var nextStates []*State
		for _, s := range states {
			var targetStates []*State
			for _,t := range s.Transitions {
				if matchers[t.Type](t,char) {
					targetStates = append(targetStates,t.State)
				}
			}
			for _, ts := range targetStates {
				closureStates := closures(ts)
				nextStates = append(nextStates, closureStates...)
			}
		}
		states = nextStates
	}

	return slices.Contains(states, n.Accept)
}

func closures(n *State) []*State {
	var states []*State
  
	var  findClosures func(childState *State)
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
