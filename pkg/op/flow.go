package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
)

func jump(m *memory.Memory, addr uint16) {
	byteAddr := addr * 2 // Addresses are word-indexed, but "memory" is byte-indexed
	m.PC = int(byteAddr)
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
