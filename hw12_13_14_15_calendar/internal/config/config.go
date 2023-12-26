package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf `toml:"logger"`
	Database   DBConf     `toml:"database"`
	HTTPServer HTTP       `toml:"http-server"`
	GRPCServer HTTP       `toml:"grpc-server"`
}

type LoggerConf struct {
	Level       string `toml:"level"`
	LogEncoding string `toml:"log_encoding"`
}

type DBConf struct {
	Storage       string `toml:"storage"`
	ConnectString string `toml:"connect_str"`
}

type HTTP struct {
	Address string `toml:"address"`
}

type GRPC struct {
	Address string `toml:"address"`
}

func NewConfig(configFile string) Config {
	var c Config
	tomlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("error reading the configuration file")
	}
	_, err = toml.Decode(string(tomlFile), &c)
	if err != nil {
		log.Fatal("error decode the configuration file")
	}
	if c.Logger.Level == "" {
		c.Logger.Level = "ERROR"
	}
	if c.Logger.LogEncoding == "" {
		c.Logger.LogEncoding = "console"
	}
	if c.Database.Storage == "" {
		c.Database.Storage = "memory"
	}
	return c
}
