package parser

import "github.com/kijimaD/gogo/token"

func (p *Parser) peekTokenIs(expect token.TokenType) bool {
	return p.peekToken.Type == expect
}
