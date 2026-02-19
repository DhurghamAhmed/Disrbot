package handlers

import (
	"context"
	"disrbot/utils"
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func globalVoiceKey(name string) string {
	return fmt.Sprintf("voice:global:%s", name)
}

const globalVoiceNamesKey = "voice_names:global"

func AddVoiceHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !utils.IsAdmin(msg.From.ID) {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["not_admin"]))
			return nil
		}

		if msg.ReplyToMessage == nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["addvoice_reply_required"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		reply := msg.ReplyToMessage

		if reply.Voice == nil && reply.Audio == nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["addvoice_voice_required"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/addvoice"))
		if name == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["addvoice_name_required"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		var fileID string
		if reply.Voice != nil {
			fileID = reply.Voice.FileID
		} else {
			fileID = reply.Audio.FileID
		}

		nameLower := strings.ToLower(name)

		utils.RDB.Set(context.Background(), globalVoiceKey(nameLower), fileID, 0)
		utils.RDB.SAdd(context.Background(), globalVoiceNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["addvoice_success"], name)
		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			Text:      response,
			ParseMode: telego.ModeMarkdown,
		})
		return nil
	}
}

func DelVoiceHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !utils.IsAdmin(msg.From.ID) {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["not_admin"]))
			return nil
		}

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/delvoice"))
		if name == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["delvoice_usage"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		nameLower := strings.ToLower(name)

		exists, _ := utils.RDB.Exists(context.Background(), globalVoiceKey(nameLower)).Result()
		if exists == 0 {
			response := fmt.Sprintf(utils.Messages[lang]["delvoice_notfound"], name)
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), response))
			return nil
		}

		utils.RDB.Del(context.Background(), globalVoiceKey(nameLower))
		utils.RDB.SRem(context.Background(), globalVoiceNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["delvoice_success"], name)
		_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), response))
		return nil
	}
}

func ListVoicesHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !utils.IsAdmin(msg.From.ID) {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["not_admin"]))
			return nil
		}

		names, err := utils.RDB.SMembers(context.Background(), globalVoiceNamesKey).Result()
		if err != nil || len(names) == 0 {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["listvoices_empty"]))
			return nil
		}

		var sb strings.Builder
		sb.WriteString(utils.Messages[lang]["listvoices_header"])
		sb.WriteString("\n\n")
		for i, name := range names {
			sb.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, name))
		}
		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			Text:      sb.String(),
			ParseMode: telego.ModeMarkdown,
		})
		return nil
	}
}

func InlineVoiceHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		query := update.InlineQuery
		if query == nil {
			return nil
		}

		rawQuery := strings.TrimSpace(query.Query)

		if !strings.HasPrefix(strings.ToLower(rawQuery), "aud") {
			_ = bot.AnswerInlineQuery(context.Background(), &telego.AnswerInlineQueryParams{
				InlineQueryID: query.ID,
				Results:       []telego.InlineQueryResult{},
				CacheTime:     1,
			})
			return nil
		}

		searchText := strings.ToLower(strings.TrimSpace(rawQuery[3:]))

		names, err := utils.RDB.SMembers(context.Background(), globalVoiceNamesKey).Result()
		if err != nil || len(names) == 0 {
			_ = bot.AnswerInlineQuery(context.Background(), &telego.AnswerInlineQueryParams{
				InlineQueryID: query.ID,
				Results:       []telego.InlineQueryResult{},
				CacheTime:     1,
			})
			return nil
		}

		var results []telego.InlineQueryResult

		for _, name := range names {
			if searchText != "" && !strings.Contains(name, searchText) {
				continue
			}

			fileID, err := utils.RDB.Get(context.Background(), globalVoiceKey(name)).Result()
			if err != nil {
				continue
			}

			results = append(results, &telego.InlineQueryResultCachedVoice{
				Type:        telego.ResultTypeVoice,
				ID:          name,
				VoiceFileID: fileID,
				Title:       name,
			})

			if len(results) >= 50 {
				break
			}
		}

		_ = bot.AnswerInlineQuery(context.Background(), &telego.AnswerInlineQueryParams{
			InlineQueryID: query.ID,
			Results:       results,
			CacheTime:     1,
		})
		return nil
	}
}
