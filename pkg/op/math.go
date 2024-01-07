package op

import "dancavallaro.com/synacor-go/pkg/memory"

func Add(m *memory.Memory, args []uint16) {
	a, b, c := memory.RegNum(args[0]), memory.ReadVal(m, args[1]), memory.ReadVal(m, args[2])
	m.GP[a] = (b + c) % memory.Modulus
}

func Eq(m *memory.Memory, args []uint16) {
	a, b, c := memory.RegNum(args[0]), memory.ReadVal(m, args[1]), memory.ReadVal(m, args[2])
	if b == c {
		m.GP[a] = 1
	} else {
		m.GP[a] = 0
	}
}
