# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/wtg/shuttletracker)](https://goreportcard.com/report/github.com/wtg/shuttletracker)

Tracking and mapping RPI's shuttles with [Go](https://golang.org/), [Polymer Web Components](https://www.polymer-project.org/), and [MongoDB](https://www.mongodb.org/).

Check it out in action at [shuttles.rpi.edu](https://shuttles.rpi.edu).

## Setting Up

1. Install Go
2. `go get github.com/wtg/shuttletracker`
3. `govendor sync`
4. Make sure you have npm, bower, golang and mongodb installed
5. Run `bower install` inside shuttle tracking directory to install dependencies listed in bower.json
6. Rename conf.json.sample to conf.json
7. Edit conf.json with the following:
   * Data Feed: API with tracking information, this is a unique API info url that we can get data from it. Since it is private, we will only put this on our private group for now (Slacks).
   * UpdateInterval: Number of seconds between each request to the data feed
   * MongoUrl: Url where MongoDB is located
   * MongoPort: Port where MongoDB is bound (default is 27017)
9. Run `bower install` inside shuttle tracking directory to install dependencies listed in bower.json
10. Rename conf.json.sample to conf.json and edit with the following:
   * Data Feed: API with tracking information (iTrak in our case), if using the dummy server, http://localhost:8081
   * UpdateInterval: Number of seconds between each request to the data feed
   * MongoUrl: Url where MongoDB is located
   * MongoPort: Port where MongoDB is bound (default is 27017)
11. Run the app using `go run main.go` in the project root directory
12. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the admin panel

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

### Environment variables

Most keys can be overridden with environment variables. The variables are usually
named `SECTION_KEY`. For example, overriding database's Mongo URL could be done with a variable named `DATABASE_MONGOURL`.
