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
	ServerHost        string `yaml:"server-host"`
	ServerPort        int    `yaml:"server-port"`
	MaxConnsPerPlayer int    `yaml:"player-maxconns"`
	Debug             bool   `yaml:"debug"`
	DBUser            string `yaml:"db-user"`
	DBPassword        string `yaml:"db-password"`
	DBHost            string `yaml:"db-host"`
	DBPort            int16  `yaml:"db-port"`
	DBName            string `yaml:"db-name"`
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
	return &Config{
		ServerHost:        "127.0.0.1",
		ServerPort:        11235,
		MaxConnsPerPlayer: 2,
		Debug:             true,
		DBUser:            "postgres",
		DBPassword:        "password",
		DBHost:            "127.0.0.1",
		DBPort:            5432,
		DBName:            "habbgo",
	}
}
