package ports

import (
	"context"
)

type TaskAnalyticSender interface {
	ActionTask(ctx context.Context, u, t, v string) error
}
