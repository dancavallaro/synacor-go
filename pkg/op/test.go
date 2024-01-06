package op

import "dancavallaro.com/synacor-go/pkg/memory"

const (
	testWordAddr = 42
	testByteAddr = 2 * testWordAddr
)

const (
	R0 = memory.RegOffset
	R1 = R0 + 1
)

func args(a ...uint16) []uint16 {
	return a
}
