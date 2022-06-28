package task

import (
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/internal/ports"
)

type Service struct {
	db ports.TaskDB
}

func New(db ports.TaskDB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) ListTasks() ([]*models.Task, error) {
	t, err := s.db.List()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) RunTask(createdTask *models.Task) error {
	err := s.db.Run(createdTask)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteTask(id string) error {
	err := s.db.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ApproveTask(id, login string) error {
	err := s.db.Approve(id, login)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeclineTask(id, login string) error {
	err := s.db.Decline(id, login)
	if err != nil {
		return err
	}
	return nil
}
