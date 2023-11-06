package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount(t *testing.T) {
	_, err := NewAccount("test", "test", "test")

	assert.Nil(t, err)
}
