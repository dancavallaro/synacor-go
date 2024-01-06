package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"fmt"
	"log"
	"os"
)

func Halt(_ []int16, _ *memory.Registers) {
	log.Println()
	log.Println("Execution halted.")
	os.Exit(0)
}

func Noop(_ []int16, _ *memory.Registers) {}

func Out(args []int16, _ *memory.Registers) {
	fmt.Print(string(args[0]))
}
