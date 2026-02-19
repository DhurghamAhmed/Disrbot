package handlers

import (
	"context"
	"crypto/md5"
	"disrbot/utils"
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// (Moved helpers to utils/redis.go)

func AddVoiceHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		if msg.ReplyToMessage == nil {
			sendText(bot, msg.Chat.ID, utils.Messages[lang]["addvoice_reply_required"])
			return nil
		}

		reply := msg.ReplyToMessage
		if reply.Voice == nil && reply.Audio == nil {
			sendText(bot, msg.Chat.ID, utils.Messages[lang]["addvoice_voice_required"])
			return nil
		}

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/addvoice"))
		if name == "" {
			sendText(bot, msg.Chat.ID, utils.Messages[lang]["addvoice_name_required"])
			return nil
		}

		var fileID string
		if reply.Voice != nil {
			fileID = reply.Voice.FileID
		} else {
			fileID = reply.Audio.FileID
		}

		nameLower := strings.ToLower(name)
		utils.RDB.Set(context.Background(), utils.GlobalVoiceKey(nameLower), fileID, 0)
		utils.RDB.SAdd(context.Background(), utils.GlobalVoiceNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["addvoice_success"], name)
		sendText(bot, msg.Chat.ID, response)
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

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/delvoice"))
		if name == "" {
			sendText(bot, msg.Chat.ID, utils.Messages[lang]["delvoice_usage"])
			return nil
		}

		nameLower := strings.ToLower(name)
		exists, _ := utils.RDB.Exists(context.Background(), utils.GlobalVoiceKey(nameLower)).Result()
		if exists == 0 {
			response := fmt.Sprintf(utils.Messages[lang]["delvoice_notfound"], name)
			sendTextPlain(bot, msg.Chat.ID, response)
			return nil
		}

		utils.RDB.Del(context.Background(), utils.GlobalVoiceKey(nameLower))
		utils.RDB.SRem(context.Background(), utils.GlobalVoiceNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["delvoice_success"], name)
		sendTextPlain(bot, msg.Chat.ID, response)
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

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		names, err := utils.RDB.SMembers(context.Background(), utils.GlobalVoiceNamesKey).Result()
		if err != nil || len(names) == 0 {
			sendTextPlain(bot, msg.Chat.ID, utils.Messages[lang]["listvoices_empty"])
			return nil
		}

		var sb strings.Builder
		sb.WriteString(utils.Messages[lang]["listvoices_header"])
		sb.WriteString("\n\n")
		for i, name := range names {
			sb.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, name))
		}
		sendText(bot, msg.Chat.ID, sb.String())
		return nil
	}
}

// InlineVoiceHandler handles inline queries for both voices ("aud" prefix)
// and IPA files ("ipa" prefix).
func InlineVoiceHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		query := update.InlineQuery
		if query == nil {
			return nil
		}

		rawQuery := strings.TrimSpace(query.Query)
		lowerQuery := strings.ToLower(rawQuery)
		lang := utils.GetLang(query.From.ID)

		var results []telego.InlineQueryResult

		// If query is empty, show help result
		if lowerQuery == "" {
			results = append(results, &telego.InlineQueryResultArticle{
				Type:        telego.ResultTypeArticle,
				ID:          "inline_help",
				Title:       utils.Messages[lang]["inline_help_title"],
				Description: utils.Messages[lang]["inline_help_desc"],
				InputMessageContent: &telego.InputTextMessageContent{
					MessageText: utils.Messages[lang]["inline_help_text"],
					ParseMode:   telego.ModeMarkdown,
				},
			})
		}

		switch {
		case strings.HasPrefix(lowerQuery, "aud"):
			results = searchVoices(rawQuery[3:])
		case strings.HasPrefix(lowerQuery, "ipa"):
			results = searchIPAs(rawQuery[3:])
		}

		_ = bot.AnswerInlineQuery(context.Background(), &telego.AnswerInlineQueryParams{
			InlineQueryID: query.ID,
			Results:       results,
			CacheTime:     1,
		})
		return nil
	}
}

func searchVoices(rawSearch string) []telego.InlineQueryResult {
	searchText := strings.ToLower(strings.TrimSpace(rawSearch))

	names, err := utils.RDB.SMembers(context.Background(), utils.GlobalVoiceNamesKey).Result()
	if err != nil || len(names) == 0 {
		return nil
	}

	var results []telego.InlineQueryResult
	for _, name := range names {
		if searchText != "" && !strings.Contains(name, searchText) {
			continue
		}

		fileID, err := utils.RDB.Get(context.Background(), utils.GlobalVoiceKey(name)).Result()
		if err != nil {
			continue
		}

		results = append(results, &telego.InlineQueryResultCachedVoice{
			Type:        telego.ResultTypeVoice,
			ID:          fmt.Sprintf("%x", md5.Sum([]byte("voice_"+name))),
			VoiceFileID: fileID,
			Title:       name,
		})

		if len(results) >= 50 {
			break
		}
	}
	return results
}

func searchIPAs(rawSearch string) []telego.InlineQueryResult {
	searchText := strings.ToLower(strings.TrimSpace(rawSearch))

	names, err := utils.RDB.SMembers(context.Background(), utils.GlobalIpaNamesKey).Result()
	if err != nil || len(names) == 0 {
		return nil
	}

	var results []telego.InlineQueryResult
	for _, name := range names {
		if searchText != "" && !strings.Contains(name, searchText) {
			continue
		}

		fileID, err := utils.RDB.Get(context.Background(), utils.GlobalIpaKey(name)).Result()
		if err != nil {
			continue
		}

		results = append(results, &telego.InlineQueryResultCachedDocument{
			Type:           telego.ResultTypeDocument,
			ID:             fmt.Sprintf("%x", md5.Sum([]byte("ipa_"+name))),
			DocumentFileID: fileID,
			Title:          name,
			Description:    name + ".ipa",
		})

		if len(results) >= 50 {
			break
		}
	}
	return results
}
