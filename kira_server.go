package kira

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lafriakh/kira/modules/config"
)

// StartServer - Start kira server
func (app *App) StartServer(server *http.Server) {
	// Gracefully shutdown
	go app.GracefullyShutdown(server)

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		app.logger.Panicf("%v", err)
	}

	// On start
	if app.lifecyles.onStart != nil {
		go app.lifecyles.onStart()
	}

	if !app.Configs.GetBool("server.tls", false) {
		app.logger.Infof("Starting HTTP server, Listening at %s", "http://"+server.Addr)

		if err := server.Serve(listener); err != http.ErrServerClosed {
			app.logger.Panicf("%v", err)
		}
	} else {
		app.logger.Infof("Starting HTTPS server, Listening at %s", "https://"+server.Addr)
		// To generate keys:
		//   - openssl genrsa -out server.key 2048
		//   - openssl ecparam -genkey -name secp384r1 -out server.key
		//   - openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
		certificateFile := app.Configs.GetString("server.tls_certificate", "./server.crt")
		keyFile := app.Configs.GetString("server.tls_key", "./server.key")
		if err := server.ServeTLS(listener, certificateFile, keyFile); err != http.ErrServerClosed {
			app.logger.Panicf("%v", err)
		}
	}
	app.logger.Info("Server closed")
}

// GracefullyShutdown the server
func (app *App) GracefullyShutdown(server *http.Server) {
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigquit
	app.logger.Infof("Signal to shutdown the server: %+v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.logger.Fatalf("Unable to shutdown server: %v", err)
	}

	// Shutdown callback
	if app.lifecyles.onShutdown != nil {
		app.lifecyles.onShutdown()
	}
}

func getServerAddr(config *config.Config) string {
	// Server HOST/PORT
	host := config.GetString("server.host", "127.0.0.1")
	port := config.GetInt("server.port", 8080)

	// Server Addr
	return fmt.Sprintf("%s:%d", host, port)
}
