package constant

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	Authorization          = "Authorization"
	Bearer                 = "Bearer"
	BearerIndex            = 0
	TokenIndex             = 1
	ServerWriteReadTimeout = 15
	ServerIdleTimeout      = 60
	ServerReleaseShotdown  = 5
	ServerContextTimeout   = 1 * time.Minute
)

type TransactionName string

type JWTClaim struct {
	CPF string `json:"cpf"`
	jwt.RegisteredClaims
}
