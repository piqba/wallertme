package r2d2

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piqba/wallertme/pkg/notify"
)

type ExternalOptions struct {
	Type              string
	DstNotificationID int64
}

type TxRepository struct {
	Option        ExternalOptions
	TGClient      *tgbotapi.BotAPI
	DiscordClient notify.DiscordClient
}

func NewTxRepository(options ExternalOptions, clients ...interface{}) TxRepository {
	var repo TxRepository
	repo.Option = options
	for _, c := range clients {
		switch c := c.(type) {
		case *tgbotapi.BotAPI:
			repo.TGClient = c
		case notify.DiscordClient:
			repo.DiscordClient = c
		}
	}
	return repo
}

func (r *TxRepository) SendNotification(data string) error {
	tx := ResultLastTxByAddr{}

	err := json.Unmarshal([]byte(data), &tx)
	if err != nil {
		return err
	}

	switch r.Option.Type {

	case notify.TELEGRAM:

		err = notify.SendMessageTG(
			r.TGClient,
			r.Option.DstNotificationID,
			tx.Hummanify(),
		)
		if err != nil {
			return err
		}

	case notify.DISCORD:

		err = notify.SendMessageDiscord(
			r.DiscordClient,
			tx.EmbedDiscord(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
