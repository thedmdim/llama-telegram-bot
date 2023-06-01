package main

import (
	"fmt"
	"llama-telegram-bot/queue"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func ProcessUpdate(update tgbotapi.Update) {
	// If we've gotten a message update.
	if update.Message != nil {

		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           update.Message.Chat.ID,
				ReplyToMessageID: 0,
			},
			DisableWebPagePreview: false,
		}

		if update.Message.Text == "/start" {
			msg.Text = "Just ask question"
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
			return
		}

		// Do enqueue task
		task := queue.Task{
			UserId: update.Message.From.ID,
			MessageId: update.Message.MessageID,
			Question: "### User: answer my next question. " + update.Message.Text + "\n### Assistant:",

			Stop: make(chan bool),
			Stream: make(chan string),
			Result: make(chan queue.Result),
		}
		
		n, err := qu.Enqueue(&task)
		if err != nil {
			if err == queue.ErrOnePerUser {
				msg.Text = "You've already asked your question. You can edit the existing one until it's your turn"
			}
			if err == queue.ErrQueueLimit {
				msg.Text = fmt.Sprintf("Now queue is full %d/%d. Wait one slot to be free at least.\nCheck queue /stats", n, qu.Limit)
			}
		}
		msg.Text = fmt.Sprintf("Your qustion registered! Your queue is %d/%d.\nYou can edit your message until it's your turn", n, qu.Limit)
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		task.AnnounceId = sent.MessageID
	}

	if update.EditedMessage != nil {
		task := queue.Task{
			UserId: update.Message.From.ID,
			MessageId: update.Message.MessageID,
			Question: update.Message.Text,
		}
		qu.Enqueue(&task)
	}

	
	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Stopping")
		_, err := bot.Request(callback)
		if err != nil {
			log.Println(err)
		}

		if update.CallbackQuery.Data == "/stop" && currentTask != nil {
			currentTask.Stop <- true
		}
	}
}