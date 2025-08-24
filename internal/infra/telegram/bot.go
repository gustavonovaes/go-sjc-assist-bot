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
	"net/url"
	"os"
	"strconv"
	"strings"

	"gustavonovaes.dev/go-sjc-assist-bot/internal/config"
)

var (
	appConfig  config.Config
	middleware func(WebhookResponse) WebhookResponse
	commands   map[string]Command
)

func init() {
	appConfig = config.New()
}

func NewWebhookServer(commandHandlers map[string]Command, middlewareFunc func(WebhookResponse) WebhookResponse) *http.ServeMux {
	middleware = middlewareFunc
	commands = commandHandlers

	webhookURL, err := url.Parse(appConfig.TELEGRAM_WEBHOOK_URL)
	if err != nil {
		log.Fatalf("ERROR: Invalid TELEGRAM_WEBHOOK_URL: %v", err)
	}

	server := http.NewServeMux()

	pattern := fmt.Sprintf("POST %s", webhookURL.Path)
	server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != appConfig.TELEGRAM_SECRET_TOKEN {
			log.Printf("WARN: Invalid secret token in request header")
			http.Error(w, "Forbidden: Invalid secret token", http.StatusForbidden)
			return
		}

		err := handleWebhookRequest(r)
		if err != nil {
			log.Printf("ERROR: Failed to handle webhook request: %v", err)
			w.WriteHeader(http.StatusOK)
			// http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Println("INFO: Webhook request handled successfully")
	})

	return server
}

func handleWebhookRequest(r *http.Request) error {
	var webhookResponse WebhookResponse
	if err := json.NewDecoder(r.Body).Decode(&webhookResponse); err != nil {
		bodyContent, _ := io.ReadAll(r.Body)
		return fmt.Errorf("failed to decode request body: %v\n%+v", err, bodyContent)
	}
	r.Body.Close()

	if middleware != nil {
		log.Println("INFO: Applying middleware to webhook response")
		webhookResponse = middleware(webhookResponse)
	}

	if len(commands) == 0 {
		log.Println("INFO: No commands available")
		return nil
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
				return err
			}

			afterCommandExecution(webhookResponse, command)

			return nil
		}
	}

	return nil
}

func beforeCommandExecution(w WebhookResponse, command string) {
	if os.Getenv("DEBUG") != "" {
		log.Printf(
			"DEBUG: User %s/%d requested command %s from in chat %d",
			w.Message.From.Username,
			w.Message.From.ID,
			command,
			w.Message.Chat.ID,
		)
	}
}

func afterCommandExecution(w WebhookResponse, command string) {
	//
}

func SetupWebhook() error {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/setWebhook?url=%s&secret_token=%s",
		appConfig.TELEGRAM_API_TOKEN,
		appConfig.TELEGRAM_WEBHOOK_URL,
		appConfig.TELEGRAM_SECRET_TOKEN,
	)

	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyContent := make([]byte, 256)
	response.Body.Read(bodyContent)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to setup webhook, status code: %d, body: %s",
			response.StatusCode,
			bodyContent,
		)
	}

	if os.Getenv("DEBUG") != "" {
		log.Printf("DEBUG: Webhook setup response: %s", bodyContent)
	}

	return nil
}

func SendMessage(chatID int, message string) error {
	res, err := http.Post(
		"https://api.telegram.org/bot"+appConfig.TELEGRAM_API_TOKEN+"/sendMessage",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{ "chat_id": "%d",  "text": "%s", "parse_mode": "%s", "disable_web_page_preview": true, }`, chatID, message, "Markdown")),
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

func SendDocument(chatID int, f *os.File) error {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	tmp, _ := writer.CreateFormFile("document", f.Name())
	f.Seek(0, 0)
	io.Copy(tmp, f)
	writer.WriteField("chat_id", strconv.Itoa(chatID))
	writer.Close()

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", appConfig.TELEGRAM_API_TOKEN),
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

func SendPhoto(chatID int, img *image.Image, caption string) error {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	tmp, _ := writer.CreateFormFile("photo", "image.png")
	if err := png.Encode(tmp, *img); err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}
	writer.WriteField("chat_id", strconv.Itoa(chatID))
	if caption != "" {
		writer.WriteField("caption", caption)
	}
	writer.Close()

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", appConfig.TELEGRAM_API_TOKEN),
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
