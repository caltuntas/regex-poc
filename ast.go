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
	str := ""
	str += fmt.Sprintln("StarNode, Child")
	str += n.Child.String()
  return str
}

type CharList struct {
	Chars []CharacterNode
}

func (n *CharList) String() string {
	return "charlist"
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
