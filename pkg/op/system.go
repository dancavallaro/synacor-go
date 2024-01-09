package op

import (
	"dancavallaro.com/synacor-go/pkg/env"
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
	fmt.Fprint(env.Config.Output, string(rune(args[0])))
}

func In(m *memory.Memory, args []uint16) {
	a := memory.RegNum(args[0])
	ch, err := env.Config.ReadChar()
	if err != nil {
		panic(err)
	}
	log.Printf("IN read '%s' (%d) from stdin\n", string(rune(ch)), rune(ch))
	m.GP[a] = ch
}
