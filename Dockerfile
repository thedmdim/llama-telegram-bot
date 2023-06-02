FROM gcc as builder
RUN apt-get update && \
    apt-get install -y \
      libboost-dev libboost-program-options-dev \
      libgtest-dev \
      cmake golang-go
ENV GOOS=linux
WORKDIR /usr/src/app
COPY . ./
RUN make && go build -v -o /usr/local/bin/app .

FROM alpine
WORKDIR /usr/local/bin/app
COPY --from=builder /usr/local/bin/app .
RUN chmod +x /usr/local/bin/app
CMD ["app"]