package lexer

import (
	"testing"
	"monkey/token"
)

func TestNextToken(t *testing.T){
    input := `=+,;(){}`
    tests := []struct{
        expectedType token.TokenType
        expectedLiteral string
    }{
        {token.ASSIGN, "="},
        {token.PLUS, "+"},
        {token.COMMA, ","},
        {token.SEMICOLON, ";"},
        {token.LPAREN, "("},
        {token.RPAREN, ")"},
        {token.LCURLY, "{"},
        {token.RCURLY, "}"},
    }
    l := New(input)

    for i,tt := range tests {
        tok := l.NextToken()
        if tok.Type != tt.expectedType{
            t.Fatalf("tests[%d] - tokentype wrong expectedType=%q, got=%q",i,tt.expectedType,tok.Type)
        }
        if tok.Literal != tt.expectedLiteral{
            t.Fatalf("tests[%d] - tokentype wrong expectedLiteral=%q, got=%q",i,tt.expectedLiteral,tok.Literal)
        }
    }
}
