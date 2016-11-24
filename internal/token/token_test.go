package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	str, err := GenerateRandomString(12)
	assert.NotEmpty(t, str)
	assert.Nil(t, err)
}
