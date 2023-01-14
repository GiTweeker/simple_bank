package gapi

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/pb"
	"github.com/techschool/simple-bank/util"
	"github.com/techschool/simple-bank/validation"
	"github.com/techschool/simple-bank/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	if violations := validateCreateUserRequest(rq); violations != nil {
		return nil, invalidArgument(violations)
	}
	hashedPassword, err := util.HashPassword(rq.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to has password : %s", err)
	}

	arg := db.CreateUserParams{
		Username:       rq.GetUsername(),
		FullName:       rq.GetFullName(),
		Email:          rq.GetEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists,
					"username already exist : %s", rq.GetUsername())

			}
		}
		return nil, status.Errorf(codes.Internal,
			"failed to create user, error : %s", err)

	}

	payloadForVerifyEmail := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	opt := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payloadForVerifyEmail, opt...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send verify email"+err.Error())
	}

	rsp := &pb.CreateUserResponse{User: convertUser(user)}

	return rsp, nil

}

func validateCreateUserRequest(rq *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(rq.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validation.ValidatePassword(rq.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := validation.ValidateFullName(rq.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := validation.ValidateEmail(rq.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return violations
}
