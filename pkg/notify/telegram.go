package notify

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piqba/wallertme/pkg/logger"
)

// TgBotOption object options for telegram service
type TgBotOption struct {
	Debug bool
	Token string
}

// GetTgBot get an instance of tgBot
func GetTgBot(option TgBotOption) *tgbotapi.BotAPI {
	var token string
	if option.Token == "" {
		token = os.Getenv("BOT_TOKEN")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.LogError(err.Error())
	}

	if option.Debug {
		bot.Debug = option.Debug
	}
	logger.LogInfo(bot.Self.FirstName)

	return bot
}
