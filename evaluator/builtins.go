package evaluator

import (
	"Cmicro-Compiler/object"
)

/**
 * @Description: 内置函数
 */

var builtins = map[string]*object.Builtin{
	//"len": { //长度函数
	//	Fn: func(args ...object.Object) object.Object {
	//		if len(args) != 1 {
	//			return newError("wrong number of arguments. got=%d, want=1", len(args))
	//		}
	//
	//		switch arg := args[0].(type) {
	//		case *object.String:
	//			return &object.Integer{Value: int64(len(arg.Value))}
	//		case *object.Array:
	//			return &object.Integer{Value: int64(len(arg.Elements))}
	//		default:
	//			return newError("argument to `len` not supported, got %s", args[0].Type())
	//		}
	//	},
	//},
	//"println": { //打印函数 换行
	//	Fn: func(args ...object.Object) object.Object {
	//		for _, arg := range args {
	//			fmt.Printf("%v\n", arg.Inspect())
	//		}
	//		return nil
	//	},
	//},
	//"print": { //打印函数
	//	Fn: func(args ...object.Object) object.Object {
	//		for _, arg := range args {
	//			fmt.Printf("%v", arg.Inspect())
	//		}
	//		return nil
	//	},
	//},
	//"input": { //输入函数
	//	Fn: func(args ...object.Object) object.Object {
	//		var input string
	//		_, err := fmt.Scanln(&input)
	//		if err != nil {
	//			return newError("argument to `input` not supported, got %s", args[0].Type())
	//		}
	//		return &object.String{Value: input}
	//	},
	//},
	//"first": { //取数组第一个元素
	//	Fn: func(args ...object.Object) object.Object {
	//		if len(args) != 1 {
	//			return newError("wrong number of arguments. got=%d, want=1", len(args))
	//		}
	//		if args[0].Type() != object.ARRAY_OBJ {
	//			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	//		}
	//
	//		arr := args[0].(*object.Array)
	//		if len(arr.Elements) > 0 {
	//			return arr.Elements[0]
	//		}
	//
	//		return NULL
	//	},
	//},
	//"last": { //取数组最后一个元素
	//	Fn: func(args ...object.Object) object.Object {
	//		if len(args) != 1 {
	//			return newError("wrong number of arguments. got=%d, want=1", len(args))
	//		}
	//
	//		if args[0].Type() != object.ARRAY_OBJ {
	//			return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
	//		}
	//
	//		arr := args[0].(*object.Array)
	//		length := len(arr.Elements)
	//		if length > 0 {
	//			return arr.Elements[length-1]
	//		}
	//
	//		return NULL
	//	},
	//},
	//"rest": { //取数组除最后一个元素
	//	Fn: func(args ...object.Object) object.Object {
	//		if len(args) != 1 {
	//			return newError("wrong number of arguments. got=%d, want=1", len(args))
	//		}
	//		if args[0].Type() != object.ARRAY_OBJ {
	//			return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	//		}
	//		arr := args[0].(*object.Array)
	//		length := len(arr.Elements)
	//		if length > 0 {
	//			newElements := make([]object.Object, length-1, length-1)
	//			copy(newElements, arr.Elements[1:length])
	//			return &object.Array{Elements: newElements}
	//		}
	//
	//		return NULL
	//	},
	//},
	//"push": { //向数组中添加元素
	//	Fn: func(args ...object.Object) object.Object {
	//		if len(args) != 2 {
	//			return newError("wrong number of arguments. got=%d, want=2", len(args))
	//		}
	//
	//		if args[0].Type() != object.ARRAY_OBJ {
	//			return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	//		}
	//
	//		arr := args[0].(*object.Array)
	//		length := len(arr.Elements)
	//
	//		newElements := make([]object.Object, length+1)
	//		copy(newElements, arr.Elements)
	//		newElements[length] = args[1]
	//
	//		return &object.Array{Elements: newElements}
	//	},
	//},
	object.BuiltinFuncNameLen:   object.GetBuiltinByName(object.BuiltinFuncNameLen),
	object.BuiltinFuncNamePuts:  object.GetBuiltinByName(object.BuiltinFuncNamePuts),
	object.BuiltinFuncNameFirst: object.GetBuiltinByName(object.BuiltinFuncNameFirst),
	object.BuiltinFuncNameLast:  object.GetBuiltinByName(object.BuiltinFuncNameLast),
	object.BuiltinFuncNameRest:  object.GetBuiltinByName(object.BuiltinFuncNameRest),
	object.BuiltinFuncNamePush:  object.GetBuiltinByName(object.BuiltinFuncNamePush),
}
