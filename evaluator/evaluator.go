package evaluator

import (
	"Cmicro-Compiler/ast"
	"Cmicro-Compiler/object"
	"fmt"
)

/**
 * @Description: 求值器
 */

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval 节点求值
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program: // 程序嵌套
		return evalProgram(node, env)
	case *ast.ExpressionStatement: // 表达式
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral: // 整型
		return &object.Integer{Value: node.Value}
	case *ast.Boolean: // 布尔值
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression: // 前缀运算符
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		result := evalPrefixExpression(node.Operator, right)
		env.Set(node.Right.String(), result)
		return result
	case *ast.InfixExpression: // 中缀运算符
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement: // 语句块
		return evalBlockStatement(node, env)
	case *ast.IfExpression: // if条件
		return evalIfExpression(node, env)
	case *ast.ReturnStatement: // 返回
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement: // 变量初始化 let
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.AssignStatement: //变量赋值
		return evalAssignStatement(node, env)
	case *ast.ForExpression: // for循环
		return evalForExpression(node, env)
	case *ast.Identifier: // 变量
		return evalIdentifier(node, env)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.FunctionLiteral: // 函数
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression: // 函数调用
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.ArrayLiteral: // 数组
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression: // 数组索引
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral: // 哈希表
		return evalHashLiteral(node, env)
	}
	return nil
}
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	// 遍历程序中的语句，处理嵌套语句块
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	// 遍历语句块中的语句
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	// 从环境变量中获取 内置函数
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	// 遍历表达式列表，在当前环境的上下文中求值，如果遇到错误，就停止求值并返回错误
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}

}
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	// 创建函数环境
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}
func unwrapReturnValue(obj object.Object) object.Object {
	// 如果返回值是ReturnValue类型，则返回其值，否则返回obj本身
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// evalIndexExpression 索引表达式匹配求值方法
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
func evalArrayIndexExpression(array, index object.Object) object.Object {
	// 数组索引
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

// evalHashLiteral 哈希表匹配求值方法
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}
func evalHashIndexExpression(hash, index object.Object) object.Object {
	// 哈希表索引
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

// nativeBoolToBooleanObject 将bool转换为Boolean对象
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// evalPrefixExpression 前缀表达式匹配求值方法
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "++":
		return evalIncrementPrefixOperatorExpression(right)
	case "--":
		return evalDecrementPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}
func evalBangOperatorExpression(right object.Object) object.Object {
	// 对逻辑值取反
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// evalMinusPrefixOperatorExpression 前缀表达式求值
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
func evalIncrementPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: ++%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: value + 1}
}
func evalDecrementPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: --%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: value - 1}
}

// evalInfixExpression 中缀表达式求值
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression 整型运算
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIfExpression If表达式求值
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return newError("else branch must be present when if condition is false")
	}
}
func isTruthy(obj object.Object) bool {
	// 判断是否为逻辑值
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// evalForStatement For语句求值
func evalForExpression(fs *ast.ForExpression, env *object.Environment) object.Object {
	var result object.Object

	if fs.Init != nil {
		Eval(fs.Init, env)
	}
	for {
		if fs.Condition != nil {
			condition := Eval(fs.Condition, env)
			if isError(condition) {
				return newError("condition must be present in for loop")
			}
			if !isTruthy(condition) {
				// 跳出循环
				break
			}
		} else {
			return newError("condition must be present in for loop")
		}

		evaluated := Eval(fs.Body, env)
		if isError(evaluated) {
			return newError("for loop body error")
		} else {
			result = evaluated
		}
		if fs.Post != nil {
			Eval(fs.Post, env)
		}
	}

	return result
}

// evalAssignStatement 赋值语句求值
func evalAssignStatement(as *ast.AssignStatement, env *object.Environment) object.Object {
	value := Eval(as.Value, env)
	if isError(value) {
		return value
	}

	name := as.Name.Value
	if _, ok := env.Get(name); ok {
		env.Set(name, value)
		return value
	}

	return newError("identifier not found: " + name)
}

// evalStringInfixExpression 字符串拼接运算
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

// 错误处理
func newError(format string, a ...interface{}) object.Object {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
func isError(obj object.Object) bool {
	// 判断是否为错误对象
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// 内置函数
func evalBuiltin(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}
