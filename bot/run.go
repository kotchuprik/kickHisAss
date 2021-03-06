package bot

import (
	"log"
	"net/http"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Config for bot for starting
type Configuration struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	BotToken string `json:"bot_token,omitempty"`
	WebHook  string `json:"web_hook,omitempty"`
}

// Run start bot
func Run(cfg *Configuration) error {

	// Start botAPI with token
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return err
	}

	// Set WebHook bots
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(cfg.WebHook + cfg.BotToken))
	if err != nil {
		return err
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(cfg.Host+":"+cfg.Port, nil)

	log.Printf("Starting, bot on: http://%s\n", cfg.Host+":"+cfg.Port)

	for update := range updates {

		if ok := checkWords(update.Message.Text); ok {

			// delete message
			if api, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			}); err != nil {
				log.Printf("Err delete message: %v\n", api.Result)
			}

			var f = false

			// restrict user
			if api, err := bot.RestrictChatMember(tgbotapi.RestrictChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: update.Message.Chat.ID,
					UserID: update.Message.From.ID,
				},
				CanSendMessages: &f,
				CanSendMediaMessages: &f,
				CanSendOtherMessages: &f,
				CanAddWebPagePreviews: &f,
				UntilDate: int64(time.Duration(5) * time.Minute),
			}); err != nil {
				log.Printf("Err restrict user: %v\n", api.Result)
			}

			// kick user
			if api, err := bot.KickChatMember(tgbotapi.KickChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: update.Message.Chat.ID,
					UserID: update.Message.From.ID,
				},
				UntilDate: int64(time.Duration(5) * time.Minute),
			}); err != nil {
				log.Printf("Err kick user: %v\n", api.Result)
			}

		}
	}
	return nil
}
