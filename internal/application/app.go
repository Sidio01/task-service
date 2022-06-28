package application

import (
	"context"
	"os"

	"gitlab.com/g6834/team26/task/internal/adapters/http"
	"gitlab.com/g6834/team26/task/internal/adapters/json_db"
	"gitlab.com/g6834/team26/task/internal/domain/task"
	"gitlab.com/g6834/team26/task/pkg/getenv"
	"gitlab.com/g6834/team26/task/pkg/logger"
	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog"
)

var (
	s *http.Server
	l *zerolog.Logger
)

func Start(ctx context.Context) {
	l = logger.New()

	// TODO: заменить на postgresql
	// pgconn := os.Getenv("PG_URL")
	// db, err := postgres.New(ctx, pgconn)
	// if err != nil {
	// 	l.Error().Msgf("db init failed: %s", err)
	// 	os.Exit(1)
	// }

	jsonconn := getenv.GetEnv("JSON_DB_FILE", "db.jsonl")
	db, err := json_db.New(jsonconn)
	if err != nil {
		l.Error().Msgf("db init failed: %s", err)
		os.Exit(1)
	}

	taskS := task.New(db)

	s, err = http.New(l, taskS)
	if err != nil {
		l.Error().Msgf("http server creating failed: %s", err)
		os.Exit(1)
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.Start()
	})

	l.Info().Msg("app is started")
	err = g.Wait()
	if err != nil {
		l.Error().Msgf("http server start failed: %s", err)
		os.Exit(1)
	}
}

func Stop() {
	_ = s.Stop(context.Background())
	l.Info().Msg("app has stopped")
}
