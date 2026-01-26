package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	/*** SET LOGGER ***/
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	// Generalize this configuration of the logger for all the project.
	// After this line, every call to slog.X call logger.
	slog.SetDefault(logger)
	defer slog.Debug("Closing app")

	/*** DATABASE INITIALIZATION ***/
	db, err := sql.Open("sqlite3", "yumsday.db")
	if err != nil {
		slog.Error("Failed to open sqlite db", "error", err)
		return
	}

	/*** CONTEXT CREATION ***/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*** SERVER INITIALIZATION ***/
	// Flag = CLI parameters;
	// When running the program, you can give it new parameters for this flags.
	addr := flag.String("addr", "", "Addresses to listen on")
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", *addr, *port),
	}

	// Goroutine waiting for a signal from the OS to shut "gracefully" the server and its working goroutines.
	go func() {
		// Create a channel that listens to os.Signal and reacts to interruption or killing of the program.
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		signal.Stop(sigCh)

		// Create a context for 1s
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
		defer shutdownCancel()

		// ShutDown() shuts the server "graefully" without interrupting any active connections.
		// It waits indefinitely for connections to go idle, so it can never returns, unless the given context expire.
		// Our context expire in 1s.
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server gracefully", "error", err)
		} else {
			slog.Info("Server stopped succesfully")
		}

		cancel()
	}()

	// Goroutine to start the server and wait for the server to be shut down.
	go func() {
		// ListenAndServe() blocks until Server.Shutdown or Server.Close,
		// then it returns the returned error is ErrServerClosed.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server stopped listening", "error", err)
		}
	}()
	slog.Info("HTTP server started", "addr", *addr, "port", *port)

	<-ctx.Done()
}
