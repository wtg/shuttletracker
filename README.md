# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/wtg/shuttletracker)](https://goreportcard.com/report/github.com/wtg/shuttletracker)

Tracking and mapping RPI's shuttles with [Go](https://golang.org/), [Polymer Web Components](https://www.polymer-project.org/), and [MongoDB](https://www.mongodb.org/).

Check it out in action at [shuttles.rpi.edu](https://shuttles.rpi.edu).

## Setting Up

1. Install Golang (https://golang.org/doc/install)
2. Run `go get github.com/wtg/shuttletracker`
3. Run `govendor sync`
4. Ensure you have NPM, Bower, and MongoDB installed.
5. Run `bower install` inside the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`) to install dependencies listed in bower.json
6. Rename `conf.json.sample` to `conf.json`
7. Edit conf.json with the following:
   * `Data Feed`: API with tracking information from iTrak... For RPI, this is a unique API URL that we can get data from. It's currently private, and we will only share it with authorized members for now.
   * `UpdateInterval`: Number of seconds between each request to the data feed
   * `MongoUrl`: Url where MongoDB is located
   * `MongoPort`: Port where MongoDB is bound (default is 27017)
8. Start MongoDB, and ensure it is running, and listening on port 27017 (or whichever port you defined in `MongoPort` within `conf.json`)
9. Add data to your database. Example DBs are provided in `example_database`, as well as a simple import/export script to setup the database for you.
10. Start the app by running `go run main.go` in the project root directory.
11. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel

## Configuration

Shuttle Tracker reads from a `conf.json` file. It can look like this:

```
{
  "Updater": {
    "DataFeed": "",
    "UpdateInterval": "3s"
  },
  "API": {
    "GoogleMapAPIKey": "",
    "GoogleMapMinDistance": 1,
    "CasURL": "https://cas-auth.rpi.edu/cas/",
    "Authenticate": true
  },
  "Database": {
    "MongoURL": "localhost:27017"
  },
  "Log": {
    "Level": "debug"
  },
  "Server": {
    "ListenURL": "localhost:8080"
  }
}
```

### Environment Variables

Most keys can be overridden with environment variables. The variables names usually take the format `SECTION_KEY`. For example, overriding database's Mongo URL could be done with a variable named `DATABASE_MONGOURL`.
