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
	Start       *State
	Accept      *State
	StateCount  int
	StatePrefix string
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
	Name        string
	Transitions []Transition
	Epsilon     []*State
}

func NewNfa(prefix string) Nfa {
	nfa := Nfa{}
	nfa.StatePrefix = prefix
	nfa.StateCount++
	nfa.NewStart()
	nfa.NewAccept()
	return nfa
}

func (s *State) AddTransition(kind TransitionType, condition string, toState *State) {
	t := Transition{Type: kind, Value: condition}
	t.State = toState
	s.Transitions = append(s.Transitions, t)
}

func (n *Nfa) NewState(name string) State {
	state := State{Name: name}
	return state
}

func (n *Nfa) FindState(name string) *State {
	seen := make(map[string]bool)
	var findState func(name string, s *State) *State
	findState = func(name string, s *State) *State {
		if seen[s.Name] {
			return nil
		}
		seen[s.Name] = true
		if s.Name == name {
			return s
		}
		for _, transition := range s.Transitions {
			stateFound := findState(name, transition.State)
			if stateFound != nil {
				return stateFound
			}
		}
		for _, state := range s.Epsilon {
			stateFound := findState(name, state)
			if stateFound != nil {
				return stateFound
			}
		}
		return nil
	}

	return findState(name, n.Start)
}

func getTransitionType(str string) TransitionType {
	if str==DOT || str==WHITESPACE {
		return Meta
	}
	return Literal
}

func CreateNfaFromString(str string) Nfa {
	nfa := Nfa{}
	lines := strings.Split(strings.TrimSpace(str), "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ",")
		fromState := parts[0]
		toState := parts[1]
		transition := parts[2]
		if nfa.Start == nil {
			start := nfa.NewStart()
			start.Name = fromState
			to := nfa.NewState(toState)
			if transition == "ε" {
				start.AddEpsilonTo(&to)
			} else {
				start.AddTransition(getTransitionType(transition),transition, &to)
			}
		} else {
			from := nfa.FindState(fromState)
			if from != nil {
				to := nfa.NewState(toState)
				if transition == "ε" {
					from.AddEpsilonTo(&to)
				} else {
					from.AddTransition(getTransitionType(transition),transition, &to)
				}
			}
		}
	}
	return nfa
}

func (n *Nfa) NewStart() *State {
	start := &State{}
	start.Name = n.StatePrefix + strconv.Itoa(n.StateCount)
	n.Start = start
	n.StateCount++
	return start
}

func (n *Nfa) AddStart(s *State) {
	n.Start = s
}

func (n *Nfa) AddAccept(s *State) {
	n.Accept = s
}

func (s *State) AddEpsilonTo(to *State) {
	s.Epsilon = append(s.Epsilon, to)
}

func (n *Nfa) NewAccept() *State {
	accept := &State{}
	accept.Name = n.StatePrefix + strconv.Itoa(n.StateCount)
	n.Accept = accept
	n.StateCount++
	return accept
}

func (n *Nfa) ToString() string {
	str := ""
	str += n.Start.ToString()
	return str
}

func (n *Nfa) Encode() string {
	encoded := ""
	encoded += n.Start.Encode()
	return encoded
}

func (s *State) Encode() string {
	seen := make(map[string]bool)

	var encode func(s *State) string
	encode = func(s *State) string {
		if seen[s.Name] {
			return "<back>"
		}
		seen[s.Name] = true

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

func stateToString(s *State, str *string, used map[string]bool) {
	isUsed, ok := used[s.Name]
	if ok && isUsed {
		return
	}
	used[s.Name] = true
	for _, t := range s.Transitions {
		*str += fmt.Sprintf("%s --> %s, %s\n", s.Name, t.State.Name, t.Value)
		stateToString(t.State, str, used)
	}
	for _, state := range s.Epsilon {
		*str += fmt.Sprintf("%s --> %s, ε\n", s.Name, state.Name)
		stateToString(state, str, used)
	}
}

func (s *State) ToString() string {
	used := make(map[string]bool)
	str := ""
	stateToString(s, &str, used)
	return str
}
