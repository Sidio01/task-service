package ports

import (
	"context"

	"gitlab.com/g6834/team26/task/internal/domain/models"
)

type Task interface {
	// Info(ctx context.Context, login string) (*models.User, error)
	// Validate(ctx context.Context, tokens models.TokenPair) (string, error)
	// Login(ctx context.Context, user, password string) (models.TokenPair, error)
	ListTasks(ctx context.Context, login string) ([]*models.Task, error) // TODO: передавать контекст
	RunTask(ctx context.Context, createdTask *models.Task) error
	DeleteTask(ctx context.Context, login, id string) error
	ApproveTask(ctx context.Context, login, id, approvalLogin string) error
	DeclineTask(ctx context.Context, login, id, approvalLogin string) error
	GrpcAuth
}
