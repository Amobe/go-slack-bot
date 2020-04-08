FROM golang:1.14-alpine3.11 AS base
WORKDIR /go/src/
ADD . /go/src/go-slack-bot
RUN cd /go/src/go-slack-bot && \
    go mod download && \
    go install

FROM alpine:3.11 AS go-slack-bot
COPY --from=base /go/bin/go-slack-bot /go-slack-bot
ENTRYPOINT /go-slack-bot
