package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/shksa/monkey/token"

	"github.com/shksa/monkey/lexer"
)

// PROMPT is the prompt message for the repl.
const PROMPT = ">> "

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

		for tok := l.NextToken(); tok != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v \n", tok)
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
	start()
}
