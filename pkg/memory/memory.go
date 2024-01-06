package memory

import "log"

const (
	Modulus      = 32768
	RegOffset    = Modulus
	MaxInt       = 32767
	NumRegisters = 8
)

type Registers struct {
	PC int
	GP [NumRegisters]uint16
}

func ReadVal(r *Registers, arg uint16) uint16 {
	if arg >= RegOffset+NumRegisters {
		log.Panicf("arg %d is not a valid value (literal or register)", arg)
	} else if arg >= RegOffset {
		return r.GP[arg-RegOffset]
	}
	return arg
}
