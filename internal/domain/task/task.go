package task

import (
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/internal/ports"
	"gitlab.com/g6834/team26/task/pkg/api"
)

type Service struct {
	db             ports.TaskDB
	grpcAuth       ports.GrpcAuth
	analyticSender ports.TaskAnalyticSender
}

func New(db ports.TaskDB, grpcAuth ports.GrpcAuth, analyticSender ports.TaskAnalyticSender) *Service {
	return &Service{
		db:             db,
		grpcAuth:       grpcAuth,
		analyticSender: analyticSender,
	}
}

func (s *Service) ListTasks(login string) ([]*models.Task, error) {
	t, err := s.db.List(login)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) RunTask(createdTask *models.Task) error { // TODO: отправлять письмо первому согласующему
	err := s.db.Run(createdTask)
	if err != nil {
		return err
	}

	err = s.analyticSender.AddTask(createdTask)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteTask(login, id string) error { // TODO: отправлять письма всем участникам об отмене операции
	err := s.db.Delete(login, id)
	if err != nil {
		return err
	}

	err = s.analyticSender.ActionTask(id, "", "delete", true)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ApproveTask(login, id, approvalLogin string) error { // TODO: отправлять письмо следующему согласующему
	err := s.db.Approve(login, id, approvalLogin)
	if err != nil {
		return err
	}

	err = s.analyticSender.ActionTask(id, approvalLogin, "approve", true)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeclineTask(login, id, approvalLogin string) error { // TODO: отправлять письма всем участникам об остановке согласования операции в связи с отклонением одним из участников
	err := s.db.Decline(login, id, approvalLogin)
	if err != nil {
		return err
	}

	err = s.analyticSender.ActionTask(id, approvalLogin, "approve", false)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Validate(refreshCookie, accessCookie string) (*api.AuthResponse, error) {
	grpcResponse, err := s.grpcAuth.Validate(refreshCookie, accessCookie)
	if err != nil {
		return nil, err
	}
	return grpcResponse, nil
}
