package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/util"
	"net/http"
	"time"
)

type userResponse struct {
	Username          string       `json:"username"`
	FullName          string       `json:"fullName"`
	Email             string       `json:"email"`
	CreatedAt         sql.NullTime `json:"createdAt"`
	PasswordChangedAt time.Time    `json:"passwordChangedAt"`
}
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password"  binding:"required,min=6"`
	FullName string `json:"full_name"  binding:"required,min=6"`
	Email    string `json:"email"  binding:"required,email"`
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password"  binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionId            uuid.UUID `json:"session_id"`
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`

	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`

	User userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var rq loginUserRequest

	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, rq.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	err = util.CheckPassword(rq.Password, user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessTokenPaylod, err := server.tokenMaker.CreateToken(
		user.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username, server.config.RefreshTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.CreateSessionParams{
		ID:           refreshTokenPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}
	session, err := server.store.CreateSession(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPaylod.ExpiredAt,
		SessionId:             session.ID,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		User:                  NewUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)

}
func (server *Server) createUser(ctx *gin.Context) {
	var rq createUserRequest
	if err := ctx.ShouldBindJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(rq.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       rq.Username,
		FullName:       rq.FullName,
		Email:          rq.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	resp := NewUserResponse(user)

	ctx.JSON(http.StatusOK, resp)

}

func NewUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}
