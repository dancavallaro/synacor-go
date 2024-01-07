package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	m := &memory.Memory{}
	Add(m, args(R0, 1, 2))
	assert.Equal(t, uint16(3), m.GP[0])
}

func TestAdd_RegArgs(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = 1
	m.GP[1] = 2
	Add(m, args(R2, R0, R1))
	assert.Equal(t, uint16(3), m.GP[2])
}

func TestEq_False(t *testing.T) {
	m := &memory.Memory{}
	Eq(m, args(R0, 1, 2))
	assert.Equal(t, uint16(0), m.GP[0])
}

func TestEq_True(t *testing.T) {
	m := &memory.Memory{}
	Eq(m, args(R0, 1, 1))
	assert.Equal(t, uint16(1), m.GP[0])
}

func TestEq_False_RegArgs(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = 1
	m.GP[1] = 2
	Eq(m, args(R2, R0, R1))
	assert.Equal(t, uint16(0), m.GP[2])
}

func TestEq_True_RegArgs(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = 1
	m.GP[1] = 1
	Eq(m, args(R2, R0, R1))
	assert.Equal(t, uint16(1), m.GP[2])
}
