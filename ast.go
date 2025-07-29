package main

import "fmt"

const (
	WHITESPACE = "\\s"
	NONWHITESPACE = "\\S"
)

type Node interface {
	String() string
}

type StarNode struct {
	Child Node
}

type SequenceNode struct {
	Children []Node;
}

func (n *SequenceNode) String() string {
	str := ""
	str += fmt.Sprintln("SequenceNode, Children")
	for i,child := range n.Children {
		str += fmt.Sprintf("Child %d = %s\n", i, child.String())
	}
	return str
}

type CharacterNode interface {
	GetValue() string
	String() string
}

type LiteralNode struct {
	Value byte
}

func (n* LiteralNode) String() string {
	return string(n.Value)
}

func (n* LiteralNode) GetValue() string {
	return n.String()
}

type MetaCharacterNode struct {
	Value string
}

func (n* MetaCharacterNode) String() string {
	return n.Value
}

func (n* MetaCharacterNode) GetValue() string {
	return n.String()
}

func (n *StarNode) String() string {
	return n.Child.String() + "*"
}

type CharList struct {
	Chars []CharacterNode
}

func (n *CharList) String() string {
	str := "["
	for _,cn := range n.Chars {
		str += cn.String()
	}
	return str + "]"
}

type NodeBuilder struct {
}

func (b NodeBuilder) Lit(c byte) *LiteralNode { return &LiteralNode{Value: c}}
func (b NodeBuilder) Star(child Node) *StarNode { return &StarNode{Child: child}}
func (b NodeBuilder) Seq(children ...Node) *SequenceNode { 
	return &SequenceNode{ Children: children}
}
func (b NodeBuilder) List(chars ...CharacterNode) *CharList { 
	return &CharList{ Chars: chars}
}

func (b NodeBuilder) Meta(s string) *MetaCharacterNode {
	return &MetaCharacterNode{Value: s}
}

func PrintAstTree(node Node, indentLevel int) {
	indentSize := 2
	indent := indentLevel * indentSize
	switch n := node.(type) {
	case *LiteralNode:
		fmt.Printf("%*sLiteral: '%s'\n", indent, "", n.String())
	case *MetaCharacterNode:
		fmt.Printf("%*sMeta: '%s'\n", indent, "", n.String())
	case *CharList:
		fmt.Printf("%*sCharList:\n", indent, "")
		for _, ch := range n.Chars {
			switch c := ch.(type) {
			case *LiteralNode:
				fmt.Printf("%*sLiteral: '%s'\n", indent+indentSize, "", c.String())
			case *MetaCharacterNode:
				fmt.Printf("%*sMeta: '%s'\n", indent+indentSize, "", c.String())
			default:
				fmt.Printf("%*sUnknown CharacterNode\n", indent+indentSize, "")
			}
		}
	case *StarNode:
		fmt.Printf("%*sStar:\n", indent, "")
		PrintAstTree(n.Child, indentLevel+1)
	case *SequenceNode:
		fmt.Printf("%*sSequence:\n", indent, "")
		for _, child := range n.Children {
			PrintAstTree(child, indentLevel+2)
		}
	default:
		fmt.Printf("%*sUnknown Node type\n", indent, "")
	}
}
