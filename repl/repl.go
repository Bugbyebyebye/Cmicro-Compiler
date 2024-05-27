package repl

import (
	"Cmicro-Compiler/compiler"
	"Cmicro-Compiler/lexer"
	"Cmicro-Compiler/object"
	"Cmicro-Compiler/parser"
	"Cmicro-Compiler/vm"
	"bufio"
	"fmt"
	"io"
)

/**
 * @Description: repl 交互式环境
 */

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// 切换为虚拟机
	//env := object.NewEnvironment()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		//将读取到的字符串 转换为token
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		//求值
		//evaluated := evaluator.Eval(program, env)
		//if evaluated != nil {
		//	io.WriteString(out, evaluated.Inspect())
		//	io.WriteString(out, "\n")
		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "编译失败:\n %s\n", err)
			continue
		}

		//machine := vm.New(comp.Bytecode())
		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "执行字节码失败:\n %s\n", err)
			continue
		}

		////stackTop := machine.StackTop()
		//io.WriteString(out, stackTop.Inspect())
		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

// 打印错误信息
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
