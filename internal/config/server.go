package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
)

type Server struct {
	HTTPAddr    string `env:"HTTP_ADDR" json:"http_addr"`
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	SSLPemPath  string `env:"SSL_PEM_PATH" json:"ssl_pem_path"`
	SSLKeyPath  string `env:"SSL_KEY_PATH" json:"ssl_key_path"`

	DataSecretKey string `env:"DATA_SECRET_KEY" json:"data_secret_key"`
}

func (s *Server) Parse() {
	flag.StringVar(&s.HTTPAddr, "http", ":4000", "Server will listen this http address")
	flag.StringVar(&s.DatabaseDSN, "dns", "", "Database connect uri")
	flag.BoolVar(&s.EnableHTTPS, "https", false, "If true, server will use HTTPS")
	flag.StringVar(&s.SSLPemPath, "pp", "", "Path to SSL pem file")
	flag.StringVar(&s.SSLKeyPath, "kp", "", "Path to SSL key file")
	flag.StringVar(&s.DataSecretKey, "data-secret", "secretttsecretttsecretttsecrettt", "Secret for crypt data")

	flag.Parse()

	if err := env.Parse(s); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
