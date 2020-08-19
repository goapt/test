package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.NotPanics(t, func() {
		Config()
	})
}
