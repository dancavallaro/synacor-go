package op

import (
	"fmt"
	"log"
	"os"
)

func Halt(_ []int16) {
	log.Println()
	log.Println("Execution halted.")
	os.Exit(0)
}

func Noop(_ []int16) {}

func Out(in []int16) {
	fmt.Print(string(in[0]))
}
