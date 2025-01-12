package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/telegram"
)

const URL_PATH = "/telegram"
const ADDR = ":8080"

func main() {
	err := telegram.SetupWebhook()
	if err != nil {
		log.Fatalf("ERROR: Fail to setup webhook: %v", err)
	}

	server := http.NewServeMux()
	server.HandleFunc(URL_PATH, telegram.HandleWebhook)

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
