package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/pkg/api"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) List(login string) ([]*models.Task, error) {
	args := d.Called(login)
	return args.Get(0).([]*models.Task), args.Error(1)
}

func (d *DbMock) Run(t *models.Task) error {
	args := d.Called(t)
	return args.Error(0)
}

func (d *DbMock) Delete(login, id string) error {
	args := d.Called(login, id)
	return args.Error(0)
}

func (d *DbMock) Approve(login, id, approvalLogin string) error {
	args := d.Called(login, id, approvalLogin)
	return args.Error(0)
}

func (d *DbMock) Decline(login, id, approvalLogin string) error {
	args := d.Called(login, id, approvalLogin)
	return args.Error(0)
}

type GrpcMock struct {
	mock.Mock
}

func (g *GrpcMock) Validate(refreshCookie, accessCookie string) (*api.AuthResponse, error) {
	args := g.Called(refreshCookie, accessCookie)
	return args.Get(0).(*api.AuthResponse), args.Error(1)
}
