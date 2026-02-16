package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Logstay/tech-challange-execution-service/config"
	"github.com/Logstay/tech-challange-execution-service/pkg/constant"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthzWithConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAuth := c.GetHeader(constant.Authorization)
		if headerAuth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: "não autorizado",
			})
			c.Abort()
			return
		}

		err := ValidateTokenWithConfig(headerAuth, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.FromContext(c.Request.Context())
		headerAuth := c.GetHeader(constant.Authorization)
		if headerAuth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: "não autorizado",
			})
			c.Abort()
			return
		}

		err := ValidateTokenWithConfig(headerAuth, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ValidateTokenWithConfig(signedToken string, cfg *config.Config) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&constant.JWTClaim{},
		func(token *jwt.Token) (any, error) {
			return []byte(cfg.JWT.Secret), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return errors.New("token expirado")
		}
		return fmt.Errorf("token invalido: %w", err)
	}

	_, ok := token.Claims.(*constant.JWTClaim)
	if !ok || !token.Valid {
		return errors.New("token invalido")
	}

	return nil
}

func ValidateToken(signedToken string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	return ValidateTokenWithConfig(signedToken, cfg)
}
