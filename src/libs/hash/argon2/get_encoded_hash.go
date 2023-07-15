package argon2

import (
	"encoding/base64"
	"fmt"
	"simrs/src/libs/random"

	"golang.org/x/crypto/argon2"
)

func GetEncodedHash(str *string) (*string, error) {
	salt, hash, err := getHash(str)
	if err != nil {
		return nil, err
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(*salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(*hash)

	encodedHash := fmt.Sprintf(
		"v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argon2Config.Memory,
		argon2Config.Iterations,
		argon2Config.Parallelism,
		b64Salt,
		b64Hash,
	)

	return &encodedHash, nil
}

func getHash(str *string) (*[]byte, *[]byte, error) {
	salt, err := random.Bytes(argon2Config.SaltLength)
	if err != nil {
		return nil, nil, err
	}

	hash := argon2.IDKey(
		[]byte(*str),
		salt,
		argon2Config.Iterations,
		argon2Config.Memory,
		argon2Config.Parallelism,
		argon2Config.KeyLength,
	)

	return &salt, &hash, nil
}
