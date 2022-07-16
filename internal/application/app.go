package application

import (
	"context"
	"os"

	"gitlab.com/g6834/team26/task/internal/adapters/grpc"
	"gitlab.com/g6834/team26/task/internal/adapters/http"
	"gitlab.com/g6834/team26/task/internal/adapters/postgres"
	"gitlab.com/g6834/team26/task/internal/domain/task"
	"gitlab.com/g6834/team26/task/pkg/config"
	"gitlab.com/g6834/team26/task/pkg/logger"
	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog"
)

var (
	s        *http.Server
	l        *zerolog.Logger
	db       *postgres.PostgresDatabase
	grpcAuth *grpc.GrpcAuth
)

func Start(ctx context.Context) {
	l = logger.New()

	c, err := config.New()
	if err != nil {
		l.Error().Msgf("Error parsing env: %s", err)
	}

	db, err = postgres.New(ctx, c.Server.PgUrl)
	if err != nil {
		l.Error().Msgf("db init failed: %s", err)
		os.Exit(1)
	}

	// db, err := json_db.New(c.Server.JsonDbFile)
	// if err != nil {
	// 	l.Error().Msgf("json db init failed: %s", err)
	// 	os.Exit(1)
	// }

	grpcAuth, err = grpc.New(c.Server.GRPCAuth)
	if err != nil {
		l.Error().Msgf("grpc auth client init failed: %s", err)
		os.Exit(1)
	}

	grpcAnalytic, err := grpc.NewAnalytic(c.Server.GRPCAnalytic)
	if err != nil {
		l.Error().Msgf("grpc analytic client init failed: %s", err)
		os.Exit(1)
	}

	taskS := task.New(db, grpcAuth, grpcAnalytic)

	s, err = http.New(l, taskS, c)
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
	ctx := context.Background()

	_ = s.Stop(ctx)
	_ = db.Stop(ctx)
	_ = grpcAuth.Stop(ctx)
	l.Info().Msg("app has stopped")
}
