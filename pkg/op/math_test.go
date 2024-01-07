package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	r := &memory.Registers{}
	Add(r, args(R0, 1, 2))
	assert.Equal(t, uint16(3), r.GP[0])
}

func TestAdd_RegArgs(t *testing.T) {
	r := &memory.Registers{}
	r.GP[0] = 1
	r.GP[1] = 2
	Add(r, args(R2, R0, R1))
	assert.Equal(t, uint16(3), r.GP[2])
}

func TestEq_False(t *testing.T) {
	r := &memory.Registers{}
	Eq(r, args(R0, 1, 2))
	assert.Equal(t, uint16(0), r.GP[0])
}

func TestEq_True(t *testing.T) {
	r := &memory.Registers{}
	Eq(r, args(R0, 1, 1))
	assert.Equal(t, uint16(1), r.GP[0])
}

func TestEq_False_RegArgs(t *testing.T) {
	r := &memory.Registers{}
	r.GP[0] = 1
	r.GP[1] = 2
	Eq(r, args(R2, R0, R1))
	assert.Equal(t, uint16(0), r.GP[2])
}

func TestEq_True_RegArgs(t *testing.T) {
	r := &memory.Registers{}
	r.GP[0] = 1
	r.GP[1] = 1
	Eq(r, args(R2, R0, R1))
	assert.Equal(t, uint16(1), r.GP[2])
}
