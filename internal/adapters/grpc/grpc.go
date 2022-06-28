package grpc

import (
	"context"
	"time"

	"gitlab.com/g6834/team26/task/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Grpc struct {
	GrpcClient api.AuthClient
}

func New(url string) (*Grpc, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	return &Grpc{GrpcClient: api.NewAuthClient(conn)}, nil
}

func (Grpc *Grpc) Validate(refreshCookie, accessCookie string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	authReq := &api.AuthRequest{Service: "task", AccessToken: accessCookie, RefreshToken: refreshCookie}
	// log.Println(authReq)
	grpcResponse, err := Grpc.GrpcClient.VerifyToken(ctx, authReq)
	if err != nil {
		return false, "", err
	}
	return grpcResponse.Result, grpcResponse.Login, nil
}
