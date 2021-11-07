package notify

import (
	"context"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	TELEGRAM = "telegram"
	DISCORD  = "discord"
	WEBHOOK  = "webhook"
	SMTP     = "smtp"
	// nameNotifier is the Tracer nameNotifier used to identify this instrumentation library.
	nameNotifier = "notify.notifier"
	// nameNotifierTG is the Tracer nameNotifierTG used to identify this instrumentation library.
	nameNotifierTG = "notify.notifier.tg"
	// nameNotifierSMTP is the Tracer nameNotifierSMTP used to identify this instrumentation library.
	nameNotifierSMTP = "notify.notifier.smtp"
	// nameNotifierDiscord is the Tracer nameNotifierDiscord used to identify this instrumentation library.
	nameNotifierDiscord = "notify.notifier.discord"
	// nameNotifierWebHook is the Tracer nameNotifierWebHook used to identify this instrumentation library.
	nameNotifierWebHook = "notify.notifier.webhook"
)

// Notifier ...
type Notifier interface {
	SendNotification(ctx context.Context, data interface{}) error
}

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// SendMessageTG send msg to telegram
func SendMessageTG(ctx context.Context, bot *tgbotapi.BotAPI, id int64, message string) error {
	_, span := otel.Tracer(nameNotifier).Start(ctx, "SendMessageTG")
	defer span.End()
	msg := tgbotapi.NewMessage(id, message)
	if _, err := bot.Send(msg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetAttributes(attribute.String("notifier.telegram", "Success"))

	return nil
}

// SendMessageDiscord send msg to discord
func SendMessageDiscord(ctx context.Context, discordClient DiscordClient, message string) error {
	_, span := otel.Tracer(nameNotifier).Start(ctx, "SendMessageDiscord")
	defer span.End()
	err := discordClient.PostMessage(context.TODO(), message)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetAttributes(attribute.String("notifier.discord", "Success"))

	return nil
}

// SendMessageWebHook send msg to webhook
func SendMessageWebHook(ctx context.Context, webhookClient WebHookClient, message string) error {
	_, span := otel.Tracer(nameNotifier).Start(ctx, "SendMessageWebHook")
	defer span.End()
	err := webhookClient.PostMessage(context.TODO(), message)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetAttributes(attribute.String("notifier.webhook", "Success"))

	return nil
}

// SendMessageSMTP send msg to SMTP
func SendMessageSMTP(ctx context.Context, sender *Sender, message string) {
	_, span := otel.Tracer(nameNotifier).Start(ctx, "SendMessageSMTP")
	defer span.End()
	//The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{os.Getenv("SMTP_EMAIL_RECEIVER")}

	Subject := "R2D2 notification service"

	bodyMessage := sender.WriteHTMLEmail(Receiver, Subject, message)

	sender.SendMail(ctx, Receiver, Subject, bodyMessage)
	span.SetAttributes(attribute.String("notifier.smtp", "Success"))

}
