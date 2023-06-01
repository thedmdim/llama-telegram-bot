FROM golang:1.20 as builder
ENV GOOS=linux
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY . ./
RUN go mod tidy
COPY . .
RUN make && go build -v -o /usr/local/bin/app .


FROM alpine
WORKDIR /usr/local/bin/app
COPY --from=builder /usr/local/bin/app .
RUN chmod +x /usr/local/bin/app
CMD ["app"]