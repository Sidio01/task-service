package ports

import "gitlab.com/g6834/team26/task/internal/domain/models"

type TaskDB interface {
	List() ([]*models.Task, error)
	Run(*models.Task) error
	Delete(id string) error
	Approve(id, login string) error
	Decline(id, login string) error
}
