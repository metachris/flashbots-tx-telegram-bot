package main

import (
	"fmt"
	"log"
	"os"
	"time"

	fbfailedtx "github.com/metachris/flashbots-failed-tx"
)

func Perror(err error) {
	if err != nil {
		panic(err)
	}
}

func runApiRequestsForFailedTx(channel chan<- fbfailedtx.BlockWithFailedTx) {
	isBlockAlreadyProcessed := make(map[int64]bool)
	isFirstRequest := true // first request is just caching all previous failed tx blocks, but not sending to telegram

	for {
		time.Sleep(5 * time.Second)

		x, e := GetFailedFlashbotsTransactions()
		if e != nil {
			log.Println(e)
		}

		for _, b := range x {
			if isBlockAlreadyProcessed[b.BlockHeight] {
				continue
			}

			if !isFirstRequest { // don't send to channel
				// log.Println("sending to channel", b.BlockHeight, len(b.FailedTx))
				channel <- b
			}

			isBlockAlreadyProcessed[b.BlockHeight] = true
		}

		isFirstRequest = false
	}
}

func main() {
	log.SetOutput(os.Stdout)

	// Start API request routine to get latest blocks with failed tx
	c := make(chan fbfailedtx.BlockWithFailedTx)
	go runApiRequestsForFailedTx(c)

	// Start bot
	log.Println("Starting bot...")
	bot, err := NewBotService(Cfg)
	Perror(err)
	log.Println(bot)

	// Start watching
	log.Println("Waiting for updates...")
	for {
		select {
		case update := <-bot.UpdateChan:
			bot.HandleUpdate(update)
		case failedTxBlock := <-c:
			log.Println("Failed tx", failedTxBlock.BlockHeight, len(failedTxBlock.FailedTx))
			if len(failedTxBlock.FailedTx) == 1 {
				msg := makeMsgForTx(failedTxBlock.FailedTx[0])
				bot.SendToSubscribers(msg)
			} else if len(failedTxBlock.FailedTx) > 1 {
				msg := fmt.Sprintf("block [%d](https://etherscan.io/block/%d) has %d failed tx:\n", failedTxBlock.BlockHeight, failedTxBlock.BlockHeight, len(failedTxBlock.FailedTx))
				for _, tx := range failedTxBlock.FailedTx {
					msg = fmt.Sprintf("%s%s", msg, makeMsgForTx(tx))
				}
				bot.SendToSubscribers(msg)
			}
		}
	}
}

func makeMsgForTx(tx fbfailedtx.FailedTx) string {
	if tx.IsFlashbots {
		return fmt.Sprintf("failed Flashbots tx [%s](https://etherscan.io/tx/%s) from [%s](https://etherscan.io/address/%s) in block [%d](https://etherscan.io/block/%d)\n", tx.Hash, tx.Hash, tx.From, tx.From, tx.Block, tx.Block)
	} else {
		return fmt.Sprintf("failed 0-gas tx [%s](https://etherscan.io/tx/%s) from [%s](https://etherscan.io/address/%s) in block [%d](https://etherscan.io/block/%d)\n", tx.Hash, tx.Hash, tx.From, tx.From, tx.Block, tx.Block)
	}
}
