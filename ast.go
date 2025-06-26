package main

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
	return ""
}

type LiteralNode struct {
	Value byte
}

func (n* LiteralNode) String() string {
	return ""
}

type DotNode struct {
	Value byte
}

func (n *DotNode) String() string {
	return "."
}

func (n *StarNode) String() string {
  return "*"
}

