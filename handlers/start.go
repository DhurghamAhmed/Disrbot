package handlers

import (
	"context"
	"disrbot/utils"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func StartHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		userID := update.Message.From.ID
		_ = utils.RDB.SAdd(context.Background(), "disr_bot_users", userID).Err()

		keyboard := tu.InlineKeyboard(tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("العربية").WithCallbackData("setlang_ar"),
			tu.InlineKeyboardButton("English").WithCallbackData("setlang_en"),
		))

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:      tu.ID(update.Message.Chat.ID),
			Text:        "أهلاً بك! يرجى اختيار لغتك:\nWelcome! Choose your language:",
			ReplyMarkup: keyboard,
		})
		return nil
	}
}
