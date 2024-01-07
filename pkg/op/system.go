package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"fmt"
	"log"
	"os"
)

func Halt(_ *memory.Memory, _ []uint16) {
	log.Println()
	log.Println("Execution halted.")
	os.Exit(0)
}

func Noop(_ *memory.Memory, _ []uint16) {}

func Out(_ *memory.Memory, args []uint16) {
	fmt.Print(string(rune(args[0])))
}
