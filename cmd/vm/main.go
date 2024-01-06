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
	trace   = flag.Bool("trace", false, "Print each opcode during execution")
	delay   = flag.Int("delay", -1, "Delay (in milliseconds) between execution of each operation")
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

	opts := vm.ExecutionOptions{
		Trace: *trace,
		Delay: *delay,
	}
	err = vm.Execute(binary, &opts)
	log.Println()
	if err != nil {
		log.Fatalf("Execution aborted: %s\n", err)
	}
	log.Println("Program completed execution normally.")
}
