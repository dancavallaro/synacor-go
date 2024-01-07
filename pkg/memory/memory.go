package memory

import "log"

const (
	Modulus      = 32768
	RegOffset    = Modulus
	MaxInt       = 32767
	NumRegisters = 8
)

type Memory struct {
	PC    int
	GP    [NumRegisters]uint16
	Stack []uint16
	Mem   [1 << 16]uint16
}

func (m *Memory) Push(val uint16) {
	m.Stack = append(m.Stack, val)
}

func (m *Memory) Pop() uint16 {
	if len(m.Stack) == 0 {
		log.Panicln("stack underflow!")
	}
	val := m.Stack[len(m.Stack)-1]
	m.Stack = m.Stack[0 : len(m.Stack)-1]
	return val
}

func (m *Memory) ReadWord(addr uint16) uint16 {
	byteAddr := addr * 2
	return (m.Mem[byteAddr+1] << 8) + m.Mem[byteAddr]
}

func (m *Memory) WriteWord(addr uint16, word uint16) {
	byteAddr := addr * 2
	m.Mem[byteAddr] = word & 0xFF
	m.Mem[byteAddr+1] = (word & 0xFF00) >> 8
}

func ReadVal(m *Memory, arg uint16) uint16 {
	if arg >= RegOffset+NumRegisters {
		log.Panicf("arg %d is not a valid value (literal or register)", arg)
	} else if arg >= RegOffset {
		return m.GP[arg-RegOffset]
	}
	return arg
}

func RegNum(arg uint16) uint16 {
	if arg < RegOffset || arg >= RegOffset+NumRegisters {
		log.Panicf("arg %d is not a valid register value", arg)
	}
	return arg - RegOffset
}
