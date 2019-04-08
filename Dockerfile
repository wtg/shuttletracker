FROM node:10 as npmenv

# Install npm dependencies and build
WORKDIR /frontend
COPY /frontend/package.json /frontend/package-lock.json ./
RUN npm ci

ADD /frontend /frontend
# this lets us override where the built assets should expect to be found
ARG STATIC_PATH
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
