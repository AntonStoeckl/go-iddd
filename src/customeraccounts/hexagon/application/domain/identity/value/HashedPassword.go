package value

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"golang.org/x/crypto/argon2"
)

type HashedPassword string

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func HashedPasswordFromPlainPassword(input PlainPassword) (HashedPassword, error) {
	config := &PasswordConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}

	pw, err := generatePassword(config, input.String())
	if err != nil {
		return "", shared.MarkAndWrapError(err, shared.ErrTechnical, "HashedPasswordFromPlainPassword")
	}

	return HashedPassword(pw), nil
}

func RebuildHashedPassword(input string) HashedPassword {
	return HashedPassword(input)
}

func (pw HashedPassword) CompareWith(plainPassword PlainPassword) bool {
	match, err := pw.comparePassword(plainPassword.String(), pw.String())
	if !match || err != nil {
		return false
	}

	return true
}

func (pw HashedPassword) Equals(other HashedPassword) bool {
	return pw.String() == other.String()
}

func (pw HashedPassword) String() string {
	return string(pw)
}

func generatePassword(config *PasswordConfig, password string) (string, error) {
	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, config.memory, config.time, config.threads, b64Salt, b64Hash)

	return full, nil
}

func (pw HashedPassword) comparePassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")

	c := &PasswordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	c.keyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
