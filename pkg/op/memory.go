package op

import "dancavallaro.com/synacor-go/pkg/memory"

func Set(m *memory.Memory, args []uint16) {
	reg, b := memory.RegNum(args[0]), memory.ReadVal(m, args[1])
	m.GP[reg] = b
}

func Push(m *memory.Memory, args []uint16) {
	a := memory.ReadVal(m, args[0])
	m.Push(a)
}

func Pop(m *memory.Memory, args []uint16) {
	a := memory.RegNum(args[0])
	m.GP[a] = m.Pop()
}

func Rmem(m *memory.Memory, args []uint16) {
	a, b := memory.RegNum(args[0]), memory.ReadVal(m, args[1])
	m.GP[a] = m.Mem[b]
}

func Wmem(m *memory.Memory, args []uint16) {
	a, b := memory.ReadVal(m, args[0]), memory.ReadVal(m, args[1])
	m.Mem[a] = b
}
