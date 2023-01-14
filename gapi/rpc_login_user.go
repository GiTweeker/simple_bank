package gapi

import (
	"context"
	"database/sql"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/pb"
	"github.com/techschool/simple-bank/util"
	"github.com/techschool/simple-bank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, rq *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	if violations := validateLoginUserRequest(rq); violations != nil {
		return nil, invalidArgument(violations)
	}

	user, err := server.store.GetUser(ctx, rq.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound,
				"username does not exist : %s", rq.GetUsername())
		}

		return nil, status.Errorf(codes.Internal,
			"failed to create user, error : %s", err)
	}

	err = util.CheckPassword(rq.GetPassword(), user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied,
			"invalid password or username")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username, server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"unable to create token")
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username, server.config.RefreshTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"unable to create token")
	}
	metadata := server.extractMetaData(ctx)

	arg := db.CreateSessionParams{
		ID:           refreshTokenPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	}
	session, err := server.store.CreateSession(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to create session")
	}
	return &pb.LoginUserResponse{
		User:                  convertUser(user),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		SessionId:             session.ID.String(),
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}, nil
}

func validateLoginUserRequest(rq *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(rq.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validation.ValidatePassword(rq.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
