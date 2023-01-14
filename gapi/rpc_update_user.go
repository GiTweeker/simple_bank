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
	"time"
)

func (server *Server) UpdateUser(ctx context.Context, rq *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	if violations := validateUpdateUserRequest(rq); violations != nil {
		return nil, invalidArgument(violations)
	}

	if authPayload.Username != rq.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update a different username")
	}
	_, err = server.store.GetUser(ctx, rq.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound,
				"username does not exist : %s", rq.GetUsername())
		}

		return nil, status.Errorf(codes.Internal,
			"unable to find user", err)
	}

	arg := db.UpdateUserParams{
		FullName: sql.NullString{
			String: rq.GetFullName(),
			Valid:  rq.FullName != nil,
		},
		Email: sql.NullString{
			String: rq.GetEmail(),
			Valid:  rq.Email != nil,
		},
		Username: rq.Username,
	}

	if rq.Password != nil {
		hashedPassword, err := util.HashPassword(rq.GetPassword())

		if err != nil {
			return nil, status.Errorf(codes.Internal,
				"unable to hash password")
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

	}
	updatedUser, err := server.store.UpdateUser(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound,
				"user not found")
		}
		return nil, status.Errorf(codes.Internal,
			"failed to update user, error : %s", err)

	}

	return &pb.UpdateUserResponse{
		User: convertUser(updatedUser),
	}, nil
}

func validateUpdateUserRequest(rq *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(rq.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if rq.Password != nil {
		if err := validation.ValidatePassword(rq.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}

	}

	if rq.FullName != nil {
		if err := validation.ValidateFullName(rq.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if rq.Email != nil {
		if err := validation.ValidateEmail(rq.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
