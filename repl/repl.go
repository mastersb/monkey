package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mastersb/monkey/evaluator"
	"github.com/mastersb/monkey/lexer"
	"github.com/mastersb/monkey/object"
	"github.com/mastersb/monkey/parser"
)

const PROMPT = ">> "
const CONT_PROMPT = "... "

const MONKEY_FACE = `            __,__
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

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	var input strings.Builder
	prompt := PROMPT

	for {
		fmt.Fprintf(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		input.WriteString(line)
		input.WriteString("\n")

		source := input.String()

		if !isBraceBalanced(source) {
			prompt = CONT_PROMPT
			continue
		}

		l := lexer.New(source)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			prompt = PROMPT
			input.Reset()
			continue
		}

		input.Reset()
		prompt = PROMPT

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func isBraceBalanced(source string) bool {
	braces := 0
	brackets := 0
	parens := 0

	inString := false
	escaped := false

	for _, ch := range source {
		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch ch {
		case '{':
			braces++
		case '}':
			braces--
		case '[':
			brackets++
		case ']':
			brackets--
		case '(':
			parens++
		case ')':
			parens--
		}
	}

	return braces == 0 && brackets == 0 && parens == 0
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
