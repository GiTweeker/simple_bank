package gapi

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type MetaData struct {
	UserAgent string
	ClientIp  string
}

const (
	MetadataForwardedForHost = "x-forwarded-for"
	MetadataGatewayUserAgent = "grpcgateway-user-agent"
	MetadataGUserAgent       = "user-agent"
)

func (server *Server) extractMetaData(ctx context.Context) *MetaData {

	getMetaDataValues(ctx, "", "")

	userAgent := getMetaDataValues(ctx, MetadataGatewayUserAgent, "")

	if len(userAgent) < 1 {
		userAgent = getMetaDataValues(ctx, MetadataGUserAgent, "")
	}
	clientIP := getMetaDataValues(ctx, MetadataForwardedForHost, "")

	if len(clientIP) < 1 {
		p, ok := peer.FromContext(ctx)

		if ok {
			clientIP = p.Addr.String()
		}
	}
	return &MetaData{
		UserAgent: userAgent,
		ClientIp:  clientIP,
	}
}

func getMetaDataValues(ctx context.Context, key string, defaultValue string) string {

	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("unable to get meta data from context for key " + key)
		return defaultValue
	}

	if v := m.Get(key); len(v) > 0 {
		return m.Get(key)[0]
	}

	return defaultValue

}
