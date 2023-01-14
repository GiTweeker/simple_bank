package gapi

import (
	"context"
	"fmt"
	"github.com/techschool/simple-bank/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorisationHeader     = "authorization"
	authorisationHeaderType = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorisationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorisation header format")
	}
	authType := strings.ToLower(fields[0])

	if authType != authorisationHeaderType {
		return nil, fmt.Errorf("unsupported authorisation header Type")
	}

	accessToken := fields[1]

	payload, err := server.tokenMaker.VerifyToken(accessToken)

	if err != nil {
		return nil, fmt.Errorf("invalid access token : %s", err)
	}

	return payload, nil

}