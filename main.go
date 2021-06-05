package main

import (
	"log"
)

func Perror(err error) {
	if err != nil {
		panic(err)
	}
}

type Participant struct {
	ChatId       int64
	Username     string
	FirstName    string
	LastName     string
	IsSubscribed bool
	CreatedAt    string
}

func main() {
	bot, err := NewBotService(Cfg)
	Perror(err)
	log.Println(bot)
}
