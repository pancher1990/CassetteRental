package hash

import (
	"crypto/sha1"
	"errors"
	"fmt"
)

type SHA1Hasher struct {
	Salt string
}

func NewHash(salt string) (*SHA1Hasher, error) {
	if salt == "" {
		return nil, errors.New("empty salt for hash")
	}
	return &SHA1Hasher{salt}, nil
}

func (h *SHA1Hasher) Hash(password string) (string, error) {
	const op = "internal/lib/hash/password/Hash"
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.Salt))), nil
}
