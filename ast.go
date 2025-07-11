package main

import "fmt"

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

type LiteralNode struct {
	Value byte
}

func (n* LiteralNode) String() string {
	return string(n.Value)
}

type DotNode struct {
	Value byte
}

func (n *DotNode) String() string {
	return "."
}

func (n *StarNode) String() string {
	str := ""
	str += fmt.Sprintln("StarNode, Child")
	str += n.Child.String()
  return str
}

type CharList struct {
	Chars []byte
}

func (n *CharList) String() string {
	return "charlist"
}

type NodeBuilder struct {
}

func (b NodeBuilder) Lit(c byte) *LiteralNode { return &LiteralNode{Value: c}}
func (b NodeBuilder) Star(child Node) *StarNode { return &StarNode{Child: child}}
func (b NodeBuilder) Dot() *DotNode { return &DotNode{Value: '.'}}
func (b NodeBuilder) Seq(children ...Node) *SequenceNode { 
	return &SequenceNode{ Children: children}
}
func (b NodeBuilder) List(chars ...byte) *CharList { 
	return &CharList{ Chars: chars}
}
