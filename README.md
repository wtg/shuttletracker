# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)

Tracking and mapping RPI's shuttles with [Go](https://golang.org/), [Vue.js](https://vuejs.org/), and [Postgres](https://www.postgresql.org).

Check it out in action at [shuttles.rpi.edu](https://shuttles.rpi.edu).

## Setting up

1. [Install Go](https://golang.org/doc/install). Shuttle Tracker targets Go 1.7 and newer, but we recommend using the latest Go stable release.
2. Ensure your `$GOPATH` is set correctly, and `$GOPATH/bin` is in your `$PATH`
3. Run `go get github.com/wtg/shuttletracker`
4. Install `govendor` by running `go get -u github.com/kardianos/govendor`
5. Switch to the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`)
6. Run `govendor sync`
7. Ensure you have [Postgres downloaded](https://www.postgresql.org/download/), installed, and running. On macOS, prefer installing it with Homebrew.
8. Run `createdb shuttletracker` to create a Postgres database.
9. Rename `conf.json.sample` to `conf.json`
10. Edit `conf.json` with the following, if necessary:
    - `API.MapboxAPIKey`: Necessary for creating routes through the admin interface. [Create your own token](https://www.mapbox.com/help/how-access-tokens-work/) or ask a Shuttle Tracker developer to provide you with one.
    - `Postgres.URL`: URL where Postgres is located. The provided default typically won't need to be modified.
11. Install Node.js and npm.
12. Switch to the `./frontend` directory.
13. Run `npm install`
14. Build the frontend using `npx vue-cli-service build --mode development`
    - _Note: if you are working on the frontend, you may instead use `npx vue-cli-service build --mode development --watch` in another terminal to continuously watch for changes and rebuild._
15. Go back up to the project root directory and build Shuttle Tracker by running `go build -o shuttletracker cmd/shuttletracker/main.go`
16. Start the app by running `./shuttletracker` in the project root directory.
17. You can add yourself as an administrator by using `./shuttletracker admins --add RCS_ID`, replacing `RCS_ID` with your RCS ID. See the "Administrators" section below for more information.
18. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel.

## Configuration

Shuttle Tracker needs configuration to run properly. The preferred method during development is to create a `conf.json` file. See `conf.json.sample` for an example of what it should contain.

`Updater.DataFeed`: API with tracking information from iTrak. For RPI, this is a unique API URL that we can get data from. It's private, and a Shuttle Tracker developer can provide it to you if necessary. However, by default, Shuttle Tracker will reach out to the instance running at shuttles.rpi.edu to piggyback off of its data feed. This means that most developers will not have to configure this key.

### Environment variables

Most keys can be overridden with environment variables. The variables names usually take the format `SECTION_KEY`. For example, overriding database's Mongo URL could be done with a variable named `POSTGRES_URL`.

## Administrators

The admin interface (at `/admin`) is only accessible to users who have been added as administrators. There is a command-line utility to do this: `shuttletracker admins`. It has two flags: `--add RCS_ID` and `--remove RCS_ID`. Replace `RCS_ID` with a valid RCS ID.

### Example usage

```
> ./shuttletracker admins
No Shuttle Tracker administrators.
> ./shuttletracker admins --add kochms
Added kochms.
> ./shuttletracker admins --add lyonj4
Added lyonj4.
> ./shuttletracker admins
kochms
lyonj4
> ./shuttletracker admins --remove lyonj4
Removed lyonj4.
> ./shuttletracker admins
kochms
```
