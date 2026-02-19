package main

import (
	"context"
	"log"
	"os"
	"strings"

	"disrbot/handlers"
	"disrbot/utils"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = godotenv.Load()
	botToken := os.Getenv("TOKEN")

	utils.InitRedis()

	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatal(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		AllowedUpdates: []string{"message", "callback_query", "inline_query"},
	})
	bh, _ := th.NewBotHandler(bot, updates)

	bh.Handle(handlers.StartHandler(bot), th.CommandEqual("start"))
	bh.Handle(handlers.LanguageHandler(bot), th.CallbackDataPrefix("setlang_"))

	bh.Handle(handlers.HelpHandler(bot), th.Or(th.CommandEqual("help"), th.TextEqual("مساعدة")))
	bh.Handle(handlers.ExplainIDHandler(bot), th.CallbackDataEqual("explain_id"))
	bh.Handle(handlers.ExplainCarbonHandler(bot), th.CallbackDataEqual("explain_carbon"))
	bh.Handle(handlers.ExplainRepliesHandler(bot), th.CallbackDataEqual("explain_replies"))
	bh.Handle(handlers.ExplainVoicesHandler(bot), th.CallbackDataEqual("explain_voices"))
	bh.Handle(handlers.BackToHelpHandler(bot), th.CallbackDataEqual("back_to_help"))

	bh.Handle(handlers.IDHandler(bot), th.Or(
		th.CommandEqual("id"),
		th.TextEqual("id"),
		th.TextEqual("ايدي"),
	))

	bh.Handle(handlers.CarbonHandler(bot), th.Or(
		th.CommandEqual("carbon"),
		func(ctx context.Context, update telego.Update) bool {
			if update.Message != nil {
				txt := strings.ToLower(update.Message.Text)
				return strings.HasPrefix(txt, "كاربون") || strings.HasPrefix(txt, "carbon")
			}
			return false
		},
	))

	bh.Handle(handlers.AddReplyHandler(bot), th.CommandEqual("addreply"))
	bh.Handle(handlers.DelReplyHandler(bot), th.CommandEqual("delreply"))
	bh.Handle(handlers.ListRepliesHandler(bot), th.CommandEqual("listreplies"))

	bh.Handle(handlers.AddVoiceHandler(bot), th.CommandEqual("addvoice"))
	bh.Handle(handlers.DelVoiceHandler(bot), th.CommandEqual("delvoice"))
	bh.Handle(handlers.ListVoicesHandler(bot), th.CommandEqual("listvoices"))

	bh.Handle(handlers.AddIpaHandler(bot), th.CommandEqual("addipa"))
	bh.Handle(handlers.DelIpaHandler(bot), th.CommandEqual("delipa"))
	bh.Handle(handlers.ListIpaHandler(bot), th.CommandEqual("listipa"))

	bh.Handle(handlers.InlineVoiceHandler(bot), th.AnyInlineQuery())
	bh.Handle(handlers.InlineIpaHandler(bot), th.AnyInlineQuery())

	bh.Handle(handlers.StateHandler(bot), th.AnyMessage())

	defer bh.Stop()
	log.Println("Bot is running...")
	_ = bh.Start()
}
