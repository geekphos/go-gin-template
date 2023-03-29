package main

import (
	"fmt"
	"os"
	"phos.cc/yoo/internal/yoo"
)
import _ "go.uber.org/automaxprocs"

func main() {
	command := yoo.NewYooCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
