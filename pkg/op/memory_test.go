package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	m := &memory.Memory{}
	Set(m, args(R0, 42))
	assert.Equal(t, uint16(42), m.GP[0])
}

func TestSet_RegVal(t *testing.T) {
	m := &memory.Memory{}
	m.GP[1] = 42
	Set(m, args(R0, R1))
	assert.Equal(t, uint16(42), m.GP[0])
}

func TestPop(t *testing.T) {
	m := &memory.Memory{}
	m.Stack = append(m.Stack, 42)
	Pop(m, args(R0))
	assert.Equal(t, uint16(42), m.GP[0])
}

func TestPush(t *testing.T) {
	m := &memory.Memory{}
	Push(m, args(42))
	assert.Equal(t, uint16(42), m.Stack[0])
}

func TestPush_RegArg(t *testing.T) {
	m := &memory.Memory{}
	m.GP[0] = 42
	Push(m, args(R0))
	assert.Equal(t, uint16(42), m.Stack[0])
}
