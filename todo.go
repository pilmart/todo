package main

// imports
import (
	"todo/cli"
	"todo/web"
)

// All code split out into packages, this is now the main entrypoint
func main() {
	go cli.StartToDo()
	web.StartMux()
}
