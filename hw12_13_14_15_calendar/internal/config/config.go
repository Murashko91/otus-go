package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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

type SchedulerConf struct {
	Scheduler Scheduler  `yaml:"mb"`
	Logger    LoggerConf `yaml:"logger"`
	Database  DBConfig   `yaml:"db"`
}

type Scheduler struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	UserName       string `yaml:"user"`
	Password       string `yaml:"password"`
	Exchange       string `yaml:"exchange"`
	ExchangeType   string `yaml:"exchangeType"`
	RoutingKey     string `yaml:"routingKey"`
	IntervalCheck  int    `yaml:"intervalCheck"`
	NotifyInterval int    `yaml:"notifyInterval"`
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

func NewSchedulerConf(configFilePath string) SchedulerConf {
	conf := &SchedulerConf{}

	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(conf); err != nil {
		panic(err)
	}

	fmt.Println(conf)

	return *conf
}

// TODO
