package main

func Compile(n Node) Nfa {
	nfa := Nfa{}
	switch n.(type) {
	case *LiteralNode:
		return compileLiteralNode(n.(*LiteralNode))
	case *SequenceNode:
		return compileSequenceNode(n.(*SequenceNode))
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

func concat(n1 Nfa, n2 Nfa) Nfa {
	nfa := Nfa{}
	nfa.Start = n1.Start
	nfa.Accept = n2.Accept
	n1.Accept.Transitions = n2.Start.Transitions
	n1.Accept = n2.Start
	return nfa
}
