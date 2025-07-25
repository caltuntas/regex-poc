package main

import (
	"fmt"
	"unicode"
)

func printPosition(input string, pos int, label string) {
	if pos >= len(input) {
		for i:=len(input); i<=pos; i++ {
			input += "_"
		}
	} else {
		input += " "
	}
	for _, ch := range input {
		fmt.Printf("%c ", ch)
	}
	fmt.Printf("--> %s\n", label)

	for i := 0; i < len(input); i++ {
		if i == pos {
			fmt.Print("^ ")
		} else {
			fmt.Print("  ")
		}
	}
	fmt.Println()
}

func matchNode(node Node, input string, pos int) (bool, int) {
	switch n := node.(type) {

	case *LiteralNode:
		printPosition(input, pos, string(n.Value))
		if pos < len(input) && input[pos] == n.Value {
			return true, pos + 1
		}
		return false, pos

	case *MetaCharacterNode:
		printPosition(input, pos,n.Value)
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
		current := pos
		for i := 0; i < len(n.Children); i++ {
			child := n.Children[i]
			if star, ok := child.(*StarNode); ok {
				printPosition(input, current, star.String())
				positions := []int{current}
				nextPos := current
				for {
					ok, next := matchNode(star.Child, input, nextPos)
					if !ok || next == nextPos {
						break
					}
					nextPos = next
					positions = append(positions, nextPos)
				}

				for j := len(positions) - 1; j >= 0; j-- {
					ok, endPos := matchNode(&SequenceNode{Children: n.Children[i+1:]}, input, positions[j])
					if ok {
						return true, endPos
					}
				}
				return false, pos
			}

			ok, next := matchNode(child, input, current)
			if !ok {
				return false, pos
			}
			current = next
		}
		return true, current

	case *CharList:
		if pos >= len(input) {
			return false, pos
		}
		for _, chNode := range n.Chars {
			ok, next := matchNode(chNode, input, pos)
			if ok {
				return true, next
			}
		}
		return false, pos
	}

	return false, pos
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
