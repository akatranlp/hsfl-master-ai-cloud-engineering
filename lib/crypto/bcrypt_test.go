package crypto

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHasher(t *testing.T) {
	hasher := NewBcryptHasher()

	t.Run("Hash", func(t *testing.T) {
		t.Run("should return hash with salt", func(t *testing.T) {
			// given
			password := []byte("password")

			// when
			hash, err := hasher.Hash(password)

			// then
			assert.NoError(t, err)
			assert.Len(t, hash, 60)
			assert.Regexp(t, regexp.MustCompile(`\$2a\$10\$(.*)`), string(hash))
		})
	})

	t.Run("Validate", func(t *testing.T) {
		t.Run("should return true if password matches hash", func(t *testing.T) {
			// given
			password := []byte("password")
			hash := []byte("$2a$10$s3BvNfI4PZO0PhcyxK4vTeu0N3Hhxo4mMgd084ENY41q/DeXhstc6")

			// when
			ok := hasher.Validate(password, hash)

			// then
			assert.True(t, ok)
		})

		t.Run("should return false if password does not match hash", func(t *testing.T) {
			// given
			password := []byte("password")
			hash := []byte("$2a$10$s3BvNfI4PZO0PhcyxK4vTeu0N3Hhxo4mMgd084ENY41q/DeXhstc7")

			// when
			ok := hasher.Validate(password, hash)

			// then
			assert.False(t, ok)
		})
	})
}
