// Package telegram fornece funcionalidades para acessar a API do Telegram.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"gustavonovaes.dev/go-sjc-assist-bot/pkg/config"
	"gustavonovaes.dev/go-sjc-assist-bot/pkg/mongodb"
)

var (
	appConfig config.Config
)

func init() {
	appConfig = config.New()
}

func SetupWebhook() error {
	res, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_TOKEN+"/setWebhook?url="+appConfig.TELEGRAM_WEBHOOK_URL,
		"application/json",
		nil,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyContent := make([]byte, 256)
	res.Body.Read(bodyContent)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to setup webhook, status code: %d, body: %s",
			res.StatusCode,
			bodyContent,
		)
	}

	return nil
}

func HandleWebhook(w http.ResponseWriter, r *http.Request, commands map[string]Command) error {
	w.WriteHeader(http.StatusOK)

	var webhookResponse WebhookResponse
	if err := json.NewDecoder(r.Body).Decode(&webhookResponse); err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}
	r.Body.Close()

	if os.Getenv("DEBUG") == "true" {
		log.Printf("DEBUG: Received message: %+v", webhookResponse.Message)
	}

	beforeHandleCommands(webhookResponse)

	if len(commands) == 0 {
		log.Println("No commands available")
	}

	for command, handler := range commands {
		if strings.Contains(webhookResponse.Message.Text, command) {
			beforeCommandExecution(webhookResponse, command)
			if err := handler(webhookResponse.Message); err != nil {
				log.Printf(
					"ERROR: Failed to execute command %s for user %s/%d in chat %d: %v",
					command,
					webhookResponse.Message.From.Username,
					webhookResponse.Message.From.ID,
					webhookResponse.Message.Chat.ID,
					err,
				)
				w.WriteHeader(http.StatusInternalServerError)
				return err
			}

			afterCommandExecution(webhookResponse, command)
		}
	}

	return nil
}

func beforeHandleCommands(w WebhookResponse) {
	_, err := mongodb.GetCollection("activities").InsertOne(context.Background(), &bson.M{
		"chat_id":   w.Message.Chat.ID,
		"user_id":   w.Message.From.ID,
		"username":  w.Message.From.Username,
		"message":   w.Message.Text,
		"timestamp": time.Now(),
	})
	if err != nil {
		log.Printf("ERROR: Failed to insert activity: %v", err)
	}
}

func beforeCommandExecution(w WebhookResponse, command string) {
	log.Printf(
		"INFO: User %s/%d requested command %s from in chat %d",
		w.Message.From.Username,
		w.Message.From.ID,
		command,
		w.Message.Chat.ID,
	)
}

func afterCommandExecution(w WebhookResponse, command string) {
	//
}

func SendMessage(chatID int, message string) error {
	res, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_TOKEN+"/sendMessage",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{
			"chat_id": %d, 
			"text": "%s",
			"parse_mode": "HTML"
		}`, chatID, message)),
	)

	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d\n%v", res.StatusCode, res.Body)
	}

	return nil
}

func SendPhoto(chatID int, image image.Image) error {
	tempFilepath := fmt.Sprintf(
		"%s/assist-bot-photo_%d_%d.png",
		os.TempDir(),
		chatID,
		time.Now().Unix(),
	)

	f, err := os.Create(tempFilepath)
	if err != nil {
		return fmt.Errorf("failed to create photo file: %v", err)
	}

	if err := png.Encode(f, image); err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)

	writer := multipart.NewWriter(buf)

	writer.CreateFormFile("photo", "photo.png")
	part, err := writer.CreateFormFile("photo", "photo.png")
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, f); err != nil {
		return fmt.Errorf("failed to copy file to part: %v", err)
	}

	body, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_TOKEN+"/sendPhoto?chat_id="+fmt.Sprintf(
			"%d",
			chatID,
		),
		writer.FormDataContentType(),
		buf,
	)
	if err != nil {
		return fmt.Errorf("failed to send photo: %v", err)
	}
	defer body.Body.Close()

	if body.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send photo, status code: %d", body.StatusCode)
	}

	return nil
}
