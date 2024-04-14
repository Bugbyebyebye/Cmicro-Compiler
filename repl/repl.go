package repl

import (
	"Cmicro-Compiler/evaluator"
	"Cmicro-Compiler/lexer"
	"Cmicro-Compiler/object"
	"Cmicro-Compiler/parser"
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
	env := object.NewEnvironment()

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
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

// 打印错误信息
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
