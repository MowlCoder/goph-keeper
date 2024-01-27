package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
)

type Server struct {
	HTTPAddr    string `env:"HTTP_ADDR" json:"http_addr"`
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`

	DataSecretKey string `env:"DATA_SECRET_KEY" json:"data_secret_key"`
}

func (s *Server) Parse() {
	flag.StringVar(&s.HTTPAddr, "http", ":4000", "Server will listen this http address")
	flag.StringVar(&s.DatabaseDSN, "dns", "", "Database connect uri")
	flag.StringVar(&s.DataSecretKey, "data-secret", "secretttsecretttsecretttsecrettt", "Secret for crypt data")

	flag.Parse()

	if err := env.Parse(s); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
