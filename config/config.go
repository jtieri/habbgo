package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

var defaultSubPath = path.Join("habbgo", "config", "config.yaml")

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
	DBDriver          string `yaml:"db-driver"`
	DBSSLMode         string `yaml:"db-ssl-mode"`
}

// LoadConfig will attempt to load the config file from the system suitable location for config files, as per
// XDG Base Directory Specification, and upon failure will initialize the default config file in this same location
func LoadConfig() (*Config, error) {
	c := &Config{}
	defaultCfgPath := path.Join(xdg.ConfigHome, defaultSubPath)

	// search for the config file & if it exists read it or else initialize default config file
	configFilePath, err := xdg.SearchConfigFile(defaultSubPath)
	if err != nil {
		log.Printf("Failed to read config file %s \n", defaultCfgPath)
		log.Printf("Creating default config file %s... \n", defaultCfgPath)

		c, err = InitConfig()
		if err != nil {
			return nil, err
		}

		log.Printf("Successfully created the default config file %s \n", defaultCfgPath)
	} else {
		data, err := ioutil.ReadFile(configFilePath)
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return nil, err
		}
		log.Printf("The file %s has been successfully loaded \n", configFilePath)
	}

	return c, nil
}

// InitConfig will attempt to initialize the default config file in a system suitable location for config files, as per
// XDG Base Directory Specification.
func InitConfig() (*Config, error) {
	var c *Config

	// retrieve a suitable location for the specified config file
	configFilePath, err := xdg.ConfigFile(defaultSubPath)
	if err != nil {
		return nil, err
	}

	// create the file on disk
	f, err := os.Create(configFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c = defaultConfig()
	bz, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}

	// write the default config to that location
	if _, err = f.Write(bz); err != nil {
		return nil, err
	}
	return c, nil
}

// defaultConfig returns a default configuration struct to be marshaled to disk
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
		DBDriver:          "postgres",
		DBSSLMode:         "disable",
	}
}
