package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
    let x = 5;
    let y = 10;
    let fuck = 838383;
    `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	/*
	   sss := program.Statements[0].(*ast.LetStatement)
	   ss := sss.Bind
	   t.Logf("kate \n %+v",ss)
	*/

	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("Parser returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements has %d members", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"fuck"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatements(t, stmt, tt.expectedIdentifier) {
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

	if letStmt.Bind.Value != bindl {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", bindl, letStmt.Bind.TokenLiteral())
		return false
	}

	if letStmt.Bind.TokenLiteral() != bindl {
		t.Errorf("s.Name not '%s'. got=%s", bindl, letStmt.Bind.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {

	input := `
    return 5;
    return 10;
    return 993322;
    `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foul;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	//    t.Logf("ast = %s", program.String())
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program statement number!=1, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not ExpressionStatement, got=%T", program.Statements[0])
	}

	i, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not identifier, got %T", stmt.Expression)
	}

	if i.Value != "foul" {
		t.Errorf("identifier value error, got %s", i.Value)
	}

	if i.TokenLiteral() != "foul" {
		t.Errorf("identifier tokenliteral error, got %s", i.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	//    t.Logf("ast = %s", program.String())
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program statement number!=1, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not ExpressionStatement, got=%T", program.Statements[0])
	}

	i, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not identifier, got %T", stmt.Expression)
	}

	if i.Value != 5 {
		t.Errorf("identifier value error, got %d", i.Value)
	}

	if i.TokenLiteral() != "5" {
		t.Errorf("identifier tokenliteral error, got %s", i.TokenLiteral())
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		on       int64
	}{
		{"!5", "!", 5},
		{"-16", "-", 16},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.On, tt.on) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
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
