package jwt

import (
	"errors"
	"time"

	"github.com/bytedance/sonic"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Create[T any](data *T) (*string, error) {
	if conf.privateKey == nil {
		logger.Panic(errors.New("cannot create jwt because of empty private key"))
	}

	token, err := jwt.NewBuilder().Expiration(time.Now().Add(*conf.duration)).Build()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	jwtDataMap := make(map[string]interface{}, 0)
	jwtDataBytes, err := sonic.Marshal(*data)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	sonic.Unmarshal(jwtDataBytes, &jwtDataMap)

	for key, value := range jwtDataMap {
		if value != nil {
			if err := token.Set(key, value); err != nil {
				logger.Error(err)
				return nil, err
			}
		}
	}

	payload, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, conf.privateKey))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	strPayload := string(payload)
	return &strPayload, nil
}
