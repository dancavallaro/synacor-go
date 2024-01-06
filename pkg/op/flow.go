package op

import "dancavallaro.com/synacor-go/pkg/memory"

func Jmp(args []int16, r *memory.Registers) {
	addr := args[0] * 2 // Addresses are word-indexed, but my "memory" is byte-indexed
	r.PC = int(addr)
}
