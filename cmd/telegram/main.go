package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/infra/mongodb"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/infra/telegram"
)

const (
	WEBHOOK_URL_PATH    = "/"
	WEBHOOK_SETUP_DELAY = 10 * time.Second
	LISTEN_ADDR         = ":443"
)

var commands = map[string]telegram.Command{
	// common
	"/start": CommandStart,
	"/ajuda": CommandStart,
	"/sobre": CommandAbout,

	// cetesb
	"/qualidadeAr": CommandQualityAir,

	// sspsp
	"/crimes":     CommandCrimes,
	"/mapaCrimes": CommandMapCrimes,
}

func main() {
	server := telegram.NewWebhookServer(commands, logUserActivityMiddleware)

	// Calls Telegram API to setup webhook after a delay to ensure the server is ready
	go func() {
		<-time.After(WEBHOOK_SETUP_DELAY)
		err := telegram.SetupWebhook()
		if err != nil {
			log.Fatalf("ERROR: Fail to setup webhook: %v", err)
		}

		log.Println("INFO: Webhook setup successfully")
	}()

	// Start the server and listen for shutdown signals to ensure graceful termination of the server
	listenWithGracefulShutdown(LISTEN_ADDR, server)
}

func logUserActivityMiddleware(wr telegram.WebhookResponse) telegram.WebhookResponse {
	activity := &bson.M{
		"chat_id":   wr.Message.Chat.ID,
		"user_id":   wr.Message.From.ID,
		"username":  wr.Message.From.Username,
		"message":   wr.Message.Text,
		"timestamp": time.Now(),
	}

	// Insert the activity into MongoDB asynchronously
	go func() {
		if err := mongodb.SaveCollection("activities", activity); err != nil {
			log.Printf("ERROR: Fail to insert activity: %v\n%+v", err, activity)
		}
	}()

	return wr
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
