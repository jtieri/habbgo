package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func configCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage the application's configuration file",
	}

	cmd.AddCommand(
		configShowCmd(a),
		configInitCmd(),
	)

	return cmd
}

type Config struct {
	Global *GlobalConfig
	Server *ServerConfig
	DB     *DatabaseConfig
}

type GlobalConfig struct {
	Debug bool `yaml:"debug" json:"debug"`
}

type ServerConfig struct {
	Host              string `yaml:"host" json:"host"`
	Port              int    `yaml:"port" json:"port"`
	MaxConnsPerPlayer int    `yaml:"player-max-conn" json:"player-max-conn"`
}

type DatabaseConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Name     string `yaml:"db-name" json:"db-name"`
	Driver   string `yaml:"driver" json:"driver"`
	SSLMode  string `yaml:"ssl-mode" json:"ssl-mode"`
}

// configInitCmd initializes an empty config at the location specified via the --home flag.
func configInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Creates a default home directory at path defined by --home",
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config init --home %s
$ %s cfg i`, appName, defaultHome, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}

			cfgDir := path.Join(home, "config")
			cfgPath := path.Join(cfgDir, "config.yaml")

			if err = createConfig(home); err == nil {
				return nil
			}

			// Otherwise, the config file exists, and an error is returned...
			return fmt.Errorf("config already exists: %s", cfgPath)
		},
	}
	return cmd
}

// configShowCmd returns the configuration file in json or yaml format.
func configShowCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s", "list", "l"},
		Short:   "Prints current configuration",
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config show --home %s
$ %s cfg list`, appName, defaultHome, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}

			cfgDir := path.Join(home, "config")
			cfgPath := path.Join(cfgDir, "config.yaml")

			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				if _, err := os.Stat(home); os.IsNotExist(err) {
					return fmt.Errorf("home path does not exist: %s", home)
				}
				return fmt.Errorf("config does not exist: %s", cfgPath)
			}

			jsn, err := cmd.Flags().GetBool(flagJSON)
			if err != nil {
				return err
			}
			yml, err := cmd.Flags().GetBool(flagYAML)
			if err != nil {
				return err
			}
			switch {
			case yml && jsn:
				return fmt.Errorf("can't pass both --json and --yaml, must pick one")
			case jsn:
				out, err := json.Marshal(a.Config)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
				return nil
			default:
				out, err := yaml.Marshal(a.Config)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
				return nil
			}
		},
	}

	return yamlFlag(a.Viper, jsonFlag(a.Viper, cmd))
}

// createConfig writes the default config file to disk in the location specified by home.
func createConfig(home string) error {
	cfgDir := path.Join(home, "config")
	cfgPath := path.Join(cfgDir, "config.yaml")

	// If the config doesn't exist...
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		// And the config folder doesn't exist...
		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			// And the home folder doesn't exist
			if _, err := os.Stat(home); os.IsNotExist(err) {
				// Create the home folder
				if err = os.Mkdir(home, os.ModePerm); err != nil {
					return err
				}
			}
			// Create the home config folder
			if err = os.Mkdir(cfgDir, os.ModePerm); err != nil {
				return err
			}
		}

		// Then create the file...
		f, err := os.Create(cfgPath)
		if err != nil {
			return err
		}
		defer f.Close()

		// And write the default config to that location...
		cfg, err := yaml.Marshal(defaultConfig())
		if err != nil {
			return err
		}

		if _, err = f.Write(cfg); err != nil {
			return err
		}
	}
	return nil
}

// initConfig reads in the config file and ENV variables if set.
// This is called as a persistent pre-run command of the root command.
func initConfig(cmd *cobra.Command, a *appState) error {
	home, err := cmd.PersistentFlags().GetString(flagHome)
	if err != nil {
		return err
	}

	cfgPath := path.Join(home, "config", "config.yaml")
	if _, err = os.Stat(cfgPath); err == nil {
		a.Viper.SetConfigFile(cfgPath)
		err = a.Viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("failed to read in config: %w", err)
		}

		// read the config file bytes
		file, err := os.ReadFile(a.Viper.ConfigFileUsed())
		if err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}

		// unmarshall them into the struct
		if err = yaml.Unmarshal(file, &a.Config); err != nil {
			return fmt.Errorf("error unmarshalling config: %w", err)
		}
	}

	return nil
}

func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		d.Username, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}

// defaultConfig returns a default configuration struct to be marshaled to disk
func defaultConfig() *Config {
	return &Config{
		Global: &GlobalConfig{
			Debug: false,
		},
		Server: &ServerConfig{
			Host:              "127.0.0.1",
			Port:              11235,
			MaxConnsPerPlayer: 2,
		},
		DB: &DatabaseConfig{
			Username: "anon",
			Password: "password",
			Host:     "127.0.0.1",
			Port:     5432,
			Name:     "habbgo",
			Driver:   "postgres",
			SSLMode:  "disable",
		},
	}
}
