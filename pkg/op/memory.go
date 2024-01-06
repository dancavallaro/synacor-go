package op

import "dancavallaro.com/synacor-go/pkg/memory"

func Set(r *memory.Registers, args []uint16) {
	reg, b := memory.RegNum(args[0]), memory.ReadVal(r, args[1])
	r.GP[reg] = b
}
