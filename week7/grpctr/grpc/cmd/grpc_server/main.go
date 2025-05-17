package main

import (
	"context"
	"flag"
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/natefinch/lumberjack"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"microservices_course/week7/grpctr/grpc/internal/client"
	otherService "microservices_course/week7/grpctr/grpc/internal/client/other_service"
	"microservices_course/week7/grpctr/grpc/internal/interceptor"
	"microservices_course/week7/grpctr/grpc/internal/logger"
	"microservices_course/week7/grpctr/grpc/internal/tracing"
	"microservices_course/week7/grpctr/grpc/pkg/note_v1"
	descOther "microservices_course/week7/grpctr/grpc/pkg/other_note_v1"
	"net"
	"os"
	"time"
)

var logLevel = flag.String("l", "info", "log level")

const (
	grpcPort         = 50051
	otherServicePort = 50052
	serviceName      = "test-service"
)

type server struct {
	note_v1.UnimplementedNoteV1Server
	otherServiceClient client.OtherServiceClient
}

func (s *server) Get(ctx context.Context, req *note_v1.GetRequest) (*note_v1.GetResponse, error) {
	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	span, ctx := opentracing.StartSpanFromContext(ctx, "get note")
	defer span.Finish()

	span.SetTag("id", req.GetId())

	note, err := s.otherServiceClient.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.WithMessage(err, "getting note")
	}

	var updatedAt *timestamppb.Timestamp
	if note.UpdatedAt.Valid {
		updatedAt = timestamppb.New(note.UpdatedAt.Time)
	}
	//logger.Info("Getting note...", zap.Int64("id", req.GetId()))

	return &note_v1.GetResponse{
		Note: &note_v1.Note{
			Id: note.ID,
			Info: &note_v1.NoteInfo{
				Title:   note.Info.Title,
				Content: note.Info.Content,
			},
			CreatedAt: timestamppb.New(note.CreatedAt),
			UpdatedAt: updatedAt,
		},
	}, nil
}

func main() {
	flag.Parse()

	logger.Init(getCore(getAtomicLevel()))
	tracing.Init(logger.Logger(), serviceName)

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", otherServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	otherServiceClient := otherService.New(descOther.NewOtherNoteV1Client(conn))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen server: %v", err)
	}

	logger.Init(getCore(getAtomicLevel()))

	s := grpc.NewServer(grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
		interceptor.LogInterceptor,
		interceptor.ValidateInterceptor,
		interceptor.ServerTracingInterceptor,
	),
	),
	)
	reflection.Register(s)
	note_v1.RegisterNoteV1Server(s, &server{otherServiceClient: otherServiceClient})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	prodCfg := zap.NewProductionEncoderConfig()
	prodCfg.TimeKey = "timestamp"
	prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	develCfg := zap.NewDevelopmentEncoderConfig()
	develCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(develCfg)
	fileEnc := zapcore.NewJSONEncoder(prodCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEnc, file, level),
	)

}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
