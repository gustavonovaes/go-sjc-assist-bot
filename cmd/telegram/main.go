package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/sspsp"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
	"gustavonovaes.dev/go-sjc-assist-bot/pkg/mongodb"
)

const (
	URL_PATH              = "/"
	ADDR                  = ":443"
	TIMEOUT_SETUP_WEBHOOK = 10 * time.Second
)

var COMMANDS = map[string]telegram.Command{
	"/start": telegram.CommandStart,
	"/ajuda": telegram.CommandStart,
	"/sobre": telegram.CommandAbout,

	// cetesb
	"/qualidadeAr": cetesb.CommandQualidadeAr,

	// sspsp
	"/crimes": sspsp.CommandCrimes,
}

func main() {
	server := http.NewServeMux()
	server.HandleFunc(URL_PATH, func(w http.ResponseWriter, r *http.Request) {
		telegram.HandleWebhook(w, r, COMMANDS)
	})

	// Calls Telegram API to setup webhook after a timeout
	go func() {
		<-time.After(TIMEOUT_SETUP_WEBHOOK)
		err := telegram.SetupWebhook()
		if err != nil {
			log.Fatalf("ERROR: Fail to setup webhook: %v", err)
		}

		log.Println("INFO: Webhook setup successfully")
	}()

	// Start the server and listen for shutdown signals to ensure graceful termination of the server
	listenWithGracefulShutdown(ADDR, server)
}

func listenWithGracefulShutdown(addr string, server http.Handler) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("INFO: Listening in %s...", addr)
		err := http.ListenAndServe(addr, server)
		if err != nil {
			log.Fatalf("ERROR: Fail to start on addr: %q", addr)
		}
	}()

	<-stopChan
	log.Println("INFO: Shutting down server...")

	log.Println("INFO: Disconnecting DB...")
	if err := mongodb.Close(); err != nil {
		panic(err)
	}

	log.Println("INFO: Server gracefully stopped")
}
