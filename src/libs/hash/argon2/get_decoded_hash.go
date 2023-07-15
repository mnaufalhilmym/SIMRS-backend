package argon2

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func GetDecodedHash(encodedHash *string) (*Config, *[]byte, *[]byte, error) {
	vals := strings.Split(*encodedHash, "$")
	if len(vals) != 4 {
		err := errors.New("invalid encoded hash")
		logger.Error(err)
		return nil, nil, nil, err
	}

	{
		version := new(int)
		_, err := fmt.Sscanf(vals[0], "v=%d", version)
		if err != nil {
			logger.Error(err)
			return nil, nil, nil, err
		}
		if *version != argon2.Version {
			err := errors.New("incompatible argon2 version")
			logger.Error(err)
			return nil, nil, nil, err
		}
	}

	argon2Config := new(Config)
	_, err := fmt.Sscanf(vals[1], "m=%d,t=%d,p=%d", &argon2Config.Memory, &argon2Config.Iterations, &argon2Config.Parallelism)
	if err != nil {
		logger.Error(err)
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[2])
	if err != nil {
		logger.Error(err)
		return nil, nil, nil, err
	}
	argon2Config.SaltLength = len(salt)

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[3])
	if err != nil {
		logger.Error(err)
		return nil, nil, nil, err
	}
	argon2Config.KeyLength = uint32(len(hash))

	return argon2Config, &salt, &hash, nil
}
