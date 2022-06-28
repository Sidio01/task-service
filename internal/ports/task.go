package ports

import "gitlab.com/g6834/team26/task/internal/domain/models"

type Task interface {
	// Info(ctx context.Context, login string) (*models.User, error)
	// Validate(ctx context.Context, tokens models.TokenPair) (string, error)
	// Login(ctx context.Context, user, password string) (models.TokenPair, error)
	ListTasks() ([]*models.Task, error)
	RunTask(createdTask *models.Task) error
	DeleteTask(id string) error
	ApproveTask(id, login string) error
	DeclineTask(id, login string) error
}
