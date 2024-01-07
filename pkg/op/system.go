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

func In(m *memory.Memory, args []uint16) {
	a := memory.RegNum(args[0])
	m.GP[a] = readChar()
}

func readChar() uint16 {
	b := make([]byte, 1)
	n, err := os.Stdin.Read(b)
	if n != 1 || err != nil {
		panic("error reading from stdin!")
	}
	return uint16(b[0])
}
