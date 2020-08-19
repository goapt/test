package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	assert.NotPanics(t, func() {
		redis := NewRedis()
		err := redis.Set("test", "123")
		assert.NoError(t, err)
		s, err := redis.Get("test")
		assert.NoError(t, err)
		assert.Equal(t, "123", s)
	})
}
