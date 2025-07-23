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
			if !unicode.IsSpace(c) {
				return true, pos + 1
			}
		}
		return false, pos

	case *SequenceNode:
		return matchSequence(n.Children, input, pos)

	case *StarNode:
		positions := []int{pos}
		current := pos
		for {
			ok, next := matchNode(n.Child, input, current)
			if !ok || next == current {
				break
			}
			current = next
			positions = append(positions, current)
		}
		return true, current

	case *CharList:
		if pos >= len(input) {
			return false, pos
		}
		current := pos
		for _, ch := range n.Chars {
			ok, next := matchNode(ch, input, current)
			if ok {
				return true, next
			}
		}
		return false, current
	}
	return false, pos
}

func matchSequence(nodes []Node, input string, pos int) (bool, int) {
	if len(nodes) == 0 {
		return true, pos
	}

	first := nodes[0]

	if star, ok := first.(*StarNode); ok {
		positions := []int{pos}
		current := pos
		for {
			ok, next := matchNode(star.Child, input, current)
			if !ok || next == current {
				break
			}
			current = next
			positions = append(positions, current)
		}

		for i := len(positions) - 1; i >= 0; i-- {
			if ok, next := matchSequence(nodes[1:], input, positions[i]); ok {
				return true, next
			}
		}
		return false, pos
	}

	ok, next := matchNode(first, input, pos)
	if !ok {
		return false, pos
	}
	return matchSequence(nodes[1:], input, next)
}

func MatchBacktrack(ast Node, input string) bool {
	ok, next := matchNode(ast, input, 0)
	return ok && next == len(input)
}

func MatchBacktrackPartial(ast Node, input string) bool {
    for start := 0; start <= len(input); start++ {
        ok, _ := matchNode(ast, input, start)
        if ok {
            return true
        }
    }
    return false
}
