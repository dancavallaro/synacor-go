package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
)

func jump(r *memory.Registers, addr uint16) {
	byteAddr := addr * 2 // Addresses are word-indexed, but "memory" is byte-indexed
	r.PC = int(byteAddr)
}

func Jmp(r *memory.Registers, args []uint16) {
	jump(r, args[0])
}

func Jt(r *memory.Registers, args []uint16) {
	// TODO: need to support the args being registers
	a, b := args[0], args[1]
	if a != 0 {
		jump(r, b)
	}
}

func Jf(r *memory.Registers, args []uint16) {
	// TODO: need to support the args being registers
	a, b := args[0], args[1]
	if a == 0 {
		jump(r, b)
	}
}
