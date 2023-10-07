package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfiguration(t *testing.T) {
	configuration, err := ReadConfiguration("test_config.json")

	assert.NoError(t, err)

	assert.True(t, configuration.IsValidConfiguration(), "Configuration read should be valid")
}
