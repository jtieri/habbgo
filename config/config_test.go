package config

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	home, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	cfg := path.Join(home, "config.yaml")
	LoadConfig(cfg)
	_, err = os.Stat(cfg)
	require.NoError(t, err)
}

func TestCreatesDefaultConfig(t *testing.T) {
	home, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	t.Log("Initializing default config file... ")
	c := InitDefaultConfig()
	bz, err := yaml.Marshal(c)
	require.NoError(t, err)

	cfg := path.Join(home, "config.yaml")
	t.Logf("Writing default config file to %s", cfg)
	require.NoError(t, os.WriteFile(cfg, bz, 0644))

	t.Logf("Checking that config.yaml exists at %s... ", cfg)
	_, err = os.Stat(cfg)
	require.NoError(t, err)
}
