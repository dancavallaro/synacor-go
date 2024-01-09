package op

import (
	"dancavallaro.com/synacor-go/internal/util"
	"dancavallaro.com/synacor-go/pkg/memory"
	"fmt"
	"io"
	"log"
	"os"
)

// TODO: add halt
// TODO: add log
// TODO: move this to vm.go?
var Environment = struct {
	Output   io.Writer
	ReadChar func() (uint16, error)
}{
	os.Stdin,
	util.ReadChar,
}

func Halt(_ *memory.Memory, _ []uint16) {
	log.Println()
	log.Println("Execution halted.")
	os.Exit(0)
}

func Noop(_ *memory.Memory, _ []uint16) {}

func Out(_ *memory.Memory, args []uint16) {
	fmt.Fprint(Environment.Output, string(rune(args[0])))
}

func In(m *memory.Memory, args []uint16) {
	a := memory.RegNum(args[0])
	ch, err := Environment.ReadChar()
	if err != nil {
		panic(err)
	}
	m.GP[a] = ch
}
