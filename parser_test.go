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

func testMetaCharacter(t *testing.T, actual *MetaCharacterNode, expected Node) bool {
	expectedLiteral, ok := expected.(*MetaCharacterNode)
	if !ok {
		t.Error("expected node is not a LiteralNode")
		return false
	}

	if expectedLiteral.Value != actual.Value {
		t.Errorf("LiteralNode values are different, expected=%s, actual=%s", expectedLiteral.Value, actual.Value)
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

func testCharListNode(t *testing.T, actual *CharList, expected Node) bool {
	expectedCharList, ok := expected.(*CharList)
	if !ok {
		t.Error("expected node is not a CharList")
		return false
	}

	for i := 0; i < len(actual.Chars); i++ {
		expectedChar := expectedCharList.Chars[i]
		actualChar := actual.Chars[i]
		if testNode(t, actualChar, expectedChar) == false {
			t.Errorf("CharList characters are different, expected %s, actual %s", expectedChar.String(), actualChar.String())
			return false
		}
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
	case *CharList:
		result = testCharListNode(t, actual.(*CharList), expected)
	case *MetaCharacterNode:
		result = testMetaCharacter(t, actual.(*MetaCharacterNode), expected)
	case *SequenceNode:
		result = testSequenceNode(t, actual.(*SequenceNode), expected)
	default:
		t.Fatalf("unknown type %T", v)
	}
	if result != true {
		t.Fatal("expected node is different than actual node")
		return false
	}
	return true
}

func testSequenceNode(t *testing.T, actual *SequenceNode, expectedNode Node) bool {
	expected, ok := expectedNode.(*SequenceNode)
	if !ok {
		t.Error("expected node is not a StarNode")
		return false
	}
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

func TestRegexParser(t *testing.T) {
	cases := map[string]Node{
		"p[ab]*c": b.Seq(
			b.Lit('p'),
			b.Star(b.List(b.Lit('a'), b.Lit('b'))),
			b.Lit('c'),
		),
		"p[\\sb]c": b.Seq(
			b.Lit('p'),
			b.List(b.Meta(WHITESPACE), b.Lit('b')),
			b.Lit('c'),
		),
		"pa\\s": b.Seq(
			b.Lit('p'),
			b.Lit('a'),
			b.Meta(WHITESPACE),
		),
		"pa[ab]c": b.Seq(
			b.Lit('p'),
			b.Lit('a'),
			b.List(b.Lit('a'), b.Lit('b')),
			b.Lit('c'),
		),
		"pa.*t": b.Seq(
			b.Lit('p'),
			b.Lit('a'),
			b.Star(b.Meta(DOT)),
			b.Lit('t'),
		),
		"pa {": b.Seq(
			b.Lit('p'),
			b.Lit('a'),
			b.Lit(' '),
			b.Lit('{'),
		),
	}
	for key, val := range cases {
		l := New(key)
		parser := NewParser(l)
		node := parser.Ast()
		ok := testNode(t, node.(*SequenceNode), val)
		if !ok {
			t.Fatalf("expected SequenceNode, got %T", node)
		}
	}

}
