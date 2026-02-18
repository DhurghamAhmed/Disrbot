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



func stateKey(userID int64) string     { return fmt.Sprintf("state:%d", userID) }
func stateDataKey(userID int64) string { return fmt.Sprintf("state_data:%d", userID) }

func clearState(userID int64) {
	utils.RDB.Del(context.Background(), stateKey(userID))
	utils.RDB.Del(context.Background(), stateDataKey(userID))
}


func AddReplyHandler(bot *telego.Bot) th.Handler {
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

		utils.RDB.Set(context.Background(), stateKey(msg.From.ID), "addreply_step1", 0)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			Text:      utils.Messages[lang]["addreply_ask_trigger"],
			ParseMode: telego.ModeMarkdown,
		})
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

		if !utils.IsAdmin(msg.From.ID) {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["not_admin"]))
			return nil
		}

		utils.RDB.Set(context.Background(), stateKey(msg.From.ID), "delreply_step1", 0)

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			Text:      utils.Messages[lang]["delreply_ask_trigger"],
			ParseMode: telego.ModeMarkdown,
		})
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

		if !utils.IsAdmin(msg.From.ID) {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["not_admin"]))
			return nil
		}

		triggers, err := utils.RDB.SMembers(context.Background(), "reply_triggers").Result()
		if err != nil || len(triggers) == 0 {
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["listreplies_empty"]))
			return nil
		}

		var sb strings.Builder
		sb.WriteString(utils.Messages[lang]["listreplies_header"])
		sb.WriteString("\n\n")
		for _, trigger := range triggers {
			reply, _ := utils.RDB.Get(context.Background(), fmt.Sprintf("reply:%s", trigger)).Result()
			sb.WriteString(fmt.Sprintf("• `%s` ← %s\n", trigger, reply))
		}

		_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID:    tu.ID(msg.Chat.ID),
			Text:      sb.String(),
			ParseMode: telego.ModeMarkdown,
		})
		return nil
	}
}


func StateHandler(bot *telego.Bot) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		msg := update.Message
		if msg == nil || msg.Text == "" {
			return nil
		}

		lang := utils.GetLang(msg.From.ID)
		text := strings.TrimSpace(msg.Text)

		state, err := utils.RDB.Get(context.Background(), stateKey(msg.From.ID)).Result()
		if err != nil {
			return checkAutoReply(bot, msg)
		}

		if strings.EqualFold(text, "توقف") || strings.EqualFold(text, "stop") {
			clearState(msg.From.ID)
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), utils.Messages[lang]["operation_cancelled"]))
			return nil
		}

		switch state {

		case "addreply_step1":
			trigger := strings.ToLower(text)
			utils.RDB.Set(context.Background(), stateDataKey(msg.From.ID), trigger, 0)
			utils.RDB.Set(context.Background(), stateKey(msg.From.ID), "addreply_step2", 0)

			response := fmt.Sprintf(utils.Messages[lang]["addreply_ask_reply"], trigger)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      response,
				ParseMode: telego.ModeMarkdown,
			})

		case "addreply_step2":
			trigger, _ := utils.RDB.Get(context.Background(), stateDataKey(msg.From.ID)).Result()
			reply := text

			utils.RDB.Set(context.Background(), fmt.Sprintf("reply:%s", trigger), reply, 0)
			utils.RDB.SAdd(context.Background(), "reply_triggers", trigger)
			clearState(msg.From.ID)

			response := fmt.Sprintf(utils.Messages[lang]["addreply_success"], trigger, reply)
			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    tu.ID(msg.Chat.ID),
				Text:      response,
				ParseMode: telego.ModeMarkdown,
			})

		case "delreply_step1":
			trigger := strings.ToLower(text)
			key := fmt.Sprintf("reply:%s", trigger)

			exists, _ := utils.RDB.Exists(context.Background(), key).Result()
			if exists == 0 {
				clearState(msg.From.ID)
				response := fmt.Sprintf(utils.Messages[lang]["delreply_notfound"], trigger)
				_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), response))
				return nil
			}

			utils.RDB.Del(context.Background(), key)
			utils.RDB.SRem(context.Background(), "reply_triggers", trigger)
			clearState(msg.From.ID)

			response := fmt.Sprintf(utils.Messages[lang]["delreply_success"], trigger)
			_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), response))
		}

		return nil
	}
}


func checkAutoReply(bot *telego.Bot, msg *telego.Message) error {
	text := strings.ToLower(strings.TrimSpace(msg.Text))
	reply, err := utils.RDB.Get(context.Background(), fmt.Sprintf("reply:%s", text)).Result()
	if err != nil {
		return nil
	}

	_, _ = bot.SendMessage(context.Background(), tu.Message(tu.ID(msg.Chat.ID), reply))
	return nil
}