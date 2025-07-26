package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/pdh9523/go-hexarch/shared/security/error_code"
	"golang.org/x/crypto/argon2"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}

type Argon2PasswordHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewArgon2PasswordHasher() *Argon2PasswordHasher {
	return &Argon2PasswordHasher{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

func (h *Argon2PasswordHasher) HashPassword(password string) (string, error) {
	salt, err := generateRandomBytes(h.saltLength)
	if err != nil {
		return "", errors.New("failed to generate salt: " + err.Error())
	}

	hash := argon2.IDKey([]byte(password), salt, h.iterations, h.memory, h.parallelism, h.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, h.memory, h.iterations, h.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

func (h *Argon2PasswordHasher) VerifyPassword(hashedPassword, password string) error {
	salt, hash, err := h.decodeHash(hashedPassword)
	if err != nil {
		return err
	}

	otherHash := argon2.IDKey([]byte(password), salt, h.iterations, h.memory, h.parallelism, h.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}
	return error_code.ErrPasswordMismatch
}

func (h *Argon2PasswordHasher) decodeHash(encodedHash string) (salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, error_code.ErrInvalidHashFormat
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", error_code.ErrHashDecodingFailed, err)
	}
	if version != argon2.Version {
		return nil, nil, error_code.ErrIncompatibleHashVersion
	}

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &h.memory, &h.iterations, &h.parallelism)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", error_code.ErrHashDecodingFailed, err)
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", error_code.ErrHashDecodingFailed, err)
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", error_code.ErrHashDecodingFailed, err)
	}

	return salt, hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
