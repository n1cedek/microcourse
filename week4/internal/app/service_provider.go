package app

import (
	"context"
	"log"
	notea "microservices_course/week4/internal/api/note"
	"microservices_course/week4/internal/client/db"
	"microservices_course/week4/internal/client/db/pg"
	"microservices_course/week4/internal/client/db/transaction"
	"microservices_course/week4/internal/closer"
	"microservices_course/week4/internal/config/env"
	repository "microservices_course/week4/internal/repo"
	noteRepo "microservices_course/week4/internal/repo/note"
	"microservices_course/week4/internal/service"
	noteService "microservices_course/week4/internal/service/note"
)

type serviceProvider struct {
	pgConfig    env.PGConfig
	grpcConfig  env.GPRCConfig
	dbClient    db.Client
	txManager   db.TxManager
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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}
	return s.dbClient
}
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}
func (s *serviceProvider) NoteRepo(ctx context.Context) repository.NoteRepo {
	if s.noteRepo == nil {
		s.noteRepo = noteRepo.NewRepo(s.DBClient(ctx))
	}
	return s.noteRepo
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(s.NoteRepo(ctx), s.TxManager(ctx))
	}
	return s.noteService
}

func (s *serviceProvider) NoteImpl(ctx context.Context) *notea.Implementation {
	if s.noteImpl == nil {
		s.noteImpl = notea.NewImplementation(s.NoteService(ctx))
	}
	return s.noteImpl
}
