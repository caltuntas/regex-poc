package main

import (
	"fmt"
	"strconv"
)

type Nfa struct {
	Start  *State
	Accept *State
	StateCount int
	StatePrefix string
}

type State struct {
	Name        string
	Transitions map[byte][]*State
	Epsilon     []*State
}

func NewNfa(prefix string) Nfa {
	nfa := Nfa{}
	nfa.StatePrefix = prefix
	nfa.StateCount++
	return nfa
}


func (n *Nfa) NewStart() *State {
	start := &State{}
	start.Name = n.StatePrefix + strconv.Itoa(n.StateCount)
	n.Start = start
	n.StateCount++
	return start
}

func (n* Nfa) AddStart(s *State) {
	n.Start = s
}

func (n* Nfa) AddAccept(s *State) {
	n.Accept = s
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

func stateToString(s *State, str *string, used map[string]bool) {
	isUsed,ok := used[s.Name]
	if ok && isUsed {
		return 
	}
	used[s.Name] = true
	for key, value := range s.Transitions {
		for _, state := range value {
			*str += fmt.Sprintf("%s --> %s, %s\n", s.Name, state.Name, string(key))
			 stateToString(state,str, used)
		}
	}
	for _, state := range s.Epsilon {
		*str += fmt.Sprintf("%s --> %s, ε\n", s.Name, state.Name)
		 stateToString(state,str, used)
	}
}

func (s *State) ToString() string {
	used := make(map[string]bool)
	str :=""
	 stateToString(s,&str,used)
	return str
}
