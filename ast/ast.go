package ast

import (
	"bytes"
    "strings"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token
	Bind  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Bind.String() + " = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) expressionNode()      {}
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral())
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	On       Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.On.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(" + ie.Left.String() + " ")
	out.WriteString(ie.Operator + " ")
	out.WriteString(ie.Right.String() + ")")

	return out.String()
}

type Boolean struct {
    Token token.Token
    Value bool
}

func (b *Boolean) expressionNode(){}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

type IfExpression struct {
    Token token.Token
    Condition Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}

func (if_ *IfExpression) expressionNode() {}
func (if_ *IfExpression) TokenLiteral() string { return if_.Token.Literal }
func (if_ *IfExpression) String() string {
    var out bytes.Buffer

    out.WriteString("if ")
    out.WriteString(if_.Condition.String() + " ")
    out.WriteString(if_.Consequence.String())

    if if_.Alternative != nil {
        out.WriteString("else " + if_.Alternative.String())
    }
    return out.String()
}

type BlockStatement struct {
    Token token.Token
    Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    for _, s := range bs.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

type FunctionLiteral struct {
    Token token.Token
    Parameters []*Identifier
    Body *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
    var out bytes.Buffer

    params := []string{}
    for _,p := range fl.Parameters {
        params = append(params, p.String())
    }
    out.WriteString(fl.TokenLiteral() + "(")
    out.WriteString(strings.Join(params,","))
    out.WriteString(")" + fl.Body.String())

    return out.String()
}

type CallExpression struct {
    Token token.Token
    Function Expression
    Arguments []Expression
}

func (c *CallExpression) expressionNode() {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
    var out bytes.Buffer

    args := []string{}
    for _, a := range c.Arguments {
        args = append(args, a.String())
    }

    out.WriteString(c.Function.String())
    out.WriteString("(" + strings.Join(args,", ") + ")")

    return out.String()
}

type StringLiteral struct {
    Token token.Token
    Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string { return sl.Token.Literal }

type ArrayLiteral struct {
    Token token.Token
    Elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
    var out bytes.Buffer

    out.WriteString("[")
    s := []string{}
    for _, e := range a.Elements {
        s = append(s, e.String())
    }
    out.WriteString(strings.Join(s,", "))
    out.WriteString("]")
    
    return out.String()
}

type IndexExpression struct {
    Token token.Token
    Left Expression
    Index Expression
}

func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(" + i.Left.String() + "[" + i.Index.String() + "])")

    return out.String()
}

type HashLiteral struct {
    Token token.Token
    Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode() {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
    var out bytes.Buffer

    pairs := []string{}
    for key, value := range h.Pairs {
        pairs = append(pairs, key.String() + ":" + value.String())
    }
    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ","))
    out.WriteString("}")

    return out.String()
}
