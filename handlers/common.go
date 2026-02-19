package handlers

import (
	"context"
	"disrbot/utils"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// sendText sends a Markdown-formatted message.
func sendText(bot *telego.Bot, chatID int64, text string) {
	_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID:    tu.ID(chatID),
		Text:      text,
		ParseMode: telego.ModeMarkdown,
	})
}

// sendTextPlain sends a plain-text message (no parse mode).
func sendTextPlain(bot *telego.Bot, chatID int64, text string) {
	_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(chatID), text))
}

// requireAdmin checks if the sender is an admin. If not, it sends a
// "not_admin" message and returns false.
func requireAdmin(bot *telego.Bot, msg *telego.Message, lang string) bool {
	if utils.IsAdmin(msg.From.ID) {
		return true
	}
	sendTextPlain(bot, msg.Chat.ID, utils.Messages[lang]["not_admin"])
	return false
}

// helpKeyboard builds the help menu inline keyboard.
func helpKeyboard(lang string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(utils.Messages[lang]["btn_id_info"]).WithCallbackData("explain_id"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(utils.Messages[lang]["btn_carbon_info"]).WithCallbackData("explain_carbon"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(utils.Messages[lang]["btn_replies_info"]).WithCallbackData("explain_replies"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(utils.Messages[lang]["btn_voices_info"]).WithCallbackData("explain_voices"),
		),
	)
}
