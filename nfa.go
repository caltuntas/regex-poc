package main

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

type Nfa struct {
	Start       *State
	Accept      *State
	StateCount  int
	StatePrefix string
}

type State struct {
	Name        string
	// TODO: Introduce a new struct for Transition
	// map key is not enough for covering different regex constructs
	Transitions map[string][]*State
	Epsilon     []*State
}

func NewNfa(prefix string) Nfa {
	nfa := Nfa{}
	nfa.StatePrefix = prefix
	nfa.StateCount++
	return nfa
}

func (s *State) AddTransition(key string, toState *State) {
	var transitions []*State
	transitions = append(transitions, toState)
	s.Transitions = make(map[string][]*State)
	s.Transitions[key] = transitions
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
		for _, states := range s.Transitions {
			for _, state := range states {
				stateFound := findState(name, state)
				if stateFound!=nil {
					return stateFound
				} 
			}
		}
		for _, state := range s.Epsilon {
			stateFound := findState(name, state)
			if stateFound!=nil {
				return stateFound
			} 
		}
		return nil
	}

	return findState(name,n.Start)
}

func CreateNfaFromString(str string) Nfa {
	nfa := Nfa{}
	lines := strings.Split(strings.TrimSpace(str),"\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line),",")
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
				start.AddTransition(transition,&to)
			}
		} else {
			from := nfa.FindState(fromState)
			if from != nil {
				to := nfa.NewState(toState)
				if transition == "ε" {
					from.AddEpsilonTo(&to)
				} else {
					from.AddTransition(transition,&to)
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

		transitionKeys := slices.Collect(maps.Keys(s.Transitions))
		slices.Sort(transitionKeys)

		for _, key := range transitionKeys {
			var encodings []string
			for _, state := range s.Transitions[key] {
				encodings = append(encodings, encode(state))
			}
			slices.Sort(encodings)
			for _, child := range encodings {
				parts = append(parts, fmt.Sprintf("(s-[%s]->%s)", string(key), child))
			}
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
	for key, value := range s.Transitions {
		for _, state := range value {
			*str += fmt.Sprintf("%s --> %s, %s\n", s.Name, state.Name, string(key))
			stateToString(state, str, used)
		}
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
