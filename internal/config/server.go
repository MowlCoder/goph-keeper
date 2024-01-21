package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
)

type Server struct {
	HTTPAddr    string `env:"HTTP_ADDR" json:"http_addr"`
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`

	LogPassSecret string `env:"LOG_PASS_SECRET" json:"log_pass_secret"`
}

func (s *Server) Parse() {
	flag.StringVar(&s.HTTPAddr, "http", ":4000", "Server will listen this http address")
	flag.StringVar(&s.DatabaseDSN, "dns", "", "Database connect uri")
	flag.StringVar(&s.LogPassSecret, "lp-secret", "secretttsecretttsecretttsecrettt", "Secret for crypt log pass data")

	flag.Parse()

	if err := env.Parse(s); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
