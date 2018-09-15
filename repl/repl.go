package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"

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

func start() {
	scanner := bufio.NewScanner(os.Stdin)
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

		fmt.Println(program.String())
	}
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	start()
}

func printParseErrors(errors []string) {
	fmt.Println(MONKEYFACE, "whoops! we ran into some monkey business!\n", "parse errors:")
	for _, errMsg := range errors {
		fmt.Println("\t", errMsg)
	}
}
