package context

import (
	"backend/constant"
	"backend/logger"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type passwordContextIE struct {
	*logger.BackendLogger
}

type passwordContext struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32

	*logger.BackendLogger
}

func newPasswordContext(passwordContextIE *passwordContextIE) (*passwordContext, error) {
	return &passwordContext{
		memory:      constant.PWD_MEMORY,
		iterations:  constant.PWD_ITERATIONS,
		parallelism: constant.PWD_PARALLELISM,
		saltLength:  constant.PWD_SALT_LENGTH,
		keyLength:   constant.PWD_KEY_LENGTH,

		BackendLogger: passwordContextIE.BackendLogger,
	}, nil
}

func (p *passwordContext) release() {
	p.PwdLog.Infoln("Release passwordContext...")

	p.ProcLog.Infoln("passwordContext released")
}

func (p *passwordContext) Hash(password string) (string, error) {
	salt := make([]byte, p.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.memory,
		p.iterations,
		p.parallelism,
		encodedSalt,
		encodedHash,
	)
	return encoded, nil

}

func (p *passwordContext) Verify(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid password hash format")
	}
	var version int
	var memory uint32
	var iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, errors.New("invalid argon2 version")
	}
	if version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}
	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&parallelism,
	); err != nil {
		return false, errors.New("invalid argon2 parameters")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("invalid salt")
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("invalid hash")
	}
	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(expectedHash)),
	)
	match := subtle.ConstantTimeCompare(actualHash, expectedHash) == 1
	return match, nil

}
