package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/internal/ports"
	"gitlab.com/g6834/team26/task/pkg/api"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) List(ctx context.Context, login string) ([]*models.Task, error) {
	args := d.Called(ctx, login)
	return args.Get(0).([]*models.Task), args.Error(1)
}

func (d *DbMock) Run(ctx context.Context, t *models.Task) error {
	args := d.Called(ctx, t)
	return args.Error(0)
}

func (d *DbMock) Update(ctx context.Context, id, login, name, text string) error {
	args := d.Called(ctx, id, login, name, text)
	return args.Error(0)
}

func (d *DbMock) Delete(ctx context.Context, login, id string) error {
	args := d.Called(ctx, login, id)
	return args.Error(0)
}

func (d *DbMock) Approve(ctx context.Context, login, id, approvalLogin string) error {
	args := d.Called(ctx, login, id, approvalLogin)
	return args.Error(0)
}

func (d *DbMock) Decline(ctx context.Context, login, id, approvalLogin string) error {
	args := d.Called(ctx, login, id, approvalLogin)
	return args.Error(0)
}

type GrpcAuthMock struct {
	mock.Mock
}

func (g *GrpcAuthMock) Validate(ctx context.Context, tokens ports.TokenPair) (*api.AuthResponse, error) {
	args := g.Called(ctx, tokens)
	return args.Get(0).(*api.AuthResponse), args.Error(1)
}

type GrpcAnalyticMock struct {
	mock.Mock
}

func (g *GrpcAnalyticMock) AddTask(ctx context.Context, t *models.Task) error {
	args := g.Called(ctx, t)
	return args.Error(0)
}

func (g *GrpcAnalyticMock) ActionTask(ctx context.Context, u, e, a string, v bool) error {
	args := g.Called(ctx, u, e, a, v)
	return args.Error(0)
}
