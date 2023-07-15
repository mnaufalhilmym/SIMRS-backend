package jwt

import (
	"errors"

	"github.com/bytedance/sonic"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Parse[T any](payload string, target *T) error {
	if conf.privateKey == nil {
		logger.Panic(errors.New("cannot parse jwt because of empty private key"))
	}

	token, err := jwt.Parse(
		[]byte(payload),
		jwt.WithKey(jwa.RS256, &conf.privateKey.PublicKey),
		jwt.WithValidate(true),
		jwt.WithVerify(true),
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	buf, err := sonic.Marshal(token)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err := sonic.Unmarshal(buf, target); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
