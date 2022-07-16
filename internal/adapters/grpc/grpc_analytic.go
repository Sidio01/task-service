package grpc

import (
	"context"
	"time"

	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/pkg/api"
	"gitlab.com/g6834/team26/task/pkg/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcAnalytic struct {
	GrpcClient api.AnalyticClient
	GrpcConn   *grpc.ClientConn
}

func NewAnalytic(url string) (*GrpcAnalytic, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	return &GrpcAnalytic{
		GrpcClient: api.NewAnalyticClient(conn),
		GrpcConn:   conn,
	}, nil
}

func (GrpcAnalytic *GrpcAnalytic) AddTask(ctx context.Context, t *models.Task) error {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	emails := make([]string, len(t.Approvals))
	for idx, email := range t.Approvals {
		emails[idx] = email.ApprovalLogin
	}

	addTaskReq := &api.AddTaskRequest{UUID: t.UUID,
		Login:       t.InitiatorLogin,
		Timestamp:   time.Now().Unix(),
		Emails:      emails,
		UUIDMessage: uuid.GenUUID(),
	}
	// log.Println(addTaskReq)
	_, err := GrpcAnalytic.GrpcClient.AddTask(ctx, addTaskReq)
	if err != nil {
		return err
	}
	return nil
}

func (GrpcAnalytic *GrpcAnalytic) ActionTask(ctx context.Context, u, e, a string, v bool) error {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	actionTaskReq := &api.ActionTaskRequest{UUID: u,
		Timestamp:   time.Now().Unix(),
		Email:       e,
		Action:      a,
		Value:       v,
		UUIDMessage: uuid.GenUUID(),
	}
	// log.Println(addTaskReq)
	_, err := GrpcAnalytic.GrpcClient.ActionTask(ctx, actionTaskReq)
	if err != nil {
		return err
	}
	return nil
}

func (GrpcAnalytic *GrpcAnalytic) StopAnalytic(ctx context.Context) error {
	err := GrpcAnalytic.GrpcConn.Close()
	if err != nil {
		return err
	}
	return nil
}
