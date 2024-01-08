package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"fmt"
	"io"
	"log"
	"os"
)

var Output io.Writer = os.Stdout

func Halt(_ *memory.Memory, _ []uint16) {
	log.Println()
	log.Println("Execution halted.")
	os.Exit(0)
}

func Noop(_ *memory.Memory, _ []uint16) {}

func Out(_ *memory.Memory, args []uint16) {
	fmt.Fprint(Output, string(rune(args[0])))
}

func In(m *memory.Memory, args []uint16) {
	//a := memory.RegNum(args[0])
	//m.GP[a] = util.ReadChar()
}
