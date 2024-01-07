package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackUnderflow(t *testing.T) {
	m := Memory{}
	assert.Panics(t, func() {
		m.Pop()
	})
}

func TestStack(t *testing.T) {
	m := Memory{}
	m.Push(1)
	m.Push(2)
	m.Push(3)
	assert.Equal(t, uint16(3), m.Pop())
	assert.Equal(t, uint16(2), m.Pop())
	assert.Equal(t, uint16(1), m.Pop())
	assert.Equal(t, 0, len(m.Stack))
}
