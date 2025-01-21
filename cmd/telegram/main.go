package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/cetesb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
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
	"/cetesb": cetesb.CommandQualidadeAr,
}

func main() {
	server := http.NewServeMux()
	server.HandleFunc(URL_PATH, func(w http.ResponseWriter, r *http.Request) {
		telegram.HandleWebhook(w, r, COMMANDS)
	})

	go func() {
		<-time.After(TIMEOUT_SETUP_WEBHOOK)
		err := telegram.SetupWebhook()
		if err != nil {
			log.Fatalf("ERROR: Fail to setup webhook: %v", err)
		}

		log.Println("INFO: Webhook setup successfully")
	}()

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

	// Add your cleanup code here

	log.Println("INFO: Server gracefully stopped")
}
