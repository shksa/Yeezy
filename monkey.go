package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/shksa/monkey/object"

	"github.com/shksa/monkey/evaluator"
	"github.com/shksa/monkey/parser"

	"github.com/shksa/monkey/lexer"
)

var (
	fileNamePtr = flag.String("file", "", "name of file to interpret")
)

// PROMPT is the prompt message for the repl.
const PROMPT = ">> "

// PEPE is displayed before error messages
const PEPE = `
⢠⠤⣤⠀⠤⡤⠄⢠⡤⢤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ 
⢸⠲⣏⠀⢀⡇⠀⢸⡗⠚⢀⣤⣶⣾⣿⣷⣶⣤⣄⠀⠀⣀⣤⣤⣴⣦⣤⡀⠀⠀⠀⠀⠀⠀⠀ 
⠈⠀⠈⠀⠉⠉⠁⠈⠁⣴⣿⣿⣿⡿⠿⣛⣛⠻⠿⣧⢻⣿⣿⣿⣿⣿⣿⣿⣄⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣿⣫⣵⣾⣿⣿⣿⡿⠷⠦⠔⣶⣶⣶⣶⣶⠶⠶⠶⠤⡀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⣿⣿⠿⠛⢁⣀⣌⣿⣿⣷⣶⣈⠿⣒⣒⣭⣭⣭⣭⣑⣒⠄⠀⠀
⠀⠀⠀⠀⠀⠀⣠⡎⣾⣿⣿⣿⣿⢫⣡⡥⠶⠿⣛⠛⠋⠳⢶⣶⣾⣜⣫⣭⣷⠖⡁⠀⠐⢶⣯⡆⠀ 
⠀⠀⠀⣰⣿⣷⣿⣿⣿⣿⣿⣷⣖⢟⡻⢿⠃⢸⠱⠶⠀⠿⠟⡻⠿⣿⡏⠀⠅⠛⠀⣘⠟⠁⠀ ⠀
⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣮⣥⣤⣴⣤⣦⠄⣠⣾⣿⡻⠿⠾⠿⠿⠟⠛⠁⠀⠀⠀ ⠀
⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⣭⣶⣿⣿⣿⣿⣿⣷⣿⣿⣿⣧⡀⠀⠀⠀⠀ ⠀
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡄⠀⠀⠀ ⠀
⢿⣿⣿⣿⣿⣿⣿⣿⣿⡿⢩⡤⠶⠶⠶⠦⠬⣉⣛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⣋⣡⠀⠀⠀ ⠀
⠘⣿⣿⣿⣿⣿⣿⣟⢿⣧⣙⠓⢒⣚⡛⠳⠶⠤⢬⣉⣉⣉⣉⣉⣉⣉⣉⣉⣉⡄⠀⠀⠀⠀ ⠀
⠀⠈⠻⢿⣿⣿⣿⣿⣶⣽⣿⣿⣿⣿⣿⣿⣷⣶⣶⣶⣤⣤⣤⣤⣤⣤⡥⠄⠀⠀⠀⠀⠀⠀ 
⠀⠀⠀⠀⠐⠒⠭⢭⣛⣛⡻⠿⠿⠿⠿⣿⣿⣿⣿⣿⠿

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
			// io.WriteString(out, program.String())
			// io.WriteString(out, "\n")
		}
	}
}

func runREPL() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	start(os.Stdin, os.Stdout)
}

func runProgramFile(fileName string) {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	env := object.NewEnvironment()
	l := lexer.New(string(fileContent))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors) != 0 {
		printParseErrors(p.Errors)
		return
	}

	evaluated := evaluator.Eval(program, env)
	fmt.Println(evaluated.Inspect())
}

func main() {
	flag.Parse()
	fileName := *fileNamePtr
	if fileName == "" {
		runREPL()
	} else {
		extOfFileName := filepath.Ext(fileName)
		if extOfFileName != ".mky" {
			fmt.Printf("Invalid file extention: Not a monkey file. got=%q \n", extOfFileName)
			return
		}
		runProgramFile(fileName)
	}
}

func printParseErrors(errors []string) {
	fmt.Println(PEPE, "whoops! PEPE died after seeing your shit code!\n", "parse errors:")
	for _, errMsg := range errors {
		fmt.Println("\t", errMsg)
	}
}
