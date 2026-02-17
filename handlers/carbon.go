package handlers

import (
	"bytes"
	"context"
	"disrbot/utils"
	"strings" 

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func CarbonHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}

		lang := utils.GetLang(msg.From.ID)

		code := msg.Text
		code = strings.TrimPrefix(code, "/carbon")
		code = strings.TrimPrefix(code, "carbon")
		code = strings.TrimPrefix(code, "كاربون")
		code = strings.TrimSpace(code) 

		if code == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["carbon_description"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		waitMsg, _ := bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: tu.ID(msg.Chat.ID),
			Text:   utils.Messages[lang]["carbon_wait"],
		})

		imageBytes, err := utils.GenerateCarbonImage(code, "dracula")
		if err != nil {
			_, _ = bot.EditMessageText(context.Background(), &telego.EditMessageTextParams{
				ChatID:    tu.ID(msg.Chat.ID),
				MessageID: waitMsg.MessageID,
				Text:      utils.Messages[lang]["carbon_error"],
			})
			return err
		}

		_, _ = bot.SendPhoto(context.Background(), &telego.SendPhotoParams{
			ChatID: tu.ID(msg.Chat.ID),
			Photo:  tu.FileFromReader(bytes.NewReader(imageBytes), "carbon.png"),
			Caption: "Done!",
		})

		_ = bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			MessageID: waitMsg.MessageID,
		})

		return nil
	}
}