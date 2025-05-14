package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"microservices_course/week6/jwt/internal/model"
	"microservices_course/week6/jwt/internal/utils"
	descAccess "microservices_course/week6/jwt/pkg/access_v1"
	descAuth "microservices_course/week6/jwt/pkg/auth_v1"
	"net"
	"strings"
	"time"
)

const (
	grpcPort               = 50051
	authPrefix             = "Bearer "
	refreshTokenSecretKey  = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey   = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 1 * time.Minute
)

var accessibleRoles map[string]string

type serverAuth struct {
	descAuth.UnimplementedAuthV1Server
}

func (s *serverAuth) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	refToken, err := utils.GenerateToken(model.UserInfo{
		Username: req.GetUsername(),
		Role:     "admin",
	}, []byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &descAuth.LoginResponse{RefreshToken: refToken}, nil
}
func (s *serverAuth) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
	clim, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to refresh token")
	}

	refToken, err := utils.GenerateToken(model.UserInfo{
		Username: clim.Username,
		Role:     "admin",
	}, []byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}
	return &descAuth.GetRefreshTokenResponse{RefreshToken: refToken}, nil
}

func (s *serverAuth) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
	clim, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to refresh token")
	}
	acsToken, err := utils.GenerateToken(model.UserInfo{
		Username: clim.Username,
		Role:     "admin",
	}, []byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		return nil, err
	}
	return &descAuth.GetAccessTokenResponse{AccessToken: acsToken}, nil
}

type serverAccess struct {
	descAccess.UnimplementedAccessV1Server
}

func (s *serverAccess) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	// 1. Извлекаем метаданные из контекста запроса
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	// 2. Проверяем наличие заголовка "authorization"
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	// 3. Проверяем, начинается ли токен с нужного префикса (например, "Bearer ")
	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	// 4. Убираем префикс, оставляя только сам токен
	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
	if err != nil {
		return nil, errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return nil, errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[req.GetEndpointAddress()]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	if role == claims.Role {
		return &emptypb.Empty{}, nil
	}

	return nil, errors.New("access denied")
}

func (s *serverAccess) accessibleRoles(ctx context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)
		accessibleRoles[model.ExamplePath] = "admin"
	}
	return accessibleRoles, nil
}

//func (s *server) Get(ctx context.Context, req *note_v1.GetRequest) (*note_v1.GetResponse, error) {
//	log.Printf("Get ID: %v", req.GetId())
//
//	return &note_v1.GetResponse{
//		Note: &note_v1.Note{
//			Id: req.GetId(),
//			Info: &note_v1.NoteInfo{
//				Title:    gofakeit.BeerName(),
//				Content:  gofakeit.IPv4Address(),
//				Author:   gofakeit.Name(),
//				IsPublic: gofakeit.Bool(),
//			},
//			CreatedAt: timestamppb.New(gofakeit.Date()),
//			UpdatedAt: timestamppb.New(gofakeit.Date()),
//		},
//	}, nil
//}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen server: %v", err)
	}

	cred, err := credentials.NewServerTLSFromFile("/home/n1cedek/GolandProjects/microcourse/week6/grpc/service.pem", "/home/n1cedek/GolandProjects/microcourse/week6/grpc/service.key")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(cred))
	reflection.Register(s)
	descAuth.RegisterAuthV1Server(s, &serverAuth{})
	descAccess.RegisterAccessV1Server(s, &serverAccess{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
