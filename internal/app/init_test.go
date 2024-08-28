package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer(t *testing.T) {
	f, _ := os.Create("config.yml")
	defer func() {
		_ = os.RemoveAll(f.Name())
	}()
	_ = os.Setenv("CONFIG_PATH", f.Name())

	_, err := NewContainer()
	assert.Nil(t, err)

	_ = os.Setenv("CONFIG_PATH", "")
	_, err = NewContainer()
	assert.Error(t, err)
}
