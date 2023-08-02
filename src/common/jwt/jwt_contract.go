package jwt

import (
	accountrole "simrs/src/common/account_role"

	"github.com/google/uuid"
)

type JwtPayload struct {
	ID         *uuid.UUID        `json:"id"`
	Role       *accountrole.Role `json:"role"`
	Expiration *int64            `json:"exp"`
}

type JwtCtxKey string
