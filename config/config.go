package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var defaultHome = os.ExpandEnv("$HOME/.habbgo")

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

func LoadConfig() (*Config, error) {
	c := &Config{}
	cfgPath := path.Join(defaultHome, "config", "config.yaml")
	data, err := ioutil.ReadFile(cfgPath)

	if err != nil {
		log.Printf("Failed to read config file %s \n", cfgPath)
		log.Printf("Creating default config file %s \n", cfgPath)
		c, err = InitConfig()
		if err != nil {
			return nil, err
		}

		log.Println("Successfully created the default config file")
	} else {
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal the config file at %s. Err: %w", cfgPath, err)
		}

		log.Println("The file 'config.yml' has been successfully loaded.")
	}

	return c, nil
}

func InitConfig() (*Config, error) {
	var c *Config
	cfgDir := path.Join(defaultHome, "config")
	cfgPath := path.Join(cfgDir, "config.yaml")

	// If the config doesn't exist...
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		// And the config folder doesn't exist...
		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			// And the home folder doesn't exist
			if _, err := os.Stat(defaultHome); os.IsNotExist(err) {
				// Create the home folder
				if err = os.Mkdir(defaultHome, os.ModePerm); err != nil {
					return nil, err
				}
			}

			// Create the home config folder
			if err = os.Mkdir(cfgDir, os.ModePerm); err != nil {
				return nil, err
			}
		}

		// Then create the file...
		f, err := os.Create(cfgPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		c = defaultConfig()
		bz, err := yaml.Marshal(c)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal the config file. Err: %w \n", err)
		}

		// And write the default config to that location...
		if _, err = f.Write(bz); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func defaultConfig() *Config {
	server := &ServerCfg{
		Host:     "127.0.0.1",
		Port:     11235,
		MaxConns: 2,
		Debug:    true,
	}
	db := &DatabaseCfg{
		User:     "root",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     3306,
		Name:     "habbgo",
	}
	return &Config{
		Server: server,
		DB:     db,
	}
}
