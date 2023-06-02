FROM gcc as c_builder
RUN apt-get update && \
    apt-get install -y \
      libboost-dev libboost-program-options-dev \
      libgtest-dev \
      cmake
WORKDIR /usr/src/app
COPY . ./
RUN make

FROM golang:1.20 as go_builder
ENV GOOS=linux
WORKDIR /usr/src/app
COPY --from=c_builder /usr/src/app .
RUN go mod tidy
RUN go build -v -o /usr/local/bin/app .


FROM alpine
WORKDIR /usr/local/bin/app
COPY --from=go_builder /usr/local/bin/app .
RUN chmod +x /usr/local/bin/app
CMD ["app"]