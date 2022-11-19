package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jtieri/habbgo/game/item"
	"github.com/jtieri/habbgo/game/jobs"
	"github.com/jtieri/habbgo/game/navigator"
	"github.com/jtieri/habbgo/game/player"
	"github.com/jtieri/habbgo/game/room"
	"github.com/jtieri/habbgo/game/scheduler"
	"github.com/jtieri/habbgo/game/service"
	"github.com/jtieri/habbgo/internal"
	"github.com/jtieri/habbgo/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// startCmd starts the game server.
func startCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"st"},
		Short:   "Start the habbgo game server.",
		Args:    cobra.ExactArgs(0),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s start
$ %s st`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			a.Log.Info("Booting up habbgo")

			if a.Debug {
				a.Config.Global.Debug = a.Debug
			}

			// Start the debug server
			debugAddr, err := cmd.Flags().GetString(flagDebugServer)
			if err != nil {
				return err
			}

			if debugAddr != "" {
				listener, err := net.Listen("tcp", debugAddr)
				if err != nil {
					return fmt.Errorf("failed to start debug server on address %s: %w", debugAddr, err)
				}
				log := a.Log.With(zap.String("sys", "debugserver"))
				log.Info("Debug server listening", zap.String("addr", debugAddr))
				internal.StartDebugServer(cmd.Context(), log, listener)
			}

			// Connect to database
			a.Log.Info("Attempting to make connection with the database",
				zap.String("host", a.Config.DB.Host),
				zap.Int("port", a.Config.DB.Port),
				zap.String("username", a.Config.DB.Username),
				zap.String("db_name", a.Config.DB.Name),
				zap.String("db_driver", a.Config.DB.Driver),
				zap.String("ssl_mode", a.Config.DB.SSLMode),
			)

			db, err := connectToDatabase(a.Config)
			if err != nil {
				panic(err)
			}

			a.Log.Info("Successfully connected to database",
				zap.String("host", a.Config.DB.Host),
				zap.Int("port", a.Config.DB.Port),
				zap.String("username", a.Config.DB.Username),
				zap.String("db_name", a.Config.DB.Name),
				zap.String("db_driver", a.Config.DB.Driver),
				zap.String("ssl_mode", a.Config.DB.SSLMode),
			)

			ctx, cancel := context.WithCancel(cmd.Context())
			sch := scheduler.NewGameScheduler(ctx, cancel, a.Log.With())

			// Build the game services
			a.Log.Info("Preparing game services")
			services := buildGameServices(ctx, a.Log.With(), db, sch, cancel)

			// Build game server
			gameServer := server.New(
				a.Log.With(),
				db,
				a.Config.Server.Host,
				a.Config.Server.Port,
				a.Config.Server.MaxConnsPerPlayer,
				a.Config.Global.Debug,
				services,
			)

			a.Log.Info("Starting the game server",
				zap.String("host", a.Config.Server.Host),
				zap.Int("port", a.Config.Server.Port),
				zap.Int("max_conns_per_player", a.Config.Server.MaxConnsPerPlayer),
				zap.Bool("debug", a.Config.Global.Debug),
			)

			// Start the game server
			errorChan := gameServer.Start(cmd.Context())
			defer gameServer.Stop()

			// Start the game scheduler
			a.Log.Info("Starting the game scheduler")
			go sch.Start()
			defer sch.Stop()

			// TODO make all job durations configurable from config file
			// Start the main game job which checks things like rewards, player count, etc.
			sch.ScheduleJob(jobs.NewGameJob(ctx, cancel, 1*time.Second))

			// Block until an error comes across the error channel
			if err := <-errorChan; err != nil && !errors.Is(err, context.Canceled) {
				a.Log.Warn("Game server failed to start",
					zap.Error(err),
				)
				return err
			}

			return nil
		},
	}
	return debugServerFlag(a.Viper, cmd)
}

func connectToDatabase(c *Config) (*sql.DB, error) {
	// Open the connection to the database.
	db, err := sql.Open(c.DB.Driver, c.DB.ConnectionString())
	if err != nil {
		return nil, err
	}

	// Check that the connection to the database is alive.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// buildGameServices initializes the game Services when starting the server.
// The server maintains a reference to these services and as new sessions are accepted,
// the players are passed a reference to the Services as well. This allows players
// to have access to the global state of the game.
func buildGameServices(
	ctx context.Context,
	log *zap.Logger,
	db *sql.DB,
	scheduler *scheduler.GameScheduler,
	cancel context.CancelFunc,
) *service.Proxies {
	ns := navigator.NewNavigatorService(
		ctx,
		log.With(zap.String("service_name", "navigator_service")),
		db,
		scheduler,
		cancel,
	)
	go ns.Start()

	rs := room.NewRoomService(
		ctx,
		log.With(zap.String("service_name", "room_service")),
		db,
		scheduler,
		cancel,
	)
	go rs.Start()

	is := item.NewItemService(
		ctx,
		log.With(zap.String("service_name", "item_service")),
		db,
		scheduler,
		cancel,
	)
	go is.Start()

	ps := player.NewPlayerService(ctx,
		log.With(zap.String("service_name", "player_service")),
		db,
		scheduler,
		cancel,
	)
	go ps.Start()

	return &service.Proxies{
		Rooms:     room.NewProxy(rs.Channels()),
		Items:     item.NewProxy(is.Channels()),
		Navigator: navigator.NewProxy(ns.Channels()),
		Players:   player.NewProxy(ps.Channels()),
	}
}
