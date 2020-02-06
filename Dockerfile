FROM golang:alpine as builder

ADD . /source/.
RUN make

FROM alpine:latest

ENV BOT_TOKEN
ENV CHAT_ID

COPY --from=builder /source/quizbot /app/quizbot
ENTRYPOINT /app/quizbot
