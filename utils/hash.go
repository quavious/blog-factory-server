package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func Hash(text string) (string, error) {
	param := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	salt := make([]byte, param.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(text), salt, param.iterations, param.memory, param.parallelism, param.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, param.memory, param.iterations, param.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

func Verify(plainText, encodedHash string) (bool, error) {
	split := strings.Split(encodedHash, "$")
	if len(split) != 6 {
		return false, errors.New("invalid hash")
	}
	version := 0
	_, err := fmt.Sscanf(split[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, errors.New("incompatible version")
	}
	param := &params{}
	_, err = fmt.Sscanf(split[3], "m=%d,t=%d,p=%d", &param.memory, &param.iterations, &param.parallelism)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.Strict().DecodeString(split[4])
	if err != nil {
		return false, err
	}
	param.saltLength = uint32(len(salt))
	hash, err := base64.RawStdEncoding.Strict().DecodeString(split[5])
	if err != nil {
		return false, err
	}
	param.keyLength = uint32(len(hash))
	otherHash := argon2.IDKey([]byte(plainText), salt, param.iterations, param.memory, param.parallelism, param.keyLength)
	if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
		return false, nil
	}
	return true, nil
}
