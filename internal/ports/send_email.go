package ports

import "gitlab.com/g6834/team26/task/internal/domain/models"

type EmailSender interface {
	StartEmailWorkers()
	SendEmail(e models.Email) error
	PushEmailToChan(e models.Email)
	GetEmailResultChan() chan map[models.Email]bool
}
