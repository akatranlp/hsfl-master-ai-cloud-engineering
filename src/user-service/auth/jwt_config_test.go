package auth

import (
	"errors"
	"os"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/utils"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestJwtConfig(t *testing.T) {
	privateKeyData, publicKeyData := utils.GenerateRSAKeyPairPem()
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyData))
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyData))

	ctrl := gomock.NewController(t)
	config := mocks.NewMockConfig(ctrl)

	t.Run("CreateTokenGenerator", func(t *testing.T) {
		t.Run("should error if privateKey cannot be read", func(t *testing.T) {
			// given
			config.
				EXPECT().
				ReadPrivateKey().
				Return(nil, errors.New("read error"))

			// when
			gen, err := NewJwtTokenGenerator(config)

			// then
			assert.Error(t, err)
			assert.Nil(t, gen)
		})

		t.Run("should error if publicKey cannot be read", func(t *testing.T) {
			// given
			config.
				EXPECT().
				ReadPrivateKey().
				Return("", nil)

			config.
				EXPECT().
				ReadPublicKey().
				Return(nil, errors.New("read error"))

			// when
			gen, err := NewJwtTokenGenerator(config)

			// then
			assert.Error(t, err)
			assert.Nil(t, gen)
		})

		t.Run("should error if privateKey cannot be casted", func(t *testing.T) {
			// given
			config.
				EXPECT().
				ReadPrivateKey().
				Return("", nil)

			config.
				EXPECT().
				ReadPublicKey().
				Return("", nil)

			// when
			gen, err := NewJwtTokenGenerator(config)

			// then
			assert.Error(t, err)
			assert.Nil(t, gen)
		})

		t.Run("should error if publicKey cannot be casted", func(t *testing.T) {
			// given
			config.
				EXPECT().
				ReadPrivateKey().
				Return(privateKey, nil)

			config.
				EXPECT().
				ReadPublicKey().
				Return("", nil)

			// when
			gen, err := NewJwtTokenGenerator(config)

			// then
			assert.Error(t, err)
			assert.Nil(t, gen)
		})

		t.Run("should error if publicKey cannot be casted", func(t *testing.T) {
			// given
			config.
				EXPECT().
				ReadPrivateKey().
				Return(privateKey, nil)

			config.
				EXPECT().
				ReadPublicKey().
				Return(publicKey, nil)

			// when
			gen, err := NewJwtTokenGenerator(config)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, gen)
		})
	})

	t.Run("ReadPrivateKey direct", func(t *testing.T) {
		t.Run("should error", func(t *testing.T) {
			// given
			config := JwtConfig{
				PrivateKey: "t",
				PublicKey:  "t",
			}

			// when
			_, err := config.ReadPrivateKey()

			// then
			assert.Error(t, err)
		})

		t.Run("should not error", func(t *testing.T) {
			// given
			config := JwtConfig{
				PrivateKey: privateKeyData,
				PublicKey:  "t",
			}

			// when
			key, err := config.ReadPrivateKey()

			// then
			assert.NoError(t, err)
			assert.Equal(t, privateKey, key)
		})
	})

	t.Run("ReadPrivateKey from file", func(t *testing.T) {
		t.Run("should error because file doesn't exist", func(t *testing.T) {
			// given
			config := JwtConfig{
				PrivateKeyPath: "t",
				PublicKeyPath:  "t",
			}

			// when
			_, err := config.ReadPrivateKey()

			// then
			assert.Error(t, err)
		})

		t.Run("should error", func(t *testing.T) {
			// given
			f, _ := os.Create("private")
			f.WriteString("t")
			f.Close()

			config := JwtConfig{
				PrivateKeyPath: "private",
				PublicKeyPath:  "t",
			}

			// when
			_, err := config.ReadPrivateKey()

			// then
			os.Remove("private")
			assert.Error(t, err)
		})

		t.Run("should not error", func(t *testing.T) {
			// given
			f, _ := os.Create("private")
			f.WriteString(privateKeyData)
			f.Close()

			config := JwtConfig{
				PrivateKeyPath: "private",
				PublicKeyPath:  "t",
			}

			// when
			key, err := config.ReadPrivateKey()

			// then
			os.Remove("private")
			assert.NoError(t, err)
			assert.Equal(t, privateKey, key)
		})
	})

	t.Run("ReadPublicKey direct", func(t *testing.T) {
		t.Run("should error", func(t *testing.T) {
			// given
			config := JwtConfig{
				PrivateKey: "t",
				PublicKey:  "t",
			}

			// when
			_, err := config.ReadPublicKey()

			// then
			assert.Error(t, err)
		})

		t.Run("should not error", func(t *testing.T) {
			// given

			config := JwtConfig{
				PrivateKey: "t",
				PublicKey:  publicKeyData,
			}

			// when
			key, err := config.ReadPublicKey()

			// then
			assert.NoError(t, err)
			assert.Equal(t, publicKey, key)
		})
	})

	t.Run("ReadPublicKey from file", func(t *testing.T) {
		t.Run("should error because file doesn't exist", func(t *testing.T) {
			// given
			config := JwtConfig{
				PrivateKeyPath: "t",
				PublicKeyPath:  "t",
			}

			// when
			_, err := config.ReadPublicKey()

			// then
			assert.Error(t, err)
		})

		t.Run("should error", func(t *testing.T) {
			// given
			f, _ := os.Create("public")
			f.WriteString("t")
			f.Close()

			config := JwtConfig{
				PrivateKeyPath: "t",
				PublicKeyPath:  "public",
			}

			// when
			_, err := config.ReadPublicKey()

			// then
			os.Remove("public")
			assert.Error(t, err)
		})

		t.Run("should not error", func(t *testing.T) {
			// given
			f, _ := os.Create("public")
			f.WriteString(publicKeyData)
			f.Close()

			config := JwtConfig{
				PrivateKeyPath: "t",
				PublicKeyPath:  "public",
			}

			// when
			key, err := config.ReadPublicKey()

			// then
			os.Remove("public")
			assert.NoError(t, err)
			assert.Equal(t, publicKey, key)
		})
	})
}
