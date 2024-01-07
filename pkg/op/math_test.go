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
