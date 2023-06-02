FROM golang:alpine as builder
RUN apk add --no-cache gcc musl-dev make cmake
ENV GOOS=linux
WORKDIR /usr/src/app
COPY . ./
RUN make && go build -v -o /usr/local/bin/app .

FROM alpine
WORKDIR /usr/local/bin/app
COPY --from=builder /usr/local/bin/app .
RUN chmod +x /usr/local/bin/app
CMD ["app"]