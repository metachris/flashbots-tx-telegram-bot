package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotService struct {
	Database     *DbService
	Bot          *tgbotapi.BotAPI
	Participants map[int64]*Participant
	UpdateChan   tgbotapi.UpdatesChannel
}

func NewBotService(cfg Config) (bot *BotService, err error) {
	bot = &BotService{
		Database:     NewDbService(cfg.Database),
		Participants: make(map[int64]*Participant),
	}

	// Get participants from DB
	participantsArr, err := bot.Database.GetParticipants()
	Perror(err)
	log.Println(len(participantsArr), "participants in DB")
	for _, p := range participantsArr {
		log.Println(p)
		bot.Participants[p.ChatId] = p
	}

	bot.Bot, err = tgbotapi.NewBotAPI(Cfg.TelegramApiKey)
	Perror(err)

	bot.Bot.Debug = Cfg.TelegramBotDebug
	log.Printf("Bot authorized on account %s", bot.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bot.UpdateChan, err = bot.Bot.GetUpdatesChan(u)
	Perror(err)

	return bot, nil
}

func (b *BotService) GetSubscribers() (ret []*Participant) {
	for _, v := range b.Participants {
		if v.IsSubscribed {
			ret = append(ret, v)
		}
	}
	return ret
}

func (b *BotService) SendToSubscribers(msg string) {
	subs := b.GetSubscribers()
	if len(subs) == 0 || b.Bot == nil {
		return
	}

	for _, v := range subs {
		if v.IsSubscribed {
			botMsg := tgbotapi.NewMessage(v.ChatId, msg)
			b.Bot.Send(botMsg)
		}
	}
}

func (b *BotService) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	msg := update.Message.Text
	log.Printf("[%s] %s", update.Message.From.UserName, msg)
	if msg == "/start" {
		entry, exists := b.Participants[update.Message.Chat.ID]
		if !exists {
			entry = &Participant{
				ChatId:   update.Message.Chat.ID,
				Username: update.Message.From.UserName,
			}
			b.Participants[update.Message.Chat.ID] = entry
		}
		entry.IsSubscribed = true
		b.Database.UpdateParticipant(entry)
		log.Println("Subscribed", entry)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are now subscribed to failed Flashbots transactions")
		msg.ReplyToMessageID = update.Message.MessageID
		b.Bot.Send(msg)

	} else if msg == "/stop" {
		entry, exists := b.Participants[update.Message.Chat.ID]
		if exists {
			entry.IsSubscribed = false
			b.Database.UpdateParticipant(entry)

			log.Println("Unsubscribed", entry)
			// log.Println(GetSubs())

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are now unsubscribed")
			msg.ReplyToMessageID = update.Message.MessageID
			b.Bot.Send(msg)
		}

	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Not understood. You can write /start to subscribe, or /stop to unsubscribe to failed Flashbots transactions.")
		msg.ReplyToMessageID = update.Message.MessageID
		b.Bot.Send(msg)
	}
}
