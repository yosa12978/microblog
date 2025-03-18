package config

import (
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name   string `yaml:"name"`
		Desc   string `yaml:"description"`
		Bottom string `yaml:"bottom"`
	} `yaml:"app"`
	Server struct {
		Addr     string `yaml:"addr"`
		Hostname string `yaml:"hostname"`
	} `yaml:"server"`
	Postgres struct {
		Addr     string `yaml:"addr"`
		DB       string `yaml:"db"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"postgres"`
}

var (
	initConfig sync.Once
	conf       Config
)

func Get() Config {
	initConfig.Do(func() {
		file, err := os.Open("config.yml")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		yaml.NewDecoder(file).Decode(&conf)
		envconfig.MustProcess("microblog", &conf)
	})
	return conf
}
