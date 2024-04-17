package evaluator

import (
	"Cmicro-Compiler/object"
	"fmt"
)

/**
 * @Description: 内置函数
 */

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			//case *object.Array:
			//return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			return args[0]
		},
	},
	"input": {
		Fn: func(args ...object.Object) object.Object {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return newError("argument to `input` not supported, got %s", args[0].Type())
			}
			return &object.String{Value: input}
		},
	},
}
