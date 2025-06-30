package main

type Nfa struct {
	Start *State
	Accept *State
}

type State struct {
	Transitions map[byte][]*State
	Epsilon []*State
}
