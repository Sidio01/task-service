package ports

import "gitlab.com/g6834/team26/task/pkg/api"

type Grpc interface {
	Validate(refreshCookie, accessCookie string) (*api.AuthResponse, error)
}
