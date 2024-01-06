package main

import (
	"dancavallaro.com/synacor-go/pkg/vm"
	"flag"
	"log"
	"os"
)

const required = "<required>"

var (
	binPath = flag.String("bin", required, "Path to executable (.bin)")
)

func main() {
	flag.Parse()
	if *binPath == required {
		log.Fatalln("-bin is required")
	}
	binary, err := os.ReadFile(*binPath)
	if err != nil {
		log.Fatalln(err)
	}
	vm.Execute(binary)
}
