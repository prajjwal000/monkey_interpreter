package token

type TokenType string
type Token struct {
    Type TokenType
    Literal string
}

const (
    Illegal="ILLEGAL"
    EOF="EOF"

    //IDENTIFIERS+ LITERALS
    IDENT="IDENT"
    INT="INT"

    //OPERATORS
    ASSIGN="="
    PLUS="+"
    MINUS="-"
    
    //DELIMITERS
    COMMA=","
    SEMICOLON=";"

    LPAREN="("
    RPAREN=")"
    LCURLY="{"
    RCURLY="}"

    //KEYWORDS
    LET="LET"
    FUNC="FUNC"
)
