package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
)

// Client - struct responsible for storing client config
type Client struct {
	ServerBaseAddr   string `env:"SERVER_BASE_ADDR" json:"server_base_addr"`
	ApiServerTimeout int    `env:"API_SERVER_TIMEOUT" json:"api_server_timeout"`
}

// Parse - parse client config from flags and envs
func (s *Client) Parse() {
	flag.StringVar(&s.ServerBaseAddr, "server", "", "Base http server address")
	flag.IntVar(&s.ApiServerTimeout, "api-timeout", 60, "Api server timeout in seconds")

	flag.Parse()

	if err := env.Parse(s); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
