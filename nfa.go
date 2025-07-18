package main

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	META    = "META"
	LITERAL = "LITERAL"
)

type Nfa struct {
	Start  *State
	Accept *State
}

type TransitionType int

const (
	Literal TransitionType = iota
	Meta
)

func (me TransitionType) String() string {
	return [...]string{"Literal", "Meta"}[me]
}

type Transition struct {
	Type  TransitionType
	Condition string
	State *State
}

type State struct {
	Transitions []Transition
	Epsilon     []*State
}

func NewNfa() Nfa {
	nfa := Nfa{}
	nfa.NewStart()
	nfa.NewAccept()
	return nfa
}

func (s *State) AddTransition(kind TransitionType, condition string, toState *State) {
	t := Transition{Type: kind, Condition: condition}
	t.State = toState
	s.Transitions = append(s.Transitions, t)
}

func (n *Nfa) NewState() *State {
	state := &State{}
	return state
}

func getTransitionType(str string) TransitionType {
	if str == DOT || str == WHITESPACE {
		return Meta
	}
	return Literal
}

func (n *Nfa) NewStart() *State {
	start := &State{}
	n.Start = start
	return start
}

func (s *State) AddEpsilonTo(to *State) {
	s.Epsilon = append(s.Epsilon, to)
}

func (n *Nfa) NewAccept() *State {
	accept := &State{}
	n.Accept = accept
	return accept
}

func (n *Nfa) Encode() string {
	encoded := ""
	encoded += n.Start.Encode()
	return encoded
}

func (s *State) Encode() string {
	seen := make(map[*State]bool)

	var encode func(s *State) string
	encode = func(s *State) string {
		if seen[s] {
			return "<back>"
		}
		seen[s] = true

		var parts []string

		sort.Slice(s.Transitions, func(i, j int) bool {
			return s.Transitions[i].Condition < s.Transitions[j].Condition
		})

		var encodings []string
		for _, t := range s.Transitions {
			encodedState := encode(t.State)
			encodings = append(encodings, encodedState)
			parts = append(parts, fmt.Sprintf("(s-[%s:%s]->%s)", t.Type, t.Condition, encodedState))
		}

		var epsilonEncodings []string
		for _, state := range s.Epsilon {
			epsilonEncodings = append(epsilonEncodings, encode(state))
		}
		slices.Sort(epsilonEncodings)
		for _, child := range epsilonEncodings {
			parts = append(parts, fmt.Sprintf("(s-[ε]->%s)", child))
		}

		return strings.Join(parts, "")
	}

	return encode(s)
}

func (n *Nfa) ToDigraph() string {
	used := make(map[*State]bool)
	names := make(map[*State]string)
	result := "digraph {\n"
	counter := 0
	var name func(s *State) string
	name = func(s *State) string {
		key, ok := names[s]
		counter++
		if ok {
			return key
		} else {
			names[s] = "s" + strconv.Itoa(counter)
			return names[s]
		}
	}
	var toEdge func(s *State)
	toEdge = func(s *State) {
		isUsed, ok := used[s]
		if ok && isUsed {
			return
		}
		used[s] = true
		for _, t := range s.Transitions {
			result += fmt.Sprintf("%s->%s [label=%s]\n", name(s), name(t.State), t.Condition)
			toEdge(t.State)
		}
		for _, state := range s.Epsilon {
			result += fmt.Sprintf("%s->%s [label=ε]\n", name(s), name(state))
			toEdge(state)
		}
	}

	toEdge(n.Start)
	result += "}\n"
	return result
}
