package unit

import (
	"testing"

	"github.com/nabilfikrisp/url-shortener/internal/common/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGenerateShortToken(t *testing.T) {
	t.Run("deterministic - same input gives same output", func(t *testing.T) {
		input := "https://test.com"

		result1 := helpers.GenerateShortToken(input)
		result2 := helpers.GenerateShortToken(input)

		assert.Equal(t, result1, result2)
	})

	t.Run("different inputs give different outputs", func(t *testing.T) {
		input1 := "https://test1.com"
		input2 := "https://test2.com"

		result1 := helpers.GenerateShortToken(input1)
		result2 := helpers.GenerateShortToken(input2)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("always returns 16 character hex string", func(t *testing.T) {
		inputs := []string{
			"a",
			"hello world",
			"https://very-long-domain-name.com/with/long/path?param=value",
			"",
			"123456789",
		}

		for _, input := range inputs {
			result := helpers.GenerateShortToken(input)

			assert.Len(t, result, 16)
			assert.Regexp(t, "^[0-9a-f]+$", result)
		}
	})
}
