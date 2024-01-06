package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
)

func jump(r *memory.Registers, addr uint16) {
	byteAddr := addr * 2 // Addresses are word-indexed, but "memory" is byte-indexed
	r.PC = int(byteAddr)
}

func Jmp(r *memory.Registers, args []uint16) {
	target := memory.ReadVal(r, args[0])
	jump(r, target)
}

func Jt(r *memory.Registers, args []uint16) {
	a, b := memory.ReadVal(r, args[0]), memory.ReadVal(r, args[1])
	if a != 0 {
		jump(r, b)
	}
}

func Jf(r *memory.Registers, args []uint16) {
	a, b := memory.ReadVal(r, args[0]), memory.ReadVal(r, args[1])
	if a == 0 {
		jump(r, b)
	}
}
