package handlers

import (
	"context"
	"disrbot/utils"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func LanguageHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		cb := update.CallbackQuery
		lang := "ar"
		if cb.Data == "setlang_en" {
			lang = "en"
		}

		_ = utils.RDB.Set(context.Background(), fmt.Sprintf("lang:%d", cb.From.ID), lang, 0).Err()
		_ = bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
			ChatID:    tu.ID(cb.Message.GetChat().ID),
			MessageID: cb.Message.GetMessageID(),
		})

		text := fmt.Sprintf(utils.Messages[lang]["lang_confirmed"], cb.From.FirstName)
		_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(cb.Message.GetChat().ID), text))
		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{CallbackQueryID: cb.ID})
		return nil
	}
}
