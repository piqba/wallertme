package r2d2

import (
	"context"
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piqba/wallertme/pkg/notify"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	// nameR2d2 is the Tracer nameR2d2 used to identify this instrumentation library.
	nameR2d2 = "R2d2.domain.tx"
)

// ExternalOptions ...
type ExternalOptions struct {
	Type              string
	DstNotificationID int64
}

// TxRepository ...
type TxRepository struct {
	Option        ExternalOptions
	TGClient      *tgbotapi.BotAPI
	DiscordClient notify.DiscordClient
	SMTPClient    notify.Sender
}

// NewTxRepository ...
func NewTxRepository(options ExternalOptions, clients ...interface{}) TxRepository {
	var repo TxRepository
	repo.Option = options
	for _, c := range clients {
		switch c := c.(type) {
		case *tgbotapi.BotAPI:
			repo.TGClient = c
		case notify.DiscordClient:
			repo.DiscordClient = c
		case notify.Sender:
			repo.SMTPClient = c
		}
	}
	return repo
}

// SendNotification ...
func (r *TxRepository) SendNotification(ctx context.Context, data interface{}) error {
	_, span := otel.Tracer(nameR2d2).Start(ctx, "SendNotification")
	defer span.End()
	t := reflect.TypeOf(data)
	if t == reflect.TypeOf(ResultLastTxADA{}) {

		tx := data.(ResultLastTxADA)
		switch r.Option.Type {

		case notify.TELEGRAM:

			err := notify.SendMessageTG(
				ctx,
				r.TGClient,
				r.Option.DstNotificationID,
				tx.TemplateTelegram(),
			)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return err
			}

		case notify.DISCORD:

			err := notify.SendMessageDiscord(
				ctx,
				r.DiscordClient,
				tx.TemplateDiscord(),
			)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return err
			}
		case notify.SMTP:
			notify.SendMessageSMTP(ctx, &r.SMTPClient, tx.TemplateSMTP())
		}
	} else if t == reflect.TypeOf(ResultLastTxSOL{}) {
		tx := data.(ResultLastTxSOL)
		switch r.Option.Type {

		case notify.TELEGRAM:

			err := notify.SendMessageTG(
				ctx,
				r.TGClient,
				r.Option.DstNotificationID,
				tx.TemplateTelegram(),
			)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return err
			}

		case notify.DISCORD:

			err := notify.SendMessageDiscord(
				ctx,
				r.DiscordClient,
				tx.TemplateDiscord(),
			)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return err
			}
		case notify.SMTP:
			notify.SendMessageSMTP(ctx, &r.SMTPClient, tx.TemplateSMTP())
		}
	}
	span.SetAttributes(attribute.String("r2d2.domain.SendNotification", "Success"))

	return nil
}
