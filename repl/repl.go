package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
    env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
        if line == "exit" {
            return
        }
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
        evaluated := eval.Eval(program,env)
        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out,"\n")
        }
	}
}

const MONKEY_FACE = `
_______AAAA____       ____AAAA________
       VVVV               VVVV        
       (__)               (__)
        \ \               / /
         \ \   \\|||//   / /
          > \   _   _   / <
 hang      > \ / \ / \ / <
  in        > \\_o_o_// <
  there...   > ( (_) ) <
              >|     |<
             / |\___/| \
             / (_____) \
             /         \
              /   o   \
               ) ___ (   
              / /   \ \  
             ( /     \ )
             ><       ><
            ///\     /\\\
            '''       '''
`
func printParserErrors(out io.Writer, errors []string) {
io.WriteString(out, MONKEY_FACE)
io.WriteString(out, "Woops! We ran into some monkey business here!\n")
io.WriteString(out, " parser errors:\n")
for _, msg := range errors {
io.WriteString(out, "\t"+msg+"\n")
}
}
