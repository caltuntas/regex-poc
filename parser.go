package main

type Parser struct {
	l            *Lexer
	currentToken Token
	nextToken    Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.readNextToken()
	p.readNextToken()
	return p
}

func (p *Parser) Ast() Node {
	return p.parseExpression()
}

func (p *Parser) parseExpression() Node {
	var node Node
	sequence := &SequenceNode{}
	node = sequence

	for p.currentToken.Type != EOF {
		term := p.parseTerm()

		if term != nil {
			sequence.Children = append(sequence.Children, term)
		}
	}
	return node
}

func (p *Parser) parseTerm() Node {
	factor := p.parseFactor()
	if p.nextToken.Type == STAR {
		star := &StarNode{}
		star.Child = factor
		factor = star
		p.readNextToken()
	}
	p.readNextToken()
	return factor

}

func (p *Parser) parseFactor() Node {
	var node Node
	if p.currentToken.Type == DOT {
		node = &MetaCharacterNode{Value: "."}
	} else if p.currentToken.Type == LITERAL {
		node = &LiteralNode{Value: p.currentToken.Value[0]}
	} else if p.currentToken.Type == ESCAPE {
		p.readNextToken()
		if p.currentToken.Value == "s" {
			return &MetaCharacterNode{Value: WHITESPACE}
		}
	} else if p.currentToken.Type == LBRACKET {
		p.readNextToken()
		charList := &CharList{}
		for p.currentToken.Type != RBRACKET {
			char := p.parseFactor()
			charList.Chars = append(charList.Chars, char.(CharacterNode))
			p.readNextToken()
		}
		return charList
	}
	return node
}

func (p *Parser) readNextToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.l.NextToken()
}
