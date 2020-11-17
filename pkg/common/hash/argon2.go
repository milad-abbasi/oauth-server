package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

var DefaultArgonConfig = ArgonConfig{
	time:    1,
	memory:  64 * 1024,
	threads: 4,
	keyLen:  32,
}

func GenerateArgon2Hash(input string) (string, error) {
	return GenerateArgon2HashWithConfig(input, &DefaultArgonConfig)
}

func GenerateArgon2HashWithConfig(input string, c *ArgonConfig) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(input), salt, c.time, c.memory, c.threads, c.keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)

	return full, nil
}

func CompareArgon2Hash(input, hash string) (bool, error) {
	parts := strings.Split(hash, "$")

	c := &ArgonConfig{}
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

	comparisonHash := argon2.IDKey([]byte(input), salt, c.time, c.memory, c.threads, c.keyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
