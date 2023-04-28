package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GRPC     int    `yaml:"grpc"`
	HTTP     int    `yaml:"http"`
	LogLevel string `yaml:"log_level"`
	TLS      struct {
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
		UseTlS   bool   `yaml:"use_tls"`
	} `yaml:"tls"`
	Cluster []string `yaml:"cluster"`
	Dir     string   `yaml:"dir"`
}

func NewConfig(filename string) *Config {
	c := &Config{}
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
