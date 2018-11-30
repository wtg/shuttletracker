FROM golang:1.11

RUN go get -u github.com/kardianos/govendor
RUN mkdir -p /go/src/github.com/wtg/shuttletracker
WORKDIR /go/src/github.com/wtg/shuttletracker
COPY vendor/vendor.json ./vendor/
RUN govendor sync
COPY . /go/src/github.com/wtg/shuttletracker
RUN go install github.com/wtg/shuttletracker/cmd/shuttletracker

# Dokku checks http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/
RUN mkdir /app
COPY CHECKS /app

EXPOSE 8080
CMD ["/go/bin/shuttletracker"]
