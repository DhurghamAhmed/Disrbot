package handlers

import (
	"disrbot/utils"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func IDHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}

		lang := utils.GetLang(msg.From.ID)

		targetUser := msg.From
		if msg.ReplyToMessage != nil {
			targetUser = msg.ReplyToMessage.From
		}

		username := "None"
		if lang == "ar" {
			username = "لا يوجد"
		}
		if targetUser.Username != "" {
			username = "@" + targetUser.Username
		}

		targetLang := utils.GetLang(targetUser.ID)

		res := fmt.Sprintf(
			utils.Messages[lang]["user_info_res"],
			targetUser.FirstName,
			username,
			targetUser.ID,
			targetLang,
		)

		sendText(bot, msg.Chat.ID, res)
		return nil
	}
}
