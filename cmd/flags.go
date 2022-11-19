package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagHome        = "home"
	flagJSON        = "json"
	flagYAML        = "yaml"
	flagDebugServer = "debug-server"
)

const (
	defaultJSON        = false
	defaultYAML        = false
	defaultDebugServer = "127.0.0.1:49152"
)

func yamlFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP(flagYAML, "y", defaultYAML, "returns the response in yaml format")
	if err := v.BindPFlag(flagYAML, cmd.Flags().Lookup(flagYAML)); err != nil {
		panic(err)
	}
	return cmd
}

func jsonFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP(flagJSON, "j", defaultJSON, "returns the response in json format")
	if err := v.BindPFlag(flagJSON, cmd.Flags().Lookup(flagJSON)); err != nil {
		panic(err)
	}
	return cmd
}

func debugServerFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(flagDebugServer, defaultDebugServer, "host address for debug server. Empty string will disable debug server.")
	if err := v.BindPFlag(flagDebugServer, cmd.Flags().Lookup(flagDebugServer)); err != nil {
		panic(err)
	}
	return cmd
}
