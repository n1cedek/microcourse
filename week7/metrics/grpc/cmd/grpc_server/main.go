package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"microservices_course/week7/grpcl/grpc/pkg/note_v1"
	"microservices_course/week7/metrics/grpc/internal/interceptor"
	"microservices_course/week7/metrics/grpc/internal/logger"
	"microservices_course/week7/metrics/grpc/internal/metric"
	"net"
	"net/http"
	"os"
	"time"
)

var logLevel = flag.String("l", "info", "log level")

const grpcPort = 50051

type server struct {
	note_v1.UnimplementedNoteV1Server
}

func (s *server) Get(ctx context.Context, req *note_v1.GetRequest) (*note_v1.GetResponse, error) {
	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}

	// rand.Intn(max - min) + min
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	logger.Info("Getting note...", zap.Int64("id", req.GetId()))

	return &note_v1.GetResponse{
		Note: &note_v1.Note{
			Id: req.GetId(),
			Info: &note_v1.NoteInfo{
				Title:    gofakeit.BeerName(),
				Content:  gofakeit.IPv4Address(),
				Author:   gofakeit.Name(),
				IsPublic: gofakeit.Bool(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	ctx := context.Background()
	err := metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen server: %v", err)
	}

	logger.Init(getCore(getAtomicLevel()))

	//s := grpc.NewServer(grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
	//	interceptor.LogInterceptor,
	//	interceptor.ValidateInterceptor,
	//),
	//),
	//)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.MetricsInterceptor))
	reflection.Register(s)
	note_v1.RegisterNoteV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		err = runPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
	}
	log.Printf("Prometheus server is running on %s", "localhost:2112")
	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
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
