package unit

import (
	"testing"

	"github.com/nabilfikrisp/url-shortener/internal/common/helpers"
	"github.com/stretchr/testify/assert"
)

func TestOurDomainValidator(t *testing.T) {
	t.Run("empty URL", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("example.com", "")
		assert.Error(t, err)
		assert.False(t, got)
	})

	t.Run("invalid URL", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("example.com", "://bad-url")
		assert.Error(t, err)
		assert.False(t, got)
	})

	t.Run("URL without hostname", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("example.com", "http://")
		assert.Error(t, err)
		assert.False(t, got)
	})

	t.Run("matching domain", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("example.com", "https://example.com/page")
		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("matching domain with port in ourDomain", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("example.com:8080", "https://example.com/abc")
		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("non-matching domain", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("example.com", "https://other.com/page")
		assert.NoError(t, err)
		assert.False(t, got)
	})

	t.Run("localhost is always allowed", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("example.com", "http://localhost:3000/page")
		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("case-insensitive match", func(t *testing.T) {
		got, err := helpers.OurDomainValidator("Example.COM", "https://example.com/")
		assert.NoError(t, err)
		assert.True(t, got)
	})
}
