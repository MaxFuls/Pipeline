package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env"`
	Name     string         `yaml:"name"`
	Registry RegistryServer `yaml:"registry"`
	GRPC     GRPCServer     `ymal:"grpc"`
}

type RegistryServer struct {
	Ip   string `yaml:"ip"`
	Port uint32 `yaml:"port"`
}

type GRPCServer struct {
	Ip      string        `yaml:"ip"`
	Port    uint32        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path invalid " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("LOG_CONFIG_PATH")
	}

	return res
}
