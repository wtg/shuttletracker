FROM node:8 as npmenv

ADD /frontend /frontend
WORKDIR /frontend

# Install npm dependencies and build
RUN npm install
RUN npm run build


FROM golang:1.12

RUN groupadd -r shuttletracker && useradd --no-log-init -r -g shuttletracker shuttletracker

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/shuttletracker

COPY --from=npmenv /static/ /app/static/

# Dokku checks http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/
COPY CHECKS /app

USER shuttletracker:shuttletracker
EXPOSE 8080
CMD ["/app/shuttletracker"]
