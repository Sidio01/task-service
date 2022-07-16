package ports

import "gitlab.com/g6834/team26/task/internal/domain/models"

type TaskAnalyticSender interface {
	AddTask(t *models.Task) error
	ActionTask(u, e, a string, v bool) error
}
