package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJump(t *testing.T) {
	m := &memory.Memory{}
	Jmp(m, args(testWordAddr))
	assert.Equal(t, testByteAddr, m.PC)
}

func TestJump_RegTarget(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = testWordAddr
	Jmp(m, args(R0))
	assert.Equal(t, testByteAddr, m.PC)
}

func TestJt_NoJump(t *testing.T) {
	m := &memory.Memory{}
	Jt(m, args(0, testWordAddr))
	assert.Equal(t, 0, m.PC)
}

func TestJt_Jump(t *testing.T) {
	m := &memory.Memory{}
	Jt(m, args(1, testWordAddr))
	assert.Equal(t, testByteAddr, m.PC)
}

func TestJt_Jump_RegTarget(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = testWordAddr
	Jt(m, args(1, R0))
	assert.Equal(t, testByteAddr, m.PC)
}

func TestJt_NoJump_RegArg(t *testing.T) {
	m := &memory.Memory{}
	Jt(m, args(R0, testWordAddr))
	assert.Equal(t, 0, m.PC)
}

func TestJf_NoJump(t *testing.T) {
	m := &memory.Memory{}
	Jf(m, args(1, testWordAddr))
	assert.Equal(t, 0, m.PC)
}

func TestJf_Jump(t *testing.T) {
	m := &memory.Memory{}
	Jf(m, args(0, testWordAddr))
	assert.Equal(t, testByteAddr, m.PC)
}

func TestJf_Jump_RegArg(t *testing.T) {
	m := &memory.Memory{}
	Jf(m, args(R0, testWordAddr))
	assert.Equal(t, testByteAddr, m.PC)
}

func TestJf_Jump_RegTarget(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = testWordAddr
	Jf(m, args(0, R0))
	assert.Equal(t, testByteAddr, m.PC)
}
