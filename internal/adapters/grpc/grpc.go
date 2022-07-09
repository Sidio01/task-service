package grpc

import (
	"context"
	"time"

	"gitlab.com/g6834/team26/task/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcAuth struct {
	GrpcClient api.AuthClient
}

func New(url string) (*GrpcAuth, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	return &GrpcAuth{GrpcClient: api.NewAuthClient(conn)}, nil
}

func (Grpc *GrpcAuth) Validate(refreshCookie, accessCookie string) (*api.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	authReq := &api.AuthRequest{Service: "task", AccessToken: accessCookie, RefreshToken: refreshCookie}
	// log.Println(authReq)
	grpcResponse, err := Grpc.GrpcClient.VerifyToken(ctx, authReq)
	if err != nil {
		return nil, err
	}
	return grpcResponse, nil
}
