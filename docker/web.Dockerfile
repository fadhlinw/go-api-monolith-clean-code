FROM golang:1.23-alpine

RUN apk add --no-cache build-base \
    inotify-tools \
    git
RUN git config --global --add safe.directory /clean_web
RUN go install github.com/rubenv/sql-migrate/...@latest

COPY . /clean_web
WORKDIR /clean_web

RUN go mod download
RUN go build -o main .
RUN apk del build-base git
RUN rm -rf /var/cache/apk/*

CMD sh /clean_web/docker/run.sh