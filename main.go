package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zouipo/yumsday/backend"
)

//go:embed backend/data/migrations
var migrationsFs embed.FS

// @title			Yumsday API
// @version			1.0
// @description		Yumsday is a meal-planning application that includes a menu, a collection of cooking recipes and a grocery list.
// @host 			localhost:8080
// @BasePath 		/api

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	// Generalize the above configuration of the logger to all the project.
	slog.SetDefault(logger)
	defer slog.Debug("Closing app")

	db, err := sql.Open("sqlite3", "yumsday.db")
	if err != nil {
		slog.Error("Failed to open sqlite db", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Flag = CLI parameters when running the program, for example: go run main.go -addr=localhost -port=8080.
	addr := flag.String("addr", "", "Addresses to listen on")
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	migrationsFs, err := fs.Sub(migrationsFs, "backend/data/migrations")
	if err != nil {
		slog.Error("Failed to load migrations filesystem", "error", err)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *addr, *port), // TCP address to listen on, in the form "host:port"
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
		}
	}()
	slog.Info("HTTP server started", "addr", *addr, "port", *port)
	slog.Info(fmt.Sprintf("Swagger docs available at: http://localhost:%d/swagger/index.html", *port))

	<-ctx.Done()
}
