package env

import (
	"errors"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
)

const (
	tlsFile = "TLS_CERT_FILE"
	tlsKey  = "TLS_KEY_FILE"
)

var _ TLSConfig = (*tlsConfig)(nil)

type tlsConfig struct {
	certFile string
	keyFile  string
}

func NewTlsConfig() (*tlsConfig, error) {
	cF := os.Getenv(tlsFile)
	if len(cF) == 0 {
		return nil, errors.New("failed to get tls file")
	}
	cK := os.Getenv(tlsKey)
	if len(cK) == 0 {
		return nil, errors.New("failed to get tls key")
	}
	return &tlsConfig{
		certFile: cF,
		keyFile:  cK,
	}, nil
}

func (t *tlsConfig) GetTLSConfig() (credentials.TransportCredentials, error) {
	cred, err := credentials.NewServerTLSFromFile(t.certFile, t.keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
		return nil, err
	}
	return cred, nil
}
