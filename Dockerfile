FROM alpine:3.7

RUN apk add --no-cache \
    ca-certificates \
    git \
    go \
    musl-dev \
    tzdata

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src/github.com/wtg/shuttletracker" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/github.com/wtg/shuttletracker
RUN go get -u github.com/kardianos/govendor

COPY . .

# Dokku checks http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/
RUN mkdir /app
COPY CHECKS /app

RUN govendor sync
RUN go build ./cmd/shuttletracker

CMD ["./shuttletracker"]
