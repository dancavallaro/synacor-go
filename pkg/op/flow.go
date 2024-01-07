package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
)

func jump(m *memory.Memory, addr uint16) {
	m.PC = int(addr)
}

func Jmp(m *memory.Memory, args []uint16) {
	target := memory.ReadVal(m, args[0])
	jump(m, target)
}

func Jt(m *memory.Memory, args []uint16) {
	a, b := memory.ReadVal(m, args[0]), memory.ReadVal(m, args[1])
	if a != 0 {
		jump(m, b)
	}
}

func Jf(m *memory.Memory, args []uint16) {
	a, b := memory.ReadVal(m, args[0]), memory.ReadVal(m, args[1])
	if a == 0 {
		jump(m, b)
	}
}

func Call(m *memory.Memory, args []uint16) {
	a := memory.ReadVal(m, args[0])
	m.Push(uint16(m.PC))
	jump(m, a)
}

func Ret(m *memory.Memory, _ []uint16) {
	a := m.Pop()
	jump(m, a)
}
