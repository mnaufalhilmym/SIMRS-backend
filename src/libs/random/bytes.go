package random

import (
	"crypto/rand"
	"errors"
)

func Bytes(n int) ([]byte, error) {
	if n <= 0 {
		err := errors.New("invalid length of random bytes")
		logger.Error(err)
		return nil, err
	}

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, nil
}
