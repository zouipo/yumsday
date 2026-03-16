package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zouipo/yumsday/backend"
	"github.com/zouipo/yumsday/internal/config"
)

//go:embed backend/data/migrations
var migrationsFs embed.FS

// @title			Yumsday API
// @version			1.0
// @description		Yumsday is a meal-planning application that includes a menu, a collection of cooking recipes and a grocery list.
// @host 			localhost:8080
// @BasePath 		/

var cmd = &cobra.Command{
	Use:   "yumsday",
	Short: "yumsday",
	Run:   run,
}

func init() {
	// Define cli flags
	cmd.PersistentFlags().String("host", "localhost", "Server host")
	cmd.PersistentFlags().Int("port", 8080, "Server port")
	cmd.PersistentFlags().String("db-path", "yumsday.db", "Path to the sqlite database")
	cmd.PersistentFlags().String("log-level", "info", "Log level")

	// Bind cli flags to viper values
	viper.BindPFlag("host", cmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", cmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("db_path", cmd.PersistentFlags().Lookup("db-path"))
	viper.BindPFlag("log_level", cmd.PersistentFlags().Lookup("log-level"))
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error(
			"Failed to load configuration",
			"error", err,
		)
		return
	}

	level := slog.LevelWarn
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
	// Generalize the above configuration of the logger to all the project.
	slog.SetDefault(logger)
	defer slog.Debug("Closing app")

	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		slog.Error("Failed to open sqlite db", "error", err)
		return
	}
	slog.Info("Opened db", "db_path", cfg.DBPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	migrationsFs, err := fs.Sub(migrationsFs, "backend/data/migrations")
	if err != nil {
		slog.Error("Failed to load migrations filesystem", "error", err)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // TCP address to listen on, in the form "host:port"
		Handler: backend.NewAPIServer(db, migrationsFs),
	}

	// Goroutine waiting for a signal from the OS to shut "gracefully" the server and its working goroutines.
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // SIGINT = Ctrl+C, SIGTERM = kill command.
		<-sigCh
		signal.Stop(sigCh)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
		defer shutdownCancel()

		// ShutDown() shuts the server "gracefully" without interrupting any active connections.
		// It waits indefinitely for connections to go idle, until the given context expires (1 second).
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server gracefully", "error", err)
		} else {
			slog.Info("Server stopped succesfully")
		}

		cancel()
	}()

	// Goroutine to start the server and wait for the server to be shut down.
	go func() {
		// ListenAndServe() blocks until Server.Shutdown or Server.Close is called,
		// then it returns the returned error ErrServerClosed.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server stopped listening", "error", err)
			cancel()
		}
	}()
	slog.Info("HTTP server started", "addr", cfg.Host, "port", cfg.Port)
	slog.Info(fmt.Sprintf("Swagger docs available at: http://%s:%d/swagger/index.html", cfg.Host, cfg.Port))

	<-ctx.Done()
}

func main() {
	cmd.Execute()
}
