package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zouipo/yumsday/backend"
	"github.com/zouipo/yumsday/internal/config"
	"github.com/zouipo/yumsday/internal/dbutils"
)

//go:embed backend/data/migrations
var migrationsFs embed.FS

// @title			Yumsday API
// @version			1.0
// @description		Yumsday is a meal-planning application that includes a menu, a collection of cooking recipes and a grocery list.
// @host 			localhost:8080
// @BasePath 		/

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error(
			"Failed to load configuration",
			"error", err,
		)
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	// Generalize the above configuration of the logger to all the project.
	slog.SetDefault(logger)
	defer slog.Debug("Closing app")

	db, err := dbutils.OpenDb(cfg.DBPath)
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

	// WaitGroup used to synchronize tasks running in dedicated goroutines
	// like the persistence of sessions in the db.
	var tasksWG sync.WaitGroup

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // TCP address to listen on, in the form "host:port"
		Handler: backend.NewAPIServer(db, migrationsFs, &tasksWG),
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
	tasksWG.Wait()
}
