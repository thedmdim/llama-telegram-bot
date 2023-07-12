[![Docker Pulls](https://img.shields.io/docker/pulls/thedmdim/llama-telegram-bot)](https://hub.docker.com/r/thedmdim/llama-telegram-bot)
[![Docker Image Size (tag)](https://img.shields.io/docker/image-size/thedmdim/llama-telegram-bot/latest)](https://hub.docker.com/r/thedmdim/llama-telegram-bot)


# 🦙 llama-telegram-bot

## What?
It's a chatbot for Telegram utilizing genius [llama.cpp](https://github.com/ggerganov/llama.cpp). Try live instance here [@telellamabot](https://t.me/telellamabot)

## How?
[llama-telegram-bot](https://github.com/thedmdim/llama-telegram-bot) is written in Go and uses [go-llama.cpp](https://github.com/go-skynet/go-llama.cpp) which is binding to [llama.cpp](https://github.com/ggerganov/llama.cpp)

## Quick Start
Let's start! Everything is simple!

Parameters are passed as env variables. Currently there are only 5 params:

1. `MODEL_PATH=/path/to/model`
2. `TG_TOKEN=your_telegram_bot_token_here`
3. `Q_SIZE=1000` - task queue limit (optional: default 1000)
4. `N_TOKENS=1024` - tokens to predict (optional: default 1024)
5. `N_CPU=4` - number of cpu to use (optional: default max available)

### Docker Compose
Local build (Prefered)
1. `git clone https://github.com/thedmdim/llama-telegram-bot`
2. `cp .env.example .env` and edit `.env` as you need
3. `docker compose up -d`

Pull from Docker Hub
1. `git clone https://github.com/thedmdim/llama-telegram-bot`
2. `cp .env.example .env` and edit `.env` as you need
3. `docker compose -f docker-compose.hub.yml up -d`

### Build and run as binary
You need to have Go and CMake installed
1. `git clone  --recurse-submodules https://github.com/thedmdim/llama-telegram-bot`
2. `cd llama-telegram-bot && make`
4. `go build .`
5. `env TG_TOKEN=<your_telegram_bot_token> MODEL_PATH=/path/to/your/model ./llama-telegram-bot`