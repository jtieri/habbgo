package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	zaplogfmt "github.com/jsternberg/zap-logfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
)

var defaultHome = filepath.Join(os.Getenv("HOME"), ".habbgo")

const appName = "habbgo"

// NewRootCmd returns the root command for habbgo.
// If log is nil, a new zap.Logger is set on the app state
// based on the command line flags regarding logging.
func NewRootCmd(log *zap.Logger) *cobra.Command {
	// Use a local app state instance scoped to the new root command,
	// so that tests don't concurrently access the state.
	a := &appState{
		Viper: viper.New(),

		Log: log,
	}

	// RootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "habbgo - habbo hotel emulation.",
		Long: strings.TrimSpace(`
   habbgo is a server application that emulates Shockwave era Habbo Hotel.`),
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// Inside persistent pre-run because this takes effect after flags are parsed.
		if log == nil {
			log, err := newRootLogger(a.Viper.GetString("log-format"), a.Viper.GetBool("debug"))
			if err != nil {
				return err
			}

			a.Log = log
		}

		// reads `homeDir/config/config.yaml` into `a.Config`
		return initConfig(rootCmd, a)
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, _ []string) {
		// Force syncing the logs before exit, if anything is buffered.
		a.Log.Sync()
	}

	// Register --home flag
	rootCmd.PersistentFlags().StringVar(&a.HomePath, flagHome, defaultHome, "set home directory")
	if err := a.Viper.BindPFlag(flagHome, rootCmd.PersistentFlags().Lookup(flagHome)); err != nil {
		panic(err)
	}

	// Register --debug flag
	rootCmd.PersistentFlags().BoolVarP(&a.Debug, "debug", "d", false, "enable debug output")
	if err := a.Viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().String("log-format", "auto", "log output format (auto, logfmt, json, or console)")
	if err := a.Viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		panic(err)
	}

	// Register subcommands
	rootCmd.AddCommand(
		configCmd(a),
		startCmd(a),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd(nil)
	rootCmd.SilenceUsage = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt) // Using signal.Notify, instead of signal.NotifyContext, in order to see details of signal.
	go func() {
		// Wait for interrupt signal.
		sig := <-sigCh

		// Stop context on root command.
		// If the invoked command respects this quickly, the main goroutine will quit right away.
		cancel()

		// Short delay before printing the received signal message.
		// This should result in cleaner output from non-interactive commands that stop quickly.
		time.Sleep(250 * time.Millisecond)
		fmt.Fprintf(os.Stderr, "Received signal %v. Attempting clean shutdown. Send interrupt again to force hard shutdown.\n", sig)

		// Dump all goroutines on panic, not just the current one.
		debug.SetTraceback("all")

		// Block waiting for a second interrupt or a timeout.
		// The main goroutine ought to finish before either case is reached.
		// But if a case is reached, panic so that we get a non-zero exit and a dump of remaining goroutines.
		select {
		case <-time.After(time.Minute):
			panic(errors.New("habbgo did not shut down within one minute of interrupt"))
		case sig := <-sigCh:
			panic(fmt.Errorf("received signal %v; forcing quit", sig))
		}
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func newRootLogger(format string, debug bool) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format("2006-01-02T15:04:05.000000Z07:00"))
	}
	config.LevelKey = "lvl"

	var enc zapcore.Encoder
	switch format {
	case "json":
		enc = zapcore.NewJSONEncoder(config)
	case "console":
		enc = zapcore.NewConsoleEncoder(config)
	case "logfmt":
		enc = zaplogfmt.NewEncoder(config)
	case "auto":
		if term.IsTerminal(int(os.Stderr.Fd())) {
			// When a user runs habbgo in the foreground, use easier to read output.
			enc = zapcore.NewConsoleEncoder(config)
		} else {
			// Otherwise, use consistent logfmt format for simplistic machine processing.
			enc = zaplogfmt.NewEncoder(config)
		}
	default:
		return nil, fmt.Errorf("unrecognized log format %q", format)
	}

	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}
	return zap.New(zapcore.NewCore(
		enc,
		os.Stderr,
		level,
	)), nil
}
