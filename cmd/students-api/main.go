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
	"github.com/zaman-hridoy/go-student-api/internal/http/handlers/student"
	"github.com/zaman-hridoy/go-student-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetStudentList(storage))


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
	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
	
}