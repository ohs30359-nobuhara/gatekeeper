package util

import (
	"github.com/go-yaml/yaml"
	"log"
	"os"
)

type Config struct {
	Server server `yaml:"server"`
	Proxy  proxy  `yaml:"proxy"`
	Redis  redis  `yaml:"redis"`
}

type server struct {
	Port string `yaml:"port"`
}

type proxy struct {
	RateLimit int `yaml:"rateLimit"`
	Timeout string `yaml:"timeout"`
	Target  struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Schema string `yaml:"schema"`
	}
}

type redis struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	No int `yaml:"no"`
	Password string `yaml:"password"`
	TimeoutMs struct  {
		Read int `yaml:"read"`
		Write int `yaml:"write"`
	}
}

func ConfigLoad() *Config {
	f, e := os.Open("config.yaml")

	if e != nil {
		log.Fatal(e.Error())
		return nil
	}

	defer f.Close()

	var c Config

	if e := yaml.NewDecoder(f).Decode(&c); e != nil {
		log.Fatal(e.Error())
		return nil
	}

	return &c
}