package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Server struct {
		Host     string `yaml:"host"`
		Port     int16  `yaml:"port"`
		MaxConns int    `yaml:"maxconns"`
	}
	Database struct {
	}
}

func LoadConfig() Config {
	data, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal("Failed to read config file 'config.yml', check that it exists & try again.", err)
	}

	config := Config{}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to unmarshal config file 'config.yml', check that its format is correct & try again.", err)
	}

	log.Println("The file 'config.yml' has been successfully loaded.")
	return config
}
