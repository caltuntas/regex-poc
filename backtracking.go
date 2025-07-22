package main

import "unicode"

func matchNode(node Node, input string, pos int) (bool, int) {
	switch n := node.(type) {

	case *LiteralNode:
		if pos < len(input) && input[pos] == n.Value {
			return true, pos + 1
		}
		return false, pos

	case *MetaCharacterNode:
		if pos >= len(input) {
			return false, pos
		}
		c := []rune(input)[pos]
		switch n.Value {
		case DOT:
			return true, pos + 1
		case WHITESPACE:
			if unicode.IsSpace(c) {
				return true, pos + 1
			}
		case NONWHITESPACE:
			if !(c == ' ' || c == '\t' || c == '\n' || c == '\r') {
				return true, pos + 1
			}
		}
		return false, pos

	case *SequenceNode:
		current := pos
		var previous Node
		previous = n.Children[0]
		for _, child := range n.Children {
			_, isStar := previous.(*StarNode) 
			ok, next := matchNode(child, input, current)
			if !ok && isStar{
				ok, next = matchNode(child, input, current-1)
			}
			if !ok {
				return false, pos // fail early, restore original pos
			}
			current = next
			previous = child
		}
		return true, current

	case *StarNode:
		current := pos
		for {
			ok, next := matchNode(n.Child, input, current)
			if !ok {
				break
			}
			current = next
		}
		return true, current

	case *CharList:
		if pos >= len(input) {
			return false, pos
		}
		current := pos
		for _, ch := range n.Chars {
			ok, next := matchNode(ch, input, current)
			if !ok {
				continue
			} else {
				return true, next
			}
		}
		return false, current
	}
	return false, pos
}

func MatchBacktrack(ast Node, input string) bool {
	ok, next := matchNode(ast, input, 0)
	return ok && next == len(input)
}
