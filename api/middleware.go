package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/techschool/simple-bank/token"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorisationPayloadKey = "authorization_payload"
)

func authMiddleware(maker token.Maker) gin.HandlerFunc {

	return func(context *gin.Context) {
		authorizationHeader := context.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorisation header format")
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorisationType := strings.ToLower(fields[0])

		if authorisationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s ", authorisationType)
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		context.Set(authorisationPayloadKey, payload)
		context.Next()
	}
}
