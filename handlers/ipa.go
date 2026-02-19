package handlers

import (
	"context"
	"disrbot/utils"
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// (Moved helpers to utils/redis.go)

func AddIpaHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/addipa"))
		var doc *telego.Document

		// 1. Check if it's a reply with a document
		if msg.ReplyToMessage != nil && msg.ReplyToMessage.Document != nil {
			doc = msg.ReplyToMessage.Document
		}

		// 2. Check if the message itself has a document (caption)
		if doc == nil && msg.Document != nil {
			doc = msg.Document
		}

		// 3. If no document found, enter state to wait for it
		if doc == nil {
			utils.RDB.Set(context.Background(), stateKey(msg.From.ID), "addipa_step1", 0)
			utils.RDB.Set(context.Background(), stateDataKey(msg.From.ID), name, 0)
			sendText(bot, msg.Chat.ID, utils.Messages[lang]["addipa_ask_doc"])
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(doc.FileName), ".ipa") {
			sendText(bot, msg.Chat.ID, utils.Messages[lang]["addipa_ipa_required"])
			return nil
		}

		if name == "" {
			name = strings.TrimSuffix(doc.FileName, ".ipa")
			name = strings.TrimSuffix(name, ".IPA")
			if name == "" {
				sendText(bot, msg.Chat.ID, utils.Messages[lang]["addipa_name_required"])
				return nil
			}
		}

		nameLower := strings.ToLower(name)
		utils.RDB.Set(context.Background(), utils.GlobalIpaKey(nameLower), doc.FileID, 0)
		utils.RDB.SAdd(context.Background(), utils.GlobalIpaNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["addipa_success"], name)
		sendText(bot, msg.Chat.ID, response)
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

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		name := strings.TrimSpace(strings.TrimPrefix(msg.Text, "/delipa"))
		if name == "" {
			sendText(bot, msg.Chat.ID, utils.Messages[lang]["delipa_usage"])
			return nil
		}

		nameLower := strings.ToLower(name)
		exists, _ := utils.RDB.Exists(context.Background(), utils.GlobalIpaKey(nameLower)).Result()
		if exists == 0 {
			response := fmt.Sprintf(utils.Messages[lang]["delipa_notfound"], name)
			sendTextPlain(bot, msg.Chat.ID, response)
			return nil
		}

		utils.RDB.Del(context.Background(), utils.GlobalIpaKey(nameLower))
		utils.RDB.SRem(context.Background(), utils.GlobalIpaNamesKey, nameLower)

		response := fmt.Sprintf(utils.Messages[lang]["delipa_success"], name)
		sendTextPlain(bot, msg.Chat.ID, response)
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

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		names, err := utils.RDB.SMembers(context.Background(), utils.GlobalIpaNamesKey).Result()
		if err != nil || len(names) == 0 {
			sendTextPlain(bot, msg.Chat.ID, utils.Messages[lang]["listipa_empty"])
			return nil
		}

		var sb strings.Builder
		sb.WriteString(utils.Messages[lang]["listipa_header"])
		sb.WriteString("\n\n")
		for i, name := range names {
			sb.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, name))
		}
		sendText(bot, msg.Chat.ID, sb.String())
		return nil
	}
}
