package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" //未知词法单元
	EOF     = "EOF"     //文件结尾

	IDENT  = "IDENT" // 标识符和字面量
	INT    = "INT"
	STRING = "STRING"

	ASSIGN    = "=" // 运算符
	PLUS      = "+"
	INCREMENT = "++"
	DECREMENT = "--"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NEQ       = "!="

	COMMA     = "," // 分隔符
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION" // 关键字
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	FOR = "FOR"
)

// 语言的关键字
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,

	"for": FOR,
}

// LookupIdent  根据标识符返回对应的TokenType
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
