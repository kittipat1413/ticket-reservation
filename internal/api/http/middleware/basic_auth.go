package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"ticket-reservation/internal/util/httpresponse"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"

	"github.com/gin-gonic/gin"
)

func (m *middleware) BasicAuth(username, password string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			// response with default unauthorized error
			httpresponse.Error(ctx, errsFramework.NewUnauthorizedError("", nil))
			return
		}
		key, secret, ok := ctx.Request.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(key))
			passwordHash := sha256.Sum256([]byte(secret))

			expectedUsernameHash := sha256.Sum256([]byte(username))
			expectedPasswordHash := sha256.Sum256([]byte(password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				ctx.Next()
			} else {
				// response with default unauthorized error
				httpresponse.Error(ctx, errsFramework.NewUnauthorizedError("", nil))
			}
		} else {
			// response with default unauthorized error
			httpresponse.Error(ctx, errsFramework.NewUnauthorizedError("", nil))
		}

	}
}
