package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `yaml:"logger"`
	// TODO
}

type LoggerConf struct {
	Level string `yaml:"level"`
	// TODO
}

func NewConfig(configFilePath string) Config {
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
