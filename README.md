# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/wtg/shuttletracker)](https://goreportcard.com/report/github.com/wtg/shuttletracker)

Tracking and mapping RPI's shuttles with [Go](https://golang.org/), [Vue.js](https://vuejs.org/), and [Postgres](https://www.postgresql.org).

Check it out in action at [shuttles.rpi.edu](https://shuttles.rpi.edu).

## Setting Up

1. [Install Go](https://golang.org/doc/install). Shuttle Tracker targets Go 1.7 and newer, but we recommend using the latest Go stable release.
2. Ensure your `$GOPATH` is set correctly, and `$GOPATH/bin` is in your `$PATH`
3. Run `go get github.com/wtg/shuttletracker`
4. Install `govendor` by running `go get -u github.com/kardianos/govendor`
5. Switch to the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`)
6. Run `govendor sync`
7. Ensure you have [Postgres downloaded and installed](https://www.postgresql.org/download/). On macOS, prefer installing it with Homebrew.
8. Rename `conf.json.sample` to `conf.json`
9. Edit conf.json with the following:
   * `DataFeed`: API with tracking information from iTrak... For RPI, this is a unique API URL that we can get data from. It's currently private, and we will only share it with authorized members for now.
   * `UpdateInterval`: Number of seconds between each request to the data feed
   * `Postgres.URL`: URL where Postgres is located
10. Start Postgres.
11. Add data to your database. Example DBs are provided in `example_database`, as well as a simple import/export script to setup the database for you.
    - If using an example database, you might need to check the name of the imported database and change `PostgresUrl` accordingly.
12. Start the app by running `go run cmd/shuttletracker/main.go` in the project root directory.
13. You can optionally add yourself as an administrator by using the `make-admin.sh` script in the example_database folder, passing it your RCS ID as the first argument.
14. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel

## Configuration

Shuttle Tracker reads from a `conf.json` file. It can look like this:

```
{
  "Updater": {
    "DataFeed": "",
    "UpdateInterval": "3s"
  },
  "API": {
    "CasURL": "https://cas-auth.rpi.edu/cas/",
    "Authenticate": true,
    "ListenURL": "127.0.0.1:8080"
  },
  "Postgres": {
    "URL": "127.0.0.1:27017"
  },
  "Log": {
    "Level": "debug"
  }
}
```

### Environment Variables

Most keys can be overridden with environment variables. The variables names usually take the format `SECTION_KEY`. For example, overriding database's Mongo URL could be done with a variable named `POSTGRES_URL`.
