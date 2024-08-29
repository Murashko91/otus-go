package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf `yaml:"logger"`
	Database DBConfig   `yaml:"db"`
	Server   Server     `yaml:"server"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
	InMemory bool   `yaml:"inMemory"`
}

type Server struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	HostGRPC string `yaml:"hostgrpc"`
	PortGRPC int    `yaml:"portgrpc"`
}

func NewCalendarConfig(configFilePath string) Config {
	conf := &Config{}

	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(conf); err != nil {
		panic(err)
	}

	return *conf
}

// TODO
