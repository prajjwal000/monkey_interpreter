package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is your lovely monkey language\n", user.Username)
	fmt.Printf("Type commands my Lord!\n")
	repl.Start(os.Stdin, os.Stdout)
}
