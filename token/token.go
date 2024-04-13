package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" //未知词法单元
	EOF     = "EOF"     //文件结尾

	IDENT = "IDENT" // 标识符和字面量
	INT   = "INT"

	ASSIGN = "=" // 运算符
	PLUS   = "+"

	COMMA     = "," // 分隔符
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION" // 关键字
	LET      = "LET"
)

// 语言的关键字
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent  根据标识符返回对应的TokenType
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
