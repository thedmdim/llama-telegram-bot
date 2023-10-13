package main

import (
	"log"
	"os"
	"runtime"
	"strconv"

	llama "github.com/go-skynet/go-llama.cpp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


var apiToken = os.Getenv("TG_TOKEN")
var modelPath = os.Getenv("MODEL_PATH")
var nTokens int
var nCpu int

var SingleMessagePrompt string
var ReplyMessagePrompt string
var StopWord = os.Getenv("STOP_WORD")

var l *llama.LLama
var bot *tgbotapi.BotAPI
var qu *TaskQueue
var currentTask *Task


func main() {
	var err error

	if apiToken == "" || modelPath == "" {
		log.Fatalln("Please provide TG_TOKEN and MODEL_PATH env variables")
	}

	// Init queue
	var queueSize = 1000
	if s := os.Getenv("Q_SIZE"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			queueSize = n
		}
	}
	qu = NewTaskQueue(queueSize)
	

	// N tokens
	nTokens = 1000
	if s := os.Getenv("N_TOKENS"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			nTokens = n
		}
	}

	// N cores
	nCpu = runtime.NumCPU()
	if s := os.Getenv("N_CPU"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			nCpu = n
		}
	}

	// Init Prompt templates
	SingleMessagePrompt = os.Getenv("SINGLE_MESSAGE_PROMPT")
	ReplyMessagePrompt = os.Getenv("REPLY_MESSAGE_PROMPT")
	if SingleMessagePrompt == "" {
		SingleMessagePrompt = "### User: Response to my next request. %s ### Assistant:"
	}
	if ReplyMessagePrompt == "" {
		ReplyMessagePrompt = "### Assistant: %s ### User: %s \n### Assistant:"
	}
	if StopWord == "" {
		StopWord = "###"
	}

	// Init LLAMA binding
	l, err = llama.New(modelPath, llama.SetContext(1024), llama.EnableEmbeddings, llama.EnableMLock)
	if err != nil {
		log.Fatalf("Loading the model failed: %s", err.Error())
	}

	// Init Telegram API client
    bot, err = tgbotapi.NewBotAPI(apiToken)
    if err != nil {
        log.Fatal(err)
    }

	// Start iterating through queue
	go ProcessQueue()

	// Receive updates
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)
    for update := range updates {
		ProcessUpdate(update)
    }
}
