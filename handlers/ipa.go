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

func globalIpaKey(name string) string {
	return fmt.Sprintf("ipa:global:%s", name)
}

const globalIpaNamesKey = "ipa_names:global"

func AddIpaHandler(bot *telego.Bot) th.Handler {
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
				Text:      utils.Messages[lang]["addipa_reply_required"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		reply := msg.ReplyToMessage

		if reply.Document == nil {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["addipa_doc_required"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		doc := reply.Document
		if !strings.HasSuffix(strings.ToLower(doc.FileName), ".ipa") {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["addipa_ipa_required"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/addipa"))
		if name == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["addipa_name_required"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		nameLower := strings.ToLower(name)

		utils.RDB.Set(context.Background(), globalIpaKey(nameLower), doc.FileID, 0)
		utils.RDB.SAdd(context.Background(), globalIpaNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["addipa_success"], name)
		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			Text:      response,
			ParseMode: telego.ModeMarkdown,
		})
		return nil
	}
}

func DelIpaHandler(bot *telego.Bot) th.Handler {
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

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/delipa"))
		if name == "" {
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      utils.Messages[lang]["delipa_usage"],
				ParseMode: telego.ModeMarkdown,
			})
			return nil
		}

		nameLower := strings.ToLower(name)
		exists, _ := utils.RDB.Exists(context.Background(), globalIpaKey(nameLower)).Result()
		if exists == 0 {
			response := fmt.Sprintf(utils.Messages[lang]["delipa_notfound"], name)
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), response))
			return nil
		}

		utils.RDB.Del(context.Background(), globalIpaKey(nameLower))
		utils.RDB.SRem(context.Background(), globalIpaNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["delipa_success"], name)
		_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), response))
		return nil
	}
}

func ListIpaHandler(bot *telego.Bot) th.Handler {
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

		names, err := utils.RDB.SMembers(context.Background(), globalIpaNamesKey).Result()
		if err != nil || len(names) == 0 {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["listipa_empty"]))
			return nil
		}

		var sb strings.Builder
		sb.WriteString(utils.Messages[lang]["listipa_header"])
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

func InlineIpaHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		query := update.InlineQuery
		if query == nil {
			return nil
		}

		rawQuery := strings.TrimSpace(query.Query)

		if !strings.HasPrefix(strings.ToLower(rawQuery), "ipa") {
			_ = bot.AnswerInlineQuery(context.Background(), &telego.AnswerInlineQueryParams{
				InlineQueryID: query.ID,
				Results:       []telego.InlineQueryResult{},
				CacheTime:     1,
			})
			return nil
		}

		searchText := strings.ToLower(strings.TrimSpace(rawQuery[3:]))

		names, err := utils.RDB.SMembers(context.Background(), globalIpaNamesKey).Result()
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

			fileID, err := utils.RDB.Get(context.Background(), globalIpaKey(name)).Result()
			if err != nil {
				continue
			}

			results = append(results, &telego.InlineQueryResultCachedDocument{
				Type:           telego.ResultTypeDocument,
				ID:             "ipa_" + name,
				DocumentFileID: fileID,
				Title:          name,
				Description:    name + ".ipa",
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
