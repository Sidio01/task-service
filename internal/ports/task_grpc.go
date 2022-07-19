package ports

import (
	"context"

	"gitlab.com/g6834/team26/task/pkg/api"
)

type GrpcAuth interface {
	Validate(ctx context.Context, refreshCookie, accessCookie string) (*api.AuthResponse, error)
}
