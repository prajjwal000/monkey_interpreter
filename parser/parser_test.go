package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
    input := `
    let x  5;
    let  = 10;
    let 838383;
    `
    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    /*
    sss := program.Statements[0].(*ast.LetStatement)
    ss := sss.BindL
    t.Logf("kate \n %+v",ss)
    */

    checkParserErrors(t, p)

    if program == nil {
        t.Fatalf("Parser returned nil")
    }
    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements has %d members",len(program.Statements))
    }

    tests := []struct{
        expectedIdentifier string
    }{
        {"x"},
        {"y"},
        {"foobar"},
    }
    
    for i,tt := range tests {
        stmt := program.Statements[i]
        if !testLetStatements(t, stmt , tt.expectedIdentifier) {
            return 
        }
    }
}

func testLetStatements(t *testing.T, s ast.Statement, bindl string) bool {
    if s.TokenLiteral() != "let" {
        t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
    }

    letStmt, ok := s.(*ast.LetStatement)
    if !ok {
        t.Errorf("s not *ast.LetStatement. got=%T", s)
        return false
    }

    if letStmt.BindL.Value != bindl {
        t.Errorf("letStmt.Name.Value not '%s'. got=%s", bindl, letStmt.BindL.TokenLiteral())
        return false
    }

    if letStmt.BindL.TokenLiteral() != bindl {
        t.Errorf("s.Name not '%s'. got=%s", bindl, letStmt.BindL.TokenLiteral())
        return false
    }

    return true 
}

func checkParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors()

    if len(errors) == 0 {
        return 
    }

    t.Errorf("parser has %d errors", len(errors))

    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }
    t.FailNow()
}

