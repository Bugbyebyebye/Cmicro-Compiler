package object

import (
	"Cmicro-Compiler/ast"
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

// 数据类型
const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_VALUE"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer 整数
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// Boolean 布尔
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// Null 空
type Null struct {
}

func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// ReturnValue 返回值
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}

// Error 错误
type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("}")

	return out.String()
}
