package v1

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pmoieni/auth/log"
)

type Server struct {
	Handler *echo.Echo
}

func (s *Server) Run(addr string) {
	s.init()
	server := http.Server{
		Addr:         addr,
		Handler:      s.Handler,
		ErrorLog:     log.StdLogger,     // set the logger for the server
		ReadTimeout:  10 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// implement graceful shutdown for unexpected signals
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Logger.Fatal("graceful shutdown timed out. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Logger.Fatal(err.Error())
		}
		serverStopCtx()

		defer cancel()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Logger.Fatal(err.Error())
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func (s *Server) init() {
	s.Handler = echo.New()
	s.Handler.HTTPErrorHandler = customHTTPErrorHandler
	s.initRoutes()
}
