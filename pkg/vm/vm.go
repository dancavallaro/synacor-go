package vm

import (
	"dancavallaro.com/synacor-go/pkg/op"
	"errors"
	"fmt"
	"log"
)

type opRef struct {
	opcode  int16
	nArgs   int
	execute func([]int16)
}

var ops = []opRef{
	{0, 0, op.Halt},
	{19, 1, op.Out},
	{21, 0, op.Noop},
}
var opRefs = map[int16]opRef{}

func init() {
	for _, o := range ops {
		opRefs[o.opcode] = o
	}
}

func readWord(bin []byte, address int) int16 {
	return (int16(bin[address+1]) << 8) + int16(bin[address])
}

func Execute(bin []byte, trace bool) error {
	for i := 0; i < len(bin)-1; i = i + 2 {
		w := readWord(bin, i)
		o, ok := opRefs[w]
		if !ok {
			return errors.New(fmt.Sprintf("invalid opcode %d", w))
		}
		if trace {
			log.Printf("Executing opcode %d\n", o.opcode)
		}

		var input []int16
		for arg := 1; arg <= o.nArgs; arg++ {
			w := readWord(bin, i+1+arg)
			input = append(input, w)
			i += 2
		}
		o.execute(input)
	}
	return nil
}
