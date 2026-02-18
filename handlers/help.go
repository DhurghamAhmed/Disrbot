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

		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["btn_id_info"]).WithCallbackData("explain_id"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["btn_carbon_info"]).WithCallbackData("explain_carbon"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["btn_replies_info"]).WithCallbackData("explain_replies"),
			),
		)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:      tu.ID(update.Message.Chat.ID),
			Text:        utils.Messages[lang]["help_main"],
			ReplyMarkup: keyboard,
			ParseMode:   telego.ModeMarkdown,
		})
		return nil
	}
}

func ExplainIDHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		cb := update.CallbackQuery
		lang := utils.GetLang(cb.From.ID)

		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["back_btn"]).WithCallbackData("back_to_help"),
			),
		)

		_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:      tu.ID(cb.Message.GetChat().ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        utils.Messages[lang]["id_description"],
			ReplyMarkup: keyboard,
			ParseMode:   telego.ModeMarkdown,
		})

		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
		})
		return nil
	}
}

func ExplainCarbonHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		cb := update.CallbackQuery
		lang := utils.GetLang(cb.From.ID)

		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["back_btn"]).WithCallbackData("back_to_help"),
			),
		)

		_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:      tu.ID(cb.Message.GetChat().ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        utils.Messages[lang]["carbon_description"],
			ReplyMarkup: keyboard,
			ParseMode:   telego.ModeMarkdown,
		})

		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
		})
		return nil
	}
}

func ExplainRepliesHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		cb := update.CallbackQuery
		lang := utils.GetLang(cb.From.ID)

		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["back_btn"]).WithCallbackData("back_to_help"),
			),
		)

		_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:      tu.ID(cb.Message.GetChat().ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        utils.Messages[lang]["replies_description"],
			ReplyMarkup: keyboard,
			ParseMode:   telego.ModeMarkdown,
		})

		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
		})
		return nil
	}
}

func BackToHelpHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		cb := update.CallbackQuery
		lang := utils.GetLang(cb.From.ID)

		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["btn_id_info"]).WithCallbackData("explain_id"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["btn_carbon_info"]).WithCallbackData("explain_carbon"),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(utils.Messages[lang]["btn_replies_info"]).WithCallbackData("explain_replies"),
			),
		)

		_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
			ChatID:      tu.ID(cb.Message.GetChat().ID),
			MessageID:   cb.Message.GetMessageID(),
			Text:        utils.Messages[lang]["help_main"],
			ReplyMarkup: keyboard,
			ParseMode:   telego.ModeMarkdown,
		})

		_ = bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
		})
		return nil
	}
}
