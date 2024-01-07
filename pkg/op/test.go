package op

import "dancavallaro.com/synacor-go/pkg/memory"

const (
	testAddr = 42
)

const (
	R0 = memory.RegOffset
	R1 = R0 + 1
	R2 = R0 + 2
)

func args(a ...uint16) []uint16 {
	return a
}
