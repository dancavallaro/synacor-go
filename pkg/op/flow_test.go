package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJump(t *testing.T) {
	r := &memory.Registers{}
	Jmp(r, args(testWordAddr))
	assert.Equal(t, testByteAddr, r.PC)
}

func TestJump_RegTarget(t *testing.T) {
	r := &memory.Registers{}
	r.GP[0] = testWordAddr
	Jmp(r, args(R0))
	assert.Equal(t, testByteAddr, r.PC)
}

func TestJt_NoJump(t *testing.T) {
	r := &memory.Registers{}
	Jt(r, args(0, testWordAddr))
	assert.Equal(t, 0, r.PC)
}

func TestJt_Jump(t *testing.T) {
	r := &memory.Registers{}
	Jt(r, args(1, testWordAddr))
	assert.Equal(t, testByteAddr, r.PC)
}

func TestJt_Jump_RegTarget(t *testing.T) {
	r := &memory.Registers{}
	r.GP[0] = testWordAddr
	Jt(r, args(1, R0))
	assert.Equal(t, testByteAddr, r.PC)
}

func TestJt_NoJump_RegArg(t *testing.T) {
	r := &memory.Registers{}
	Jt(r, args(R0, testWordAddr))
	assert.Equal(t, 0, r.PC)
}

func TestJf_NoJump(t *testing.T) {
	r := &memory.Registers{}
	Jf(r, args(1, testWordAddr))
	assert.Equal(t, 0, r.PC)
}

func TestJf_Jump(t *testing.T) {
	r := &memory.Registers{}
	Jf(r, args(0, testWordAddr))
	assert.Equal(t, testByteAddr, r.PC)
}

func TestJf_Jump_RegArg(t *testing.T) {
	r := &memory.Registers{}
	Jf(r, args(R0, testWordAddr))
	assert.Equal(t, testByteAddr, r.PC)
}

func TestJf_Jump_RegTarget(t *testing.T) {
	r := &memory.Registers{}
	r.GP[0] = testWordAddr
	Jf(r, args(0, R0))
	assert.Equal(t, testByteAddr, r.PC)
}
