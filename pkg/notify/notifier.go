package notify

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	TELEGRAM = "telegram"
	DISCORD  = "discord"
	SMTP     = "smtp"
)

type Notifier interface {
	SendNotification(data interface{}) error
}

func SendMessageTG(bot *tgbotapi.BotAPI, id int64, message string) error {
	msg := tgbotapi.NewMessage(id, message)
	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func SendMessageDiscord(discordClient DiscordClient, message string) error {

	err := discordClient.PostMessage(context.TODO(), message)
	if err != nil {
		return err
	}

	return nil
}
