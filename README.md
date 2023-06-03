# llama-telegram-bot

## What?
This project aims to liberate using of LLM providing the most simple way for use and deploy your own instance of AI chatbot.

## How?
[llama-telegram-bot](https://github.com/thedmdim/llama-telegram-bot) is written in Go and built on top of [go-llama.cpp](https://github.com/go-skynet/go-llama.cpp) which is binding to [llama.cpp](https://github.com/ggerganov/llama.cpp)

## Quick Start

Parameters to executable are passed as env variables:
```bash
MODEL_PATH=/path/to/model
TG_TOKEN=your_telegram_bot_token_here
Q_SIZE=1000 # task queue limit (optional)
N_TOKENS=1024 # tokens to predict (optional)
```
### Docker
```bash
docker run -d --privileged \
    --name llama-telegram-bot \
    -v /path/to/your/model.bin:/model.bin
    -e MODEL_PATH=/model.bin
    -e TG_TOKEN=your_telegram_api_token \
    thedmdim/llama-telegram-bot
```

Example:
```bash
docker run -t --privileged \
    -v /root/stable-vicuna-13B.ggmlv3.q8_0.bin:/model.bin \
    -e MODEL_PATH=/model.bin \
    -e TG_TOKEN=6082407582:AAОS2uRCE-nlM3tkКdofW_EBTSdVI4_OV8g \
    thedmdim/llama-telegram-bot
```
`--privileged` here used for container to be able to mlock for loading whole model into RAM. So ensure that your RAM is bigger than model you gonna use.

### Building
You need to have Go and CMake installed
1. git clone https://github.com/thedmdim/llama-telegram-bot
2. git submodule update --init --recursive
3. make
4. go build .
5. env TG_TOKEN=<your_telegram_bot_token> MODEL_PATH=/path/to/your/model llama-telegram-bot