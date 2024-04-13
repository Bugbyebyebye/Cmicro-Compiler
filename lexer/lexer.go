package lexer

import "Cmicro-Compiler/token"

/**
 * @File: lexer
 * @Description:词法分析器
 */

type Lexer struct {
	input        string //输入的代码字符串
	position     int    //指向当前字符
	readPosition int    //指向下一个字符
	ch           byte   //当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 读取input下一个字符
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1 //readPosition 始终指向下一个字符
}

// NextToken 用于获取下一个token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// 只要l.ch不是前面定义的可以识别的字符，就检查是不是标识符
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			return tok //提前退出，因为readIdentifier中继续调用了l.readChar()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// 创建token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 读取标识符（变量）
// 并前移词法分析器的扫描位置，直到遇见非字母字符
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 判断是否为字母
// （决定了语言能够处理的语言形式，即变量的命名规则）
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
