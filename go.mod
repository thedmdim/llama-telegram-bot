module llama-telegram-bot

go 1.20

require (
	github.com/go-skynet/go-llama.cpp v0.0.0-20230514113756-e9533336d695
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
)

replace github.com/go-skynet/go-llama.cpp => ./go-llama.cpp
