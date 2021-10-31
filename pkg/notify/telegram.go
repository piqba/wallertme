package notify

import (
	"context"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TgBotOption object options for telegram service
type TgBotOption struct {
	Debug bool
	Token string
}

// GetTgBot get an instance of tgBot
func GetTgBot(ctx context.Context, option TgBotOption) *tgbotapi.BotAPI {
	_, span := otel.Tracer(nameNotifierTG).Start(ctx, "GetTgBot")
	defer span.End()
	var token string
	if option.Token == "" {
		token = os.Getenv("BOT_TOKEN")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	if option.Debug {
		bot.Debug = option.Debug
	}
	span.SetAttributes(attribute.String("notifier.telegram.client", bot.Self.FirstName))

	return bot
}
