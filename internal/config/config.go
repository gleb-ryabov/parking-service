package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represent app settings.
type Config struct {
	Log   LogConfig  `yaml:"log"`
	DbURL string     `yaml:"db_url"`
	HTTP  HTTPConfig `yaml:"http"`
	GRPC  GRPCConfig `yaml:"grpc"`
}

// HTTPConfig represents HTTP config.
type HTTPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// GRPCConfig represents GRPC config.
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// LogConfig represents logging config.
type LogConfig struct {
	Level      string `yaml:"level"`
	FileName   string `yaml:"file"`
	MaxAge     int    `yaml:"maxAge"`
	MaxBackups int    `yaml:"maxBackups"`
}

// MustLoad loads config by env/flag path or panic.
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is not set")
	}

	return MustLoadByPath(path)
}

// MustLoadByPath loads config by path or panic.
func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed read config file: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath returns config path from flag or env variable.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
