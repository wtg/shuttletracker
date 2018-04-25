FROM golang:1.10

RUN go get -u github.com/kardianos/govendor
RUN mkdir -p /go/src/github.com/wtg/shuttletracker
WORKDIR /go/src/github.com/wtg/shuttletracker
COPY vendor/vendor.json ./vendor/
RUN govendor sync
COPY . /go/src/github.com/wtg/shuttletracker
RUN go install github.com/wtg/shuttletracker

# Dokku checks http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/
RUN mkdir /app
COPY CHECKS /app

CMD ["/go/bin/shuttletracker"]
