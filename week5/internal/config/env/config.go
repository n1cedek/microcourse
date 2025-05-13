package env

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc/credentials"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}

type PGConfig interface {
	DSN() string
}

type GPRCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}
type TLSConfig interface {
	GetTLSConfig() (credentials.TransportCredentials, error)
}
