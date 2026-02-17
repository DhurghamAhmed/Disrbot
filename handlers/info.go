package handlers

import (
	"context"
	"disrbot/utils"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func InfoHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		lang := utils.GetLang(userID)

		_, err := bot.SendMessage(context.Background(), tu.Message(
			tu.ID(chatID),
			utils.Messages[lang]["info_msg"], 
		))
		
		return err
	}
}