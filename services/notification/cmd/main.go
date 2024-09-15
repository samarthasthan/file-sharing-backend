package main

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/env"
	"github.com/samarthasthan/21BRS1248_Backend/common/kafka"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	mail "github.com/samarthasthan/21BRS1248_Backend/services/notification/internal"
)

var (
	SMTP_SERVER    string
	SMTP_PORT      string
	SMTP_LOGIN     string
	SMTP_PASSWORD  string
	KAFKA_PORT     string
	KAFKA_HOST     string
)

func init() {
	SMTP_SERVER = env.GetEnv("SMTP_SERVER", "smtp-relay.sendinblue.com")
	SMTP_PORT = env.GetEnv("SMTP_PORT", "587")
	SMTP_LOGIN = env.GetEnv("SMTP_LOGIN", "use your own sender")
	SMTP_PASSWORD = env.GetEnv("SMTP_PASSWORD", "use your own key")
	KAFKA_PORT = env.GetEnv("KAFKA_PORT", "9092")
	KAFKA_HOST = env.GetEnv("KAFKA_HOST", "localhost")
}

func main() {
	// Initialising Custom Logger
	log := logger.NewLogger("Mail")

	// Initialising Kafka Consumer
	k := kafka.NewKafkaConsumer(KAFKA_HOST, KAFKA_PORT)

	// Initialising Mail Handler
	m := mail.NewMailHandler(k, log, SMTP_SERVER, SMTP_PORT, SMTP_LOGIN, SMTP_PASSWORD)
	// Start sending mails
	m.SendMails()
}