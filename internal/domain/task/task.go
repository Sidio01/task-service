package task

import (
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/internal/ports"
)

type Service struct {
	db   ports.TaskDB
	grpc ports.Grpc
}

func New(db ports.TaskDB, grpc ports.Grpc) *Service {
	return &Service{
		db:   db,
		grpc: grpc,
	}
}

func (s *Service) ListTasks(login string) ([]*models.Task, error) {
	t, err := s.db.List(login)
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

func (s *Service) DeleteTask(login, id string) error {
	err := s.db.Delete(login, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ApproveTask(login, id, approvalLogin string) error {
	err := s.db.Approve(login, id, approvalLogin)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeclineTask(login, id, approvalLogin string) error {
	err := s.db.Decline(login, id, approvalLogin)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Validate(refreshCookie, accessCookie string) (bool, string, error) {
	result, login, err := s.grpc.Validate(refreshCookie, accessCookie)
	if err != nil {
		return false, "", err
	}
	return result, login, nil
}
