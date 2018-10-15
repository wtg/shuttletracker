# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/wtg/shuttletracker)](https://goreportcard.com/report/github.com/wtg/shuttletracker)

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
   * `API.MapboxAPIKey`: Necessary for creating routes through the admin interface. [Create your own token](https://www.mapbox.com/help/how-access-tokens-work/) or ask a Shuttle Tracker developer to provide you with one.
   * `Postgres.URL`: URL where Postgres is located. The provided default typically won't need to be modified.

11. Add data to your database. Example DBs are provided in `example_database`, as well as a simple import/export script to setup the database for you.
    - If using an example database, you might need to check the name of the imported database and change the Postgres URL accordingly.
12. Start the app by running `go run cmd/shuttletracker/main.go` in the project root directory.
13. You can optionally add yourself as an administrator by using the `make-admin.sh` script in the example_database folder, passing it your RCS ID as the first argument.
14. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel

## Configuration

Shuttle Tracker needs configuration to run properly. The preferred method during development is to create a `conf.json` file. See `conf.json.sample` for an example of what it should contain.

`Updater.DataFeed`: API with tracking information from iTrak. For RPI, this is a unique API URL that we can get data from. It's private, and a Shuttle Tracker developer can provide it to you if necessary. However, by default, Shuttle Tracker will reach out to the instance running at shuttles.rpi.edu to piggyback off of its data feed. This means that most developers will not have to configure this key.

### Environment variables

Most keys can be overridden with environment variables. The variables names usually take the format `SECTION_KEY`. For example, overriding database's Mongo URL could be done with a variable named `POSTGRES_URL`.


## Setting up (Windows)

1. [Download Go](https://golang.org/dl/). Shuttle Tracker targets Go version 1.7 and newer, but we recommend using the latest stable release of Go.
![Step 1](https://i.gyazo.com/9287b1920748f974cf15ee9f4222f3a1.png)
2. Ensure your `$GOPATH` is set correctly, and `$GOPATH/bin` is in your `$PATH`, this should be done for you already.
![Step 2](https://i.gyazo.com/e178a37f2ff99ac937dd55e85038553a.png)
![Step 2b](https://i.gyazo.com/fdd1f602e929de3943414034003ff6ba.png)
3. Open a command prompt by pressing windows + r, then type `cmd` and hit ok, or search for command prompt
![Step 3](https://i.gyazo.com/8534b060f86f8888ef77e98da6d03790.png)
4. Run `go get github.com/wtg/shuttletracker`
5. Install `govendor` by running `go get -u github.com/kardianos/govendor`
6. Switch to the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`)
7. Run `govendor sync`
![Steps 4-6](https://i.gyazo.com/04ff6e707cab9bcc51faf97cfa6c89b0.png)
8. Download and run the installer PostgreSQL [from here](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads). Select the latest version. When prompted to set a password, make it something simple, as you will be using this later, for example `shuttle`.  Use default options for everything else. Write down your password.
![Step 8](https://i.gyazo.com/c4d522d1aea67aedfc74160a490ab788.png)
![Step 8b](https://i.gyazo.com/0c59be022c6a67305f54e42fd0e507c5.png)
![Step 8c](https://i.gyazo.com/ff6cb14302cd39246aafea2eb33b30b3.png)
9. When complete, open pgAdmin from your search bar. In the Object Browser, open Servers, then open your PostgreSQL server. You will need to enter your password from step 8. Once this is done, right click on Databases and select New Database. Name it `shuttletracker` and hit Ok.

![Step 9](https://i.gyazo.com/f369818d8dd146fd459c1841280b3a20.png)
![Step 9b](https://i.gyazo.com/1be05d25b102d2bfd8b05ae18e02758c.png)
![Step 9c](https://i.gyazo.com/2fd7541f8c504963891ee5361298dc49.png)

10. Navigate to your shuttle tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`) and rename `conf.json.sample` to `conf.json`
11. Edit `conf.json` with the following, if necessary:
   * `API.MapboxAPIKey`: Necessary for creating routes through the admin interface. [Create your own token](https://www.mapbox.com/help/how-access-tokens-work/) or ask a Shuttle Tracker developer to provide you with one.
   * `Postgres.URL`: URL where Postgres is located which will be the default with postgres:password@localhost added before the database name, where password is your password from step 8. For example: `"URL": "postgres://postgres:shuttle@localhost/shuttletracker?sslmode=disable"`
![Step 11](https://i.gyazo.com/8520e284349c82830b57508c40de08e5.png)
12. Open pgAdmin and click on your shuttletracker database once to highlight it. Now select the query button above which has the letters SQL written on it. In the text field write `INSERT INTO users (id, username) VALUES(0, 'Your RCS ID');` and make sure to fill in Your RCS ID which is the same as your email without @rpi.edu. Then hit the execute button which is a green arrow. This will add you as an admin on your local shuttle tracker.
![Step 12](https://i.gyazo.com/052cf283931cfac08216b2261c724895.png)
![Step 12b](https://i.gyazo.com/5d6fc2f1e4fd145664d443a2f025327f.png)
13. Start the app by running `go run cmd/shuttletracker/main.go` in the command prompt in the project root directory (`$GOPATH/src/github.com/wtg/shuttletracker`)
![Step 13](https://i.gyazo.com/40efd8043ff37471c6fd59b98bf2d569.png)
15. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel
16. Copy the information from [vehicles](https://shuttles.rpi.edu/vehicles), [routes](https://shuttles.rpi.edu/routes), and [stops](https://shuttles.rpi.edu/stops) into the admin panel if you want to mimic the current shuttle tracker site.

