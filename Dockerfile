FROM ubuntu:latest
#FROM alpine:3.6

# RUN apk add --no-cache \
#     ca-certificates \
#     git \
#     go \
#     musl-dev \
#     nodejs-npm \
#     tzdata
RUN apt-get update && apt-get install -y \
    git \
    golang \
    nodejs \
    npm \
    tzdata

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src/github.com/wtg/shuttletracker" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/github.com/wtg/shuttletracker
RUN go get -u github.com/kardianos/govendor

RUN npm install -g bower

RUN ln -s /usr/bin/nodejs /usr/bin/node
COPY ./bower.json .
RUN bower install --allow-root

COPY . .

# Dokku checks http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/
RUN mkdir /app
COPY CHECKS /app

RUN govendor sync
RUN go build -o shuttletracker

CMD ["./shuttletracker"]
