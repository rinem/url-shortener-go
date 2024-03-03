package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
  "github.com/rinem/url-shortener-go/internal/handlers"
)

func main() {
	r := chi.NewRouter()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	svr := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		sig := <-killSig
		logger.Info("Received kill signal, shutting down!", slog.String("signal", sig.String()))
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 5*time.Second)

		go func() {
			<-shutdownCtx.Done()

			// Print then exit with an error
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Shutdown Deadline Exceeded!")
			}
		}()

		err := svr.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		serverStopCtx()
		logger.Info("Server shutting down!")
		// Cancelling shutdown context
		cancel()
	}()

	go func() {
		err := svr.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

  r.Get("/healthcheck", handlers.NewHealthHandler().ServeHTTP)

	logger.Info("Server Up!")

	<-serverCtx.Done()
}
