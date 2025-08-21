package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zaman-hridoy/go-student-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to students api updaed"))
	})


	// setup server
	server := &http.Server {
		Addr: cfg.Address,
		Handler: router,
	}

	



	slog.Info("server started", slog.String("address", cfg.Address))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func ()  {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("Error", err)
			log.Fatal("Fialed to start server")
		}
	}()

	<- done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
	
}