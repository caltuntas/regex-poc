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
	Value string
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
	t := Transition{Type: kind, Value: condition}
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

func CreateNfaFromString(str string) Nfa {
	var stateMap map[string]*State
	stateMap =make(map[string]*State)
	nfa := Nfa{}
	lines := strings.Split(strings.TrimSpace(str), "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ",")
		fromState := parts[0]
		toState := parts[1]
		transition := parts[2]
		if nfa.Start == nil {
			start := nfa.NewStart()
			stateMap[fromState] = start
			to := State{}
			stateMap[toState] = &to
			if transition == "ε" {
				start.AddEpsilonTo(&to)
			} else {
				start.AddTransition(getTransitionType(transition), transition, &to)
			}
		} else {
			from := stateMap[fromState]
			if from == nil {
				state := nfa.NewState()
				stateMap[fromState] = state
				from = state
			}
			to := stateMap[toState]
			if to == nil {
				state := nfa.NewState()
			 stateMap[toState] = state
				to = state
			}
			if transition == "ε" {
				from.AddEpsilonTo(to)
			} else {
				from.AddTransition(getTransitionType(transition), transition, to)
			}
		}
	}
	return nfa
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
			return s.Transitions[i].Value < s.Transitions[j].Value
		})

		var encodings []string
		for _, t := range s.Transitions {
			encodedState := encode(t.State)
			encodings = append(encodings, encodedState)
			parts = append(parts, fmt.Sprintf("(s-[%s:%s]->%s)", t.Type, t.Value, encodedState))
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
			result += fmt.Sprintf("%s->%s [label=%s]\n", name(s), name(t.State), t.Value)
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
