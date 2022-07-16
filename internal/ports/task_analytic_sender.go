package ports

import (
	"context"

	"gitlab.com/g6834/team26/task/internal/domain/models"
)

type TaskAnalyticSender interface {
	AddTask(ctx context.Context, t *models.Task) error
	ActionTask(ctx context.Context, u, e, a string, v bool) error
}
