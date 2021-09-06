package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Server *ServerCfg
	DB     *DatabaseCfg
}

type ServerCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	MaxConns int    `yaml:"maxconns"`
	Debug    bool   `yaml:"debug"`
}

type DatabaseCfg struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int16  `yaml:"port"`
	Name     string `yaml:"name"`
}

func LoadConfig(file string) *Config {
	c := &Config{}
	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println("Failed to read config file 'config.yml'")
		log.Println("Creating default config file...")
		c = InitDefaultConfig()

		bz, err := yaml.Marshal(c)
		if err != nil {
			log.Fatal("An error occured while writing the default config file...")
		}

		err = os.WriteFile(file, bz, 0644)
		if err != nil {
			log.Fatal("An error occured while writing the default config file...")
		}

		log.Println("Successfully created the default config file")
	} else {
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			log.Fatal("Failed to unmarshal config file 'config.yml', check that its format is correct & try again.", err)
		}

		log.Println("The file 'config.yml' has been successfully loaded.")
	}

	return c
}

func InitDefaultConfig() *Config {
	server := &ServerCfg{
		Host:     "127.0.0.1",
		Port:     11235,
		MaxConns: 2,
		Debug:    true,
	}
	db := &DatabaseCfg{
		User:     "",
		Password: "",
		Host:     "",
		Port:     0,
		Name:     "",
	}
	return &Config{
		Server: server,
		DB:     db,
	}
}
