package main

import (
	"fmt"
	"slices"
)

func Compile(n Node) Nfa {
	nfa := NewNfa("s")
	return compileNode(&nfa, n)
}

func compileNode(nfa *Nfa, n Node) Nfa {
	switch n.(type) {
	case *LiteralNode:
		return compileLiteralNode(nfa, n.(*LiteralNode))
	case *SequenceNode:
		return compileSequenceNode(nfa, n.(*SequenceNode))
	case *StarNode:
		return compileStartNode(nfa, n.(*StarNode))
	case *CharList:
		return compileCharList(nfa, n.(*CharList))
	case *MetaCharacterNode:
		return compileMetaCharacter(nfa, n.(*MetaCharacterNode))
	default:
		panic(fmt.Sprintf("Unknown node type %T", n))
	}
}

func compileLiteralNode(parentNfa *Nfa, n *LiteralNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
	accept := nfa.NewAccept()

	var transitions []*State
	transitions = append(transitions, accept)
	start.Transitions = make(map[string][]*State)
	start.Transitions[string(n.Value)] = transitions
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

	var transitions []*State
	transitions = append(transitions, accept)
	start.Transitions = make(map[string][]*State)
	start.Transitions[n.Value] = transitions
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

func compileStartNode(parentNfa *Nfa, n *StarNode) Nfa {
	nfa := NewNfa(parentNfa.StatePrefix)
	nfa.StateCount = parentNfa.StateCount
	start := nfa.NewStart()
	accept := nfa.NewAccept()
	parentNfa.StateCount = nfa.StateCount
	childNfa := compileNode(parentNfa, n.Child)
	childStart := childNfa.Start
	childAccept := childNfa.Accept
	start.Epsilon = append(start.Epsilon, childStart)
	start.Epsilon = append(start.Epsilon, accept)
	childAccept.Epsilon = append(childAccept.Epsilon, accept)
	childAccept.Epsilon = append(childAccept.Epsilon, childStart)
	nfa.Start = start
	nfa.Accept = accept
	return nfa
}

func compileCharList(parentNfa *Nfa, n *CharList) Nfa {
	var charListNfa Nfa
	for i:=0; i<len(n.Chars); i++ {
		char := n.Chars[i]
		ln := &LiteralNode{Value: char }
		if charListNfa == (Nfa{}) {
			charListNfa = compileLiteralNode(parentNfa, ln)
		} else {
			nfa := compileLiteralNode(parentNfa, ln)
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

func Match(n Nfa, input string) bool {
	states := closures(n.Start)
	for i := 0; i < len(input); i++ {
		var nextStates []*State
		char := input[i]
		for _, s := range states {
			var targetStates []*State
			charStates, keyFound := s.Transitions[string(char)]
			if keyFound {
				targetStates = charStates
			} else {
				dotStates, dotFound := s.Transitions["."]
				if dotFound {
					targetStates = dotStates
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
