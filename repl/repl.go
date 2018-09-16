package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"

	"github.com/shksa/monkey/object"

	"github.com/shksa/monkey/evaluator"
	"github.com/shksa/monkey/parser"

	"github.com/shksa/monkey/lexer"
)

// PROMPT is the prompt message for the repl.
const PROMPT = ">> "

// MONKEYFACE is displayed before error messages
const MONKEYFACE = `            
						__,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		didScan := scanner.Scan()
		if !didScan {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors) != 0 {
			printParseErrors(p.Errors)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
		}
	}
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	start(os.Stdin, os.Stdout)
}

func printParseErrors(errors []string) {
	fmt.Println(MONKEYFACE, "whoops! we ran into some monkey business!\n", "parse errors:")
	for _, errMsg := range errors {
		fmt.Println("\t", errMsg)
	}
}
