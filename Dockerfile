FROM node:8 as npmenv

ADD /frontend /frontend

WORKDIR /frontend

# Install npm dependencies and build
RUN npm install
RUN npm run build

FROM golang:1.11

RUN go get -u github.com/kardianos/govendor
RUN mkdir -p /go/src/github.com/wtg/shuttletracker
WORKDIR /go/src/github.com/wtg/shuttletracker
COPY vendor/vendor.json ./vendor/
RUN govendor sync
COPY . /go/src/github.com/wtg/shuttletracker
RUN go install github.com/wtg/shuttletracker/cmd/shuttletracker

COPY --from=npmenv /static/ /go/src/github.com/wtg/shuttletracker/static/

# Dokku checks http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/
RUN mkdir /app
COPY CHECKS /app

EXPOSE 8080
CMD ["/go/bin/shuttletracker"]
