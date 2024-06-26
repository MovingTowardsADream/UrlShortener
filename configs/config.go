package configs

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

const (
	defaultConfigPath = "./configs/config.yaml"
	defaultEnvPath    = ".env"
)

type (
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Redis `yaml:"redis"`
		Mongo `yaml:"mongo"`
		Log   `yaml:"logger"`
	}

	App struct {
		Name    string `yaml:"name" default:"url-shortener"`
		Version string `yaml:"version" default:"v1.0.0"`
	}

	HTTP struct {
		Port    string        `env:"HTTP_PORT" env-default:":8080" yaml:"port"`
		Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"5s"    yaml:"timeout"`
	}

	Redis struct {
		Address  string `env:"REDIS_ADDRESS"`
		Password string `env:"REDIS_PASSWORD"`
		PoolMax  int    `env:"REDIS_POOL_MAX" env-default:"100"     yaml:"pool_max"`
		IdleConn int    `env:"REDIS_IDLE_CONN" env-default:"10"     yaml:"idle_conn"`
	}

	Mongo struct {
		Url      string `env:"MONGO_URL"`
		PoolMax  int    `env:"MONGO_POOL_MAX" env-default:"100"     yaml:"pool_max"`
		IdleConn int    `env:"MONGO_IDLE_CONN" env-default:"10"     yaml:"idle_conn"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug" yaml:"log_level"`
	}
)

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath, defaultEnvPath)
}

func MustLoadPath(configPath, envPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	_ = godotenv.Load(envPath)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = defaultConfigPath
	}

	return res
}
