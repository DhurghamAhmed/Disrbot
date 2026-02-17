package handlers

import (
	"context"
	"disrbot/utils"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func IDHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {

		msg := update.Message
		if msg == nil {
			return nil
		}

		lang := utils.GetLang(msg.From.ID)

		var targetUser *telego.User

		if msg.ReplyToMessage != nil {
			targetUser = msg.ReplyToMessage.From
		} else {
			targetUser = msg.From
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

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			Text:      res,
			ParseMode: telego.ModeMarkdown,
		})

		return nil
	}
}
