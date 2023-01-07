package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type renewAccessResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
type renewAccessRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var rq renewAccessRequest

	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshTokenPayload, err := server.tokenMaker.VerifyToken(rq.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshTokenPayload.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if session.Username != refreshTokenPayload.Username {
		err := fmt.Errorf("incorrect user session")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if session.RefreshToken != rq.RefreshToken {
		err := fmt.Errorf("mismatch session refresh token")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.
		CreateToken(
			refreshTokenPayload.Username,
			server.config.RefreshTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rsp := renewAccessResponse{
		AccessToken:          refreshToken,
		AccessTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)

}
