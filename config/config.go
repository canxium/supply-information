package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/joho/godotenv"
)

// App config struct
type Config struct {
	Server ServerConfig
	Supply Supply
}

// Server config struct
type ServerConfig struct {
	AppVersion        string        `json:"app_version" conf:"default:1.0.0,env:APP_VERSION"`
	Addr              string        `json:"app_port" conf:"default:0.0.0.0:8080,env:APP_ADDR"`
	Mode              string        `json:"app_mode" conf:"default:production,env:APP_MODE"`
	ReadTimeout       time.Duration `json:"read_timeout" conf:"default:5s,env:APP_READ_TIMEOUT"`
	WriteTimeout      time.Duration `json:"write_timeout" conf:"default:5s,env:APP_WRITE_TIMEOUT"`
	CtxDefaultTimeout time.Duration `json:"ctx_timeout" conf:"default:12s,env:APP_CTX_TIMEOUT"`
	Postgres          string        `json:"postgres" conf:"env:POSTGRES"`
}

type Supply struct {
	RpcApi    string `json:"rpc_api" conf:"default:https://rpc.canxium.org,env:RPC_API_URL"`
	Addresses string `json:"addresses" conf:"default:0xBd65D6efb2C3e6B4dD33C664643BEB8e5E133055,env:ADDRESSES"`
}

// Parse config file
func LoadConfig() (*Config, error) {
	godotenv.Load()
	cfg := &Config{}

	if help, err := conf.Parse("", cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
		}

		log.Fatal(err)
	}

	return cfg, nil
}
