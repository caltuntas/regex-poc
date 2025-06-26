package main

import (
	"testing"
)

func testLiteralNode(t *testing.T, actual *LiteralNode, expected Node) bool {
	expectedLiteral, ok := expected.(*LiteralNode)
	if !ok {
		t.Error("expected node is not a LiteralNode")
		return false
	}

	if expectedLiteral.Value != actual.Value {
		t.Errorf("LiteralNode values are different, expected=%b, actual=%b", expectedLiteral.Value, actual.Value)
		return false
	}
	return true
}

func testStartNode(t *testing.T, actual *StarNode, expected Node) bool {
	expectedStar, ok := expected.(*StarNode)
	if !ok {
		t.Error("expected node is not a StarNode")
		return false
	}

	return testNode(t, actual.Child, expectedStar.Child)
}

func testDotNode(t *testing.T, actual *DotNode, expected Node) bool {
	_, ok := expected.(*DotNode)
	if !ok {
		t.Error("expected node is not a DotNode")
		return false
	}
	return true
}

func testNode(t *testing.T, actual Node, expected Node) bool {
	var result bool
	switch v := actual.(type) {
	case *LiteralNode:
		result = testLiteralNode(t, actual.(*LiteralNode), expected)
	case *StarNode:
		result = testStartNode(t, actual.(*StarNode), expected)
	case *DotNode:
		result = testDotNode(t, actual.(*DotNode), expected)
	default:
		t.Fatalf("unknown type %T", v)
	}
	if result != true {
		t.Fatal("expected node is different than actual node")
		return false
	}
	return true
}

func testSequenceNode(t *testing.T, actual *SequenceNode, expected *SequenceNode) bool {
	if len(actual.Children) != len(expected.Children) {
		t.Fatalf("expected %d children in SequenceNode, got %d", len(expected.Children), len(actual.Children))
		return false
	}
	for i, child := range actual.Children {
		expectedLiteral := expected.Children[i]
		if testNode(t, child, expectedLiteral) == false {
			return false
		}
	}

	return true
}

var b NodeBuilder
func TestParse(t *testing.T) {
	expected := b.Seq(
		b.Lit('p'),
		b.Lit('a'),
		b.Star(b.Dot()),
		b.Lit('t'),
	);
	l := New("pa.*t")
	parser := NewParser(l)
	node := parser.Ast()

	ok := testSequenceNode(t, node.(*SequenceNode), expected)
	if !ok {
		t.Fatalf("expected SequenceNode, got %T", node)
	}
}
