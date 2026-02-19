package handlers

import (
	"context"
	"disrbot/utils"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HelpHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		userID := update.Message.From.ID
		lang := utils.GetLang(userID)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:      tu.ID(update.Message.Chat.ID),
			Text:        utils.Messages[lang]["help_main"],
			ReplyMarkup: helpKeyboard(lang),
			ParseMode:   telego.ModeMarkdown,
		})
		return nil
	}
}

// explainHandler is a generic factory for all "explain" callback handlers.
// Each explain button shows a specific message with a "back" button.
func explainHandler(bot *telego.Bot, messageKey string) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		cb := update.CallbackQuery
		lang := utils.GetLang(cb.From.ID)

		backKeyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["back_btn"]).WithCallbackData("back_to_help"),
			),
		)

		_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:      tu.ID(cb.Message.GetChat().ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        utils.Messages[lang][messageKey],
			ReplyMarkup: backKeyboard,
			ParseMode:   telego.ModeMarkdown,
		})
		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
		})
		return nil
	}
}

func ExplainIDHandler(bot *telego.Bot) th.Handler { return explainHandler(bot, "id_description") }
func ExplainCarbonHandler(bot *telego.Bot) th.Handler {
	return explainHandler(bot, "carbon_description")
}
func ExplainRepliesHandler(bot *telego.Bot) th.Handler {
	return explainHandler(bot, "replies_description")
}
func ExplainVoicesHandler(bot *telego.Bot) th.Handler {
	return explainHandler(bot, "voices_description")
}

func BackToHelpHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		cb := update.CallbackQuery
		lang := utils.GetLang(cb.From.ID)

		_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:      tu.ID(cb.Message.GetChat().ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        utils.Messages[lang]["help_main"],
			ReplyMarkup: helpKeyboard(lang),
			ParseMode:   telego.ModeMarkdown,
		})
		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
		})
		return nil
	}
}
