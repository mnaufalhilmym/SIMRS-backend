package argon2

import (
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

func CompareStringAndEncodedHash(str, encodedHash *string) (bool, error) {
	argon2Params, salt, hash, err := GetDecodedHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey(
		[]byte(*str),
		*salt,
		argon2Params.Iterations,
		argon2Params.Memory,
		argon2Params.Parallelism,
		argon2Params.KeyLength,
	)

	if subtle.ConstantTimeCompare(*hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
