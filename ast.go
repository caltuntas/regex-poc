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

type NodeBuilder struct {
}

func (b NodeBuilder) Lit(c byte) *LiteralNode { return &LiteralNode{Value: c}}
func (b NodeBuilder) Star(child Node) *StarNode { return &StarNode{Child: child}}
func (b NodeBuilder) Dot() *DotNode { return &DotNode{}}
func (b NodeBuilder) Seq(children ...Node) *SequenceNode { 
	return &SequenceNode{ Children: children}
}
