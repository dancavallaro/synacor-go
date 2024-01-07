package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJump(t *testing.T) {
	m := &memory.Memory{}
	Jmp(m, args(testAddr))
	assert.Equal(t, testAddr, m.PC)
}

func TestJump_RegTarget(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = testAddr
	Jmp(m, args(R0))
	assert.Equal(t, testAddr, m.PC)
}

func TestJt_NoJump(t *testing.T) {
	m := &memory.Memory{}
	Jt(m, args(0, testAddr))
	assert.Equal(t, 0, m.PC)
}

func TestJt_Jump(t *testing.T) {
	m := &memory.Memory{}
	Jt(m, args(1, testAddr))
	assert.Equal(t, testAddr, m.PC)
}

func TestJt_Jump_RegTarget(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = testAddr
	Jt(m, args(1, R0))
	assert.Equal(t, testAddr, m.PC)
}

func TestJt_NoJump_RegArg(t *testing.T) {
	m := &memory.Memory{}
	Jt(m, args(R0, testAddr))
	assert.Equal(t, 0, m.PC)
}

func TestJf_NoJump(t *testing.T) {
	m := &memory.Memory{}
	Jf(m, args(1, testAddr))
	assert.Equal(t, 0, m.PC)
}

func TestJf_Jump(t *testing.T) {
	m := &memory.Memory{}
	Jf(m, args(0, testAddr))
	assert.Equal(t, testAddr, m.PC)
}

func TestJf_Jump_RegArg(t *testing.T) {
	m := &memory.Memory{}
	Jf(m, args(R0, testAddr))
	assert.Equal(t, testAddr, m.PC)
}

func TestJf_Jump_RegTarget(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = testAddr
	Jf(m, args(0, R0))
	assert.Equal(t, testAddr, m.PC)
}
