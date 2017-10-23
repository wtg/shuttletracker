# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/wtg/shuttletracker)](https://goreportcard.com/report/github.com/wtg/shuttletracker)

Tracking and mapping RPI's shuttles with [Go](https://golang.org/), [Polymer Web Components](https://www.polymer-project.org/), and [MongoDB](https://www.mongodb.org/).

Check it out in action at [shuttles.rpi.edu](https://shuttles.rpi.edu).

## Setting Up

1. Install Go (https://golang.org/doc/install)
2. Ensure your `$GOPATH` is set correctly, and is apart of your `$PATH`
3. Run `go get github.com/wtg/shuttletracker`
4. Install `govendor`  by running `go get -u github.com/kardianos/govendor`
5. Switch to the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`)
6. Run `govendor sync`
7. Ensure you have NPM, Bower, and MongoDB installed.
8. Run `bower install` inside the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`) to install dependencies listed in bower.json
9. Rename `conf.json.sample` to `conf.json`
10. Edit conf.json with the following:
   * `DataFeed`: API with tracking information from iTrak... For RPI, this is a unique API URL that we can get data from. It's currently private, and we will only share it with authorized members for now.
   * `UpdateInterval`: Number of seconds between each request to the data feed
   * `MongoUrl`: URL where MongoDB is located
   * `MongoPort`: Port where MongoDB is bound (default is 27017)
11. Start MongoDB, and ensure it is running, and listening on port 27017 (or whichever port you defined in `MongoPort` within `conf.json`)
12. Add data to your database. Example DBs are provided in `example_database`, as well as a simple import/export script to setup the database for you.
    - If using an example database, you might need to check the name of the imported database, and change `MongoUrl` accordingly.
13. Start the app by running `go run main.go` in the project root directory.
14. You can optionally add yourself as an administrator by using the `make-admin` script in the example_database folder, passing it your RCS ID as the first argument.
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
    "GoogleMapAPIKey": "",
    "GoogleMapMinDistance": 1,
    "CasURL": "https://cas-auth.rpi.edu/cas/",
    "Authenticate": true,
    "ListenURL": "127.0.0.1:8080"
  },
  "Database": {
    "MongoURL": "127.0.0.1:27017"
  },
  "Log": {
    "Level": "debug"
  }
}
```

### Environment Variables

Most keys can be overridden with environment variables. The variables names usually take the format `SECTION_KEY`. For example, overriding database's Mongo URL could be done with a variable named `DATABASE_MONGOURL`.
