package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type EmailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

var (
	dialer *gomail.Dialer
	once   sync.Once
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Erro ao carregar arquivo .env")
	}
}

func getDialer() *gomail.Dialer {
	once.Do(func() {
		loadEnv()
		dialer = gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_SENDER_PASSWORD"))
	})
	return dialer
}

func sendEmail(email *gomail.Message, d *gomail.Dialer, done chan<- bool) {
	if err := d.DialAndSend(email); err != nil {
		panic(err)
	}
	done <- true
}

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	loadEnv()

	// Criar dialer SMTP
	d := getDialer()

	http.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var emailRequest EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&emailRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		email := gomail.NewMessage()
		email.SetHeader("From", os.Getenv("EMAIL_SENDER_ADDRESS"))
		email.SetHeader("To", os.Getenv("EMAIL_RECEIVER_ADDRESS"))
		email.SetHeader("Subject", emailRequest.Subject)
		email.SetBody("text/plain", emailRequest.Body)

		done := make(chan bool)
		go sendEmail(email, d, done)

		<-done
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Email enviado com sucesso"}`))
	})

	http.ListenAndServe(":8080", nil)
}
