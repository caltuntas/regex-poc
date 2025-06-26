package main

type Parser struct {
	l *Lexer
	currentToken Token
	nextToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.currentToken = l.NextToken()
	return p
}

func (p *Parser) Ast() Node {
	return p.parseExpression()
}

func (p *Parser) parseExpression() Node {
	var node Node
	sequence := &SequenceNode{}
	node = sequence

	for p.l.HasMore() {
		term := p.parseTerm()

		if term != nil {
			sequence.Children = append(sequence.Children, term)
		}
	}
	return node
}

func (p *Parser) parseTerm() Node {
	factor := p.parseFactor()
	if p.l.NextChar() == '*' {
		p.readNextToken()
		p.readNextToken()
		star := &StarNode{}
		star.Child = factor
		return star
	} 
	p.readNextToken()
	return factor
	
}

func (p *Parser) parseFactor() Node {
	var node Node
	if p.currentToken.Type == DOT {
		node = &DotNode{}
	}else if p.currentToken.Type== LITERAL {
		node = &LiteralNode{Value: p.currentToken.Value[0]}
	}
	return node
}

func (p *Parser) isNextToken(t TokenType) bool {
	return p.nextToken.Type == t
}

func (p *Parser) readNextToken() {
	//p.currentToken = p.nextToken
	p.currentToken = p.l.NextToken()
}
