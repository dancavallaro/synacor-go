package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testWordAddr = 42
	testByteAddr = 2 * testWordAddr
)

func TestJump(t *testing.T) {
	r := &memory.Registers{}
	Jmp(r, args(testWordAddr))
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

func args(a ...uint16) []uint16 {
	return a
}
