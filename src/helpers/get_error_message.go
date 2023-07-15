package helpers

import (
	"simrs/src/constants"
	"simrs/src/libs/env"
)

func GetErrorMessage(errMsg ...string) *string {
	if len(errMsg) > 1 && env.Get(env.APP_MODE) == constants.APP_MODE_RELEASE {
		return &errMsg[1]
	}
	return &errMsg[0]
}
