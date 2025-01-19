package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram/cetesb"
)

const (
	URL_PATH = "/"
	ADDR     = ":443"
)

var COMMANDS = map[string]telegram.Command{
	"/start": telegram.CommandStart,
	"/ajuda": telegram.CommandStart,
	"/sobre": telegram.CommandAbout,

	// cetesb
	"/cetesb:qualidade-ar": cetesb.CommandQualidadeAr,
}

func main() {
	err := telegram.SetupWebhook()
	if err != nil {
		log.Fatalf("ERROR: Fail to setup webhook: %v", err)
	}

	server := http.NewServeMux()
	server.HandleFunc(URL_PATH, func(w http.ResponseWriter, r *http.Request) {
		telegram.HandleWebhook(w, r, COMMANDS)
	})

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
