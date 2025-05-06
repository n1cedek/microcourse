package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	notea "microservices_course/week3/internal/api/note"
	"microservices_course/week3/internal/closer"
	"microservices_course/week3/internal/config/env"
	repository "microservices_course/week3/internal/repo"
	noteRepo "microservices_course/week3/internal/repo/note"
	"microservices_course/week3/internal/service"
	noteService "microservices_course/week3/internal/service/note"
)

type serviceProvider struct {
	pgConfig   env.PGConfig
	grpcConfig env.GPRCConfig
	pgPool     *pgxpool.Pool

	noteService service.NoteService
	noteRepo    repository.NoteRepo
	noteImpl    *notea.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() env.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewDsnConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() env.GPRCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGrpcConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to dstabase: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping: %v", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})
		s.pgPool = pool
	}
	return s.pgPool
}

func (s *serviceProvider) NoteRepo(ctx context.Context) repository.NoteRepo {
	if s.noteRepo == nil {
		s.noteRepo = noteRepo.NewRepo(s.PgPool(ctx))
	}
	return s.noteRepo
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(s.NoteRepo(ctx))
	}
	return s.noteService
}

func (s *serviceProvider) NoteImpl(ctx context.Context) *notea.Implementation {
	if s.noteImpl == nil {
		s.noteImpl = notea.NewImplementation(s.NoteService(ctx))
	}
	return s.noteImpl
}
