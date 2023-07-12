package main

import (
	"llama-telegram-bot/queue"
	"log"
	"strings"
	"time"

	llama "github.com/go-skynet/go-llama.cpp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


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


type Result struct {
	Text string
	Err error
}


func Predict(task *queue.Task) (chan string, chan Result) {

	stream := make(chan string)
	result := make(chan Result)

	go func(){
		callback := func(token string) bool {
			select {
			case stream <- token:
				return true
			case <- task.Stop:
				return false
			}
		}
	
		text, err := l.Predict(
			task.Question,
			llama.Debug,
			llama.SetTokenCallback(callback),
			llama.SetTokens(nTokens), 
			llama.SetThreads(nCpu),
			llama.SetTopK(90),
			llama.SetTopP(0.86),
			llama.SetStopWords("###"),
		)
		close(stream)
		result <- Result{text, err}
	}()
	
	return stream, result
}

// This function is a mess
func ProcessTask(task *queue.Task) {

	log.Printf("Start processing task from user %d\n", task.UserID)
	log.Printf("The prompt is:\n%s\n", task.Question)

	// Start prediction
	stream, result :=  Predict(task)

	// Resulting generated text
	var answer string

	var counter int
	var issent bool
	for {
		select {
		case token := <- stream: 
			if !issent && strings.TrimSpace(token) != "" {
				answer += token
				msg := tgbotapi.NewMessage(task.UserID, answer)
				msg.ReplyMarkup = &stopButton
				sent, err := bot.Send(msg)
				if err != nil {
					log.Println("[ProcessTask] error sending answer:", err)
					continue
				}
				// Save answer message ID to stream tokens to it
				task.MessageID = sent.MessageID
				issent = true
				continue
			}

			answer += token
			counter++
			if counter == 6 {
				edited := tgbotapi.NewEditMessageText(task.UserID, task.MessageID, answer)
				edited.ReplyMarkup = &stopButton
				_, err := bot.Send(edited)
				if err != nil {
					log.Println("[ProcessTask] error streaming answer:", err)
				}
				counter = 0
				
			}

		case prediction := <- result:

			delete := tgbotapi.NewDeleteMessage(task.UserID, task.AnnounceID)
			_, err := bot.Request(delete)
			if err != nil {
				log.Println("Couldn't delete announce message:", err)
			}
			

			if prediction.Err != nil || strings.TrimSpace(prediction.Text) == "" || task.MessageID == 0 {
				log.Println("[ProcessTask] prediction error:", prediction.Err, prediction.Text)
				failure := tgbotapi.NewMessage(task.UserID, "Sorry, couldn't generate answer")
				_, err := bot.Send(failure)
				if err != nil {
					log.Println("[ProcessTask] error sending failure message:", err)
				}
				return
			}


			edited := tgbotapi.NewEditMessageText(task.UserID, task.MessageID, prediction.Text)
			// Set parse mode to Markdown if it's backticks there
			if nBackticks := strings.Count(prediction.Text, "`"); nBackticks > 0 && nBackticks % 2 == 0 {
				edited.ParseMode = "Markdown"
			}
			_, err = bot.Send(edited)
			if err != nil {
				log.Println("[ProcessTask] error sending answer:", err)
			}

			log.Printf("Generated answer is:\n%s\n", prediction.Text)

			return
		}
	}
}