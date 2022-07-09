package ports

import "gitlab.com/g6834/team26/task/internal/domain/models"

type Task interface {
	// Info(ctx context.Context, login string) (*models.User, error)
	// Validate(ctx context.Context, tokens models.TokenPair) (string, error)
	// Login(ctx context.Context, user, password string) (models.TokenPair, error)
	ListTasks(login string) ([]*models.Task, error) // TODO: передавать контекст
	RunTask(createdTask *models.Task) error
	DeleteTask(login, id string) error
	ApproveTask(login, id, approvalLogin string) error
	DeclineTask(login, id, approvalLogin string) error
	GrpcAuth
}
