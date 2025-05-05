package config

import "github.com/joho/godotenv"

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}

type PGConfig interface 
type PGConfig interface {
	DSN() string
}

type GPRCConfig interface {
	Address() string
}
