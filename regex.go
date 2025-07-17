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
var matchers map[TransitionType]matcherFunc

func initMatchers() {
	matchers = map[TransitionType]matcherFunc {
		Literal: matchLiteral,
		Meta: matchMeta,
	}
}

func compileNode(nfa *Nfa, n Node) Nfa {
	switch n:=n.(type) {
	case *LiteralNode:
		return compileLiteral(nfa, n)
	case *SequenceNode:
		return compileSequence(nfa, n)
	case *StarNode:
		return compileStar(nfa, n)
	case *CharList:
		return compileCharList(nfa, n)
	case *MetaCharacterNode:
		return compileMetaCharacter(nfa, n)
	default:
		panic(fmt.Sprintf("Unknown node type %T", n))
	}
}

func compileLiteral(parentNfa *Nfa, n *LiteralNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
	accept := nfa.NewAccept()
	start.AddTransition(Literal,string(n.Value),accept)
	parentNfa.StateCount = nfa.StateCount
	return nfa
}

func compileMetaCharacter(parentNfa *Nfa, n *MetaCharacterNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
	accept := nfa.NewAccept()
	start.AddTransition(Meta, n.Value,accept)
	parentNfa.StateCount = nfa.StateCount
	return nfa
}

func compileSequence(parentNfa *Nfa, n *SequenceNode) Nfa {
	nfa := compileNode(parentNfa, n.Children[0])
	for i := 1; i < len(n.Children); i++ {
		childNfa := compileNode(parentNfa, n.Children[i])
		nfa = concat(parentNfa, nfa, childNfa)
	}
	return nfa
}

func compileStar(parentNfa *Nfa, n *StarNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	parentNfa.StateCount = nfa.StateCount
	childNfa := compileNode(parentNfa, n.Child)
	// TODO: changing the order of following epsilon transitions breaks the encoding
	nfa.Start.AddEpsilonTo(childNfa.Start)
	nfa.Start.AddEpsilonTo(nfa.Accept)
	childNfa.Accept.AddEpsilonTo(nfa.Accept)
	childNfa.Accept.AddEpsilonTo(childNfa.Start)
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
	nfa.NewStart()
  nfa.NewAccept()
	nfa.Start.AddEpsilonTo(n1.Start)
	nfa.Start.AddEpsilonTo(n2.Start)
	n1.Accept.AddEpsilonTo(nfa.Accept)
	n2.Accept.AddEpsilonTo(nfa.Accept)
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
