// Package telegram fornece funcionalidades para acessar a API do Telegram.
package telegram

import (
	"bytes"
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

	"gustavonovaes.dev/go-sjc-assist-bot/internal/config"
)

var (
	appConfig config.Config
	debug     func(WebhookResponse)
)

func init() {
	appConfig = config.New()
}

func SetupWebhook(d func(WebhookResponse)) error {
	debug = d

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

func HandleWebhook(
	w http.ResponseWriter,
	r *http.Request,
	commands map[string]Command,
) error {
	w.WriteHeader(http.StatusOK)

	var webhookResponse WebhookResponse
	if err := json.NewDecoder(r.Body).Decode(&webhookResponse); err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}
	r.Body.Close()

	if os.Getenv("DEBUG") == "true" {
		log.Printf("DEBUG: Received message: %+v", webhookResponse.Message)
	}

	if debug != nil {
		debug(webhookResponse)
	}

	if len(commands) == 0 {
		log.Println("No commands available")
	}

	for command, handler := range commands {
		if strings.Contains(webhookResponse.Message.Text, command) {
			beforeCommandExecution(webhookResponse, command)
			if err := handler(&webhookResponse.Message); err != nil {
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

func SendMessage(chatID string, message string) error {
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

func SendDocument(chatID string, f *os.File) error {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	tmp, _ := writer.CreateFormFile("document", f.Name())
	f.Seek(0, 0)
	io.Copy(tmp, f)
	writer.WriteField("chat_id", chatID)
	writer.Close()

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", appConfig.TELEGRAM_TOKEN),
		buf,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send document: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to send document, status code: %d",
			response.StatusCode,
		)
	}

	return nil
}

func SendPhoto(chatID string, img *image.Image, caption string) error {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	tmp, _ := writer.CreateFormFile("photo", "image.png")
	if err := png.Encode(tmp, *img); err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}
	writer.WriteField("chat_id", chatID)
	if caption != "" {
		writer.WriteField("caption", caption)
	}
	writer.Close()

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", appConfig.TELEGRAM_TOKEN),
		buf,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send photo: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to send photo, status code: %d",
			response.StatusCode,
		)
	}

	return nil
}
