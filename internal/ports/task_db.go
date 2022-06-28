package ports

import "gitlab.com/g6834/team26/task/internal/domain/models"

type TaskDB interface {
	List(login string) ([]*models.Task, error)
	Run(t *models.Task) error
	Delete(login, id string) error
	Approve(login, id, approvalLogin string) error
	Decline(login, id, approvalLogin string) error
}
