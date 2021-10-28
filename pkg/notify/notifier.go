package notify

import (
	"context"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	TELEGRAM = "telegram"
	DISCORD  = "discord"
	SMTP     = "smtp"
)

// Notifier ...
type Notifier interface {
	SendNotification(data interface{}) error
}

// SendMessageTG send msg to telegram
func SendMessageTG(bot *tgbotapi.BotAPI, id int64, message string) error {
	msg := tgbotapi.NewMessage(id, message)
	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil
}

// SendMessageDiscord send msg to discord
func SendMessageDiscord(discordClient DiscordClient, message string) error {

	err := discordClient.PostMessage(context.TODO(), message)
	if err != nil {
		return err
	}

	return nil
}

// SendMessageSMTP send msg to SMTP
func SendMessageSMTP(sender *Sender, message string) {
	//The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{os.Getenv("SMTP_EMAIL_RECEIVER")}

	Subject := "R2D2 notification service"

	bodyMessage := sender.WriteHTMLEmail(Receiver, Subject, message)

	sender.SendMail(Receiver, Subject, bodyMessage)
}
