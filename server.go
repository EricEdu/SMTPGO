package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type EmailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Erro ao carregar arquivo .env")
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_SENDER_PASSWORD"))

	router := gin.Default()

	router.POST("/send-email", func(c *gin.Context) {
		var emailRequest EmailRequest

		if err := c.ShouldBindJSON(&emailRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email := gomail.NewMessage()
		email.SetHeader("From", os.Getenv("EMAIL_SENDER_ADDRESS"))
		email.SetHeader("To", os.Getenv("EMAIL_RECEIVER_ADDRESS"))
		email.SetHeader("Subject", emailRequest.Subject)
		email.SetBody("text/plain", emailRequest.Body)

		if err := d.DialAndSend(email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email enviado com sucesso"})
	})

	router.Run(":8080")
}
