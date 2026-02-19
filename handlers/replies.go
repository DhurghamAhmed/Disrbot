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

// --- Redis key helpers ---

func stateKey(userID int64) string     { return fmt.Sprintf("state:%d", userID) }
func stateDataKey(userID int64) string { return fmt.Sprintf("state_data:%d", userID) }

func replyKey(chatID int64, trigger string) string {
	return fmt.Sprintf("reply:%d:%s", chatID, trigger)
}

func triggersKey(chatID int64) string {
	return fmt.Sprintf("reply_triggers:%d", chatID)
}

func clearState(userID int64) {
	utils.RDB.Del(context.Background(), stateKey(userID))
	utils.RDB.Del(context.Background(), stateDataKey(userID))
}

// --- Command Handlers ---

func AddReplyHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		utils.RDB.Set(context.Background(), stateKey(msg.From.ID), "addreply_step1", 0)
		utils.RDB.Set(context.Background(), stateDataKey(msg.From.ID), fmt.Sprintf("%d", msg.Chat.ID), 0)

		sendText(bot, msg.Chat.ID, utils.Messages[lang]["addreply_ask_trigger"])
		return nil
	}
}

func DelReplyHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		utils.RDB.Set(context.Background(), stateKey(msg.From.ID), "delreply_step1", 0)
		utils.RDB.Set(context.Background(), stateDataKey(msg.From.ID), fmt.Sprintf("%d", msg.Chat.ID), 0)

		sendText(bot, msg.Chat.ID, utils.Messages[lang]["delreply_ask_trigger"])
		return nil
	}
}

func ListRepliesHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}
		lang := utils.GetLang(msg.From.ID)

		if !requireAdmin(bot, msg, lang) {
			return nil
		}

		triggers, err := utils.RDB.SMembers(context.Background(), triggersKey(msg.Chat.ID)).Result()
		if err != nil || len(triggers) == 0 {
			sendTextPlain(bot, msg.Chat.ID, utils.Messages[lang]["listreplies_empty"])
			return nil
		}

		var sb strings.Builder
		sb.WriteString(utils.Messages[lang]["listreplies_header"])
		sb.WriteString("\n\n")
		for _, trigger := range triggers {
			reply, _ := utils.RDB.Get(context.Background(), replyKey(msg.Chat.ID, trigger)).Result()
			sb.WriteString(fmt.Sprintf("• `%s` ← %s\n", trigger, reply))
		}
		sendText(bot, msg.Chat.ID, sb.String())
		return nil
	}
}

// --- State Machine ---

func StateHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil {
			return nil
		}

		// State machine needs either text (for triggers/replies) or document (for IPA)
		if msg.Text == "" && msg.Caption == "" && msg.Document == nil {
			return nil
		}

		lang := utils.GetLang(msg.From.ID)
		text := strings.TrimSpace(msg.Text)
		if text == "" {
			text = strings.TrimSpace(msg.Caption)
		}

		state, err := utils.RDB.Get(context.Background(), stateKey(msg.From.ID)).Result()
		if err != nil {
			return checkAutoReply(bot, msg)
		}

		// Cancel command
		if strings.EqualFold(text, "توقف") || strings.EqualFold(text, "stop") {
			clearState(msg.From.ID)
			sendTextPlain(bot, msg.Chat.ID, utils.Messages[lang]["operation_cancelled"])
			return nil
		}

		savedData, _ := utils.RDB.Get(context.Background(), stateDataKey(msg.From.ID)).Result()

		switch state {
		case "addreply_step1":
			trigger := strings.ToLower(text)
			utils.RDB.Set(context.Background(), stateDataKey(msg.From.ID), fmt.Sprintf("%s:%s", savedData, trigger), 0)
			utils.RDB.Set(context.Background(), stateKey(msg.From.ID), "addreply_step2", 0)

			response := fmt.Sprintf(utils.Messages[lang]["addreply_ask_reply"], trigger)
			sendText(bot, msg.Chat.ID, response)

		case "addreply_step2":
			parts := strings.SplitN(savedData, ":", 2)
			if len(parts) != 2 {
				clearState(msg.From.ID)
				return nil
			}

			var chatID int64
			fmt.Sscanf(parts[0], "%d", &chatID)
			trigger := parts[1]

			utils.RDB.Set(context.Background(), replyKey(chatID, trigger), text, 0)
			utils.RDB.SAdd(context.Background(), triggersKey(chatID), trigger)
			clearState(msg.From.ID)

			response := fmt.Sprintf(utils.Messages[lang]["addreply_success"], trigger, text)
			sendText(bot, msg.Chat.ID, response)

		case "delreply_step1":
			var chatID int64
			fmt.Sscanf(savedData, "%d", &chatID)
			trigger := strings.ToLower(text)
			key := replyKey(chatID, trigger)

			exists, _ := utils.RDB.Exists(context.Background(), key).Result()
			if exists == 0 {
				clearState(msg.From.ID)
				response := fmt.Sprintf(utils.Messages[lang]["delreply_notfound"], trigger)
				sendTextPlain(bot, msg.Chat.ID, response)
				return nil
			}

			utils.RDB.Del(context.Background(), key)
			utils.RDB.SRem(context.Background(), triggersKey(chatID), trigger)
			clearState(msg.From.ID)

			response := fmt.Sprintf(utils.Messages[lang]["delreply_success"], trigger)
			sendTextPlain(bot, msg.Chat.ID, response)

		case "addipa_step1":
			doc := msg.Document
			if doc == nil {
				sendText(bot, msg.Chat.ID, utils.Messages[lang]["addipa_doc_required"])
				return nil
			}

			if !strings.HasSuffix(strings.ToLower(doc.FileName), ".ipa") {
				sendText(bot, msg.Chat.ID, utils.Messages[lang]["addipa_ipa_required"])
				return nil
			}

			name := savedData
			if name == "" {
				name = strings.TrimSuffix(doc.FileName, ".ipa")
				name = strings.TrimSuffix(name, ".IPA")
			}

			nameLower := strings.ToLower(name)
			utils.RDB.Set(context.Background(), utils.GlobalIpaKey(nameLower), doc.FileID, 0)
			utils.RDB.SAdd(context.Background(), utils.GlobalIpaNamesKey, nameLower)
			clearState(msg.From.ID)

			response := fmt.Sprintf(utils.Messages[lang]["addipa_success"], name)
			sendText(bot, msg.Chat.ID, response)
		}

		return nil
	}
}

// --- Auto Reply ---

func checkAutoReply(bot *telego.Bot, msg *telego.Message) error {
	text := strings.ToLower(strings.TrimSpace(msg.Text))

	reply, err := utils.RDB.Get(context.Background(), replyKey(msg.Chat.ID, text)).Result()
	if err != nil {
		return nil
	}

	_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), reply))
	return nil
}
