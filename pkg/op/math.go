package op

import "dancavallaro.com/synacor-go/pkg/memory"

func Add(r *memory.Registers, args []uint16) {
	a, b, c := memory.RegNum(args[0]), memory.ReadVal(r, args[1]), memory.ReadVal(r, args[2])
	r.GP[a] = (b + c) % memory.Modulus
}
