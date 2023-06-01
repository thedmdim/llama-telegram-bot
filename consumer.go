package main

import (
	"llama-telegram-bot/queue"
	"log"
	"runtime"
	"strings"
	"time"

	llama "github.com/go-skynet/go-llama.cpp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Result struct {
	Text chan string
	Err chan error
}

var stopButton = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Stop", "/stop"),
	),
)

func ProcessQueue() {
	for {
		task, err := qu.Dequeue()
		currentTask = task
		if err == queue.ErrQueueEmpty {
			time.Sleep(time.Second * 2)
			continue
		}
		ProcessTask(task)
	}
}

func Predict(task *queue.Task) {

	callback := func(token string) bool {
		select {
		case task.Stream <- token:
			return true
		case <- task.Stop:
			return false
		}
	}

	text, err := l.Predict(
		task.WrapInRoles(),
		llama.Debug,
		llama.SetTokenCallback(callback),
		llama.SetTokens(nTokens), 
		llama.SetThreads(runtime.NumCPU()),
		llama.SetTopK(90),
		llama.SetTopP(0.86),
		llama.SetStopWords("###"),
	)
	close(task.Stream)
	task.Result <- queue.Result{text, err}

}

func ProcessTask(task *queue.Task) {

	// generated text
	var answer string

	// Start prediction
	go Predict(task)

	defer func(){
		msg := tgbotapi.NewEditMessageText(task.UserId, task.MessageId, answer)
		msg.BaseEdit.ReplyMarkup = nil
		bot.Send(msg)
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}()

	// Send first tokens
	for token := range task.Stream {
		if strings.TrimSpace(token) != "" {
			answer += token
			break
		}
	}

	// Delete previous message notification
	delete := tgbotapi.NewDeleteMessage(task.UserId, task.AnnounceId)
	bot.Send(delete)

	if answer == "" {
		answer = "Couldn't generate answer, sorry"
		return
	}

	msg := tgbotapi.NewMessage(task.UserId, answer)
	msg.BaseChat.ReplyMarkup = &stopButton
	sent, _ := bot.Send(msg)

	// Save answer message ID to stream tokens to it
	task.MessageId = sent.MessageID

	edited := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:           task.UserId,
			MessageID: task.MessageId,
			ReplyMarkup: &stopButton,
		},
	}

	// Start streaming tokens
	var counter int
	for token := range task.Stream {
		answer += token
		counter++
		if counter == 10 {
			edited.Text = answer
			bot.Send(edited)
			counter = 0
		}
	}
	if counter != 0 {
		edited.Text = answer
		bot.Send(edited)
	}

	// Send resulting text
	result := <- task.Result
	if result.Err != nil {
		log.Println(result.Err)
		return
	}
	answer = result.Text
}
