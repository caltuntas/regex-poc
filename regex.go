package main

func Compile(n Node) Nfa {
	nfa := Nfa{}
	switch n.(type) {
	case *LiteralNode:
		return compileLiteralNode(n.(*LiteralNode))
	case *SequenceNode:
		return compileSequenceNode(n.(*SequenceNode))
	case *StarNode:
		return compileStartNode(n.(*StarNode))
	}
	return nfa
}

func compileLiteralNode(n *LiteralNode) Nfa {
	nfa := Nfa{}
	start := &State{}
	accept := &State{}

	var transitions []*State
	transitions = append(transitions, accept)
	start.Transitions = make(map[byte][]*State)
	start.Transitions[n.Value] = transitions
	nfa.Start = start
	nfa.Accept = accept
	return nfa
}

func compileSequenceNode(n *SequenceNode) Nfa {
	nfa := Compile(n.Children[0])

	for i:=1; i<len(n.Children); i++ {
		childNfa := Compile(n.Children[i])
		nfa = concat(nfa, childNfa)
	}

	return nfa
}

func compileStartNode(n *StarNode) Nfa {
	nfa := Nfa{}
	start := &State{}
	accept := &State{}
	childNfa := Compile(n.Child)
	childStart := childNfa.Start
	childAccept := childNfa.Accept
	start.Epsilon = append(start.Epsilon, childStart)
	start.Epsilon = append(start.Epsilon, childAccept)
	childAccept.Epsilon = append(childAccept.Epsilon, accept)
	childAccept.Epsilon = append(childAccept.Epsilon, childStart)
	nfa.Start = start
	nfa.Accept = accept
	return nfa
}

func concat(n1 Nfa, n2 Nfa) Nfa {
	nfa := Nfa{}
	nfa.Start = n1.Start
	nfa.Accept = n2.Accept
	n1.Accept.Transitions = n2.Start.Transitions
	n1.Accept = n2.Start
	return nfa
}

func Match(n Nfa, input string) bool {
	states := closures(n.Start)
	for i :=0; i<len(input); i++ {
		var nextStates []*State
		char := input[i]
		for _,s := range states {
			if targetStates, ok := s.Transitions[char]; ok {
				for _,ts :=range targetStates {
					nextStates = append(nextStates, ts)
				}
			}
		}
		states = nextStates
	}

	for _,s := range states {
		if n.Accept == s {
			return true
		}
	}
	return false
}

func closures(n *State) []*State {
	var states []*State
	states = append(states, n)
	for _, s := range n.Epsilon {
    states = append(states, s)
	}
	return states
}
