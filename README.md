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

1. Uninstall any older version of Go that may be on your computer. [Click here](https://golang.org/doc/install#uninstall) for a guide on uninstalling Go.
2. [Download Go](https://golang.org/dl/). Shuttle Tracker targets Go version 1.7 and newer, but we recommend using the latest stable release of Go.
![Step 2](https://user-images.githubusercontent.com/6248819/47017483-5e891180-d120-11e8-9157-6c5823ecb13a.png)
3. Open your System Properties by searching `Edit the system environment variables` then press `Environment Variables...`
 * Ensure your `GOPATH` variable is set correctly in the `User variables for (Username)`.
 * Select `Path` under `User variables for (Username)` and make sure `%GOPATH%\bin` is on the list.
 * Make sure `GOROOT` is set correctly under `System variables`. 
 * Select `Path` under `System variables` Make sure `GOROOT\bin` is on the list.
 
 This step should be done for you already.
![Step 3](https://user-images.githubusercontent.com/6248819/47713951-e4718600-dc11-11e8-8ebc-73425eea8384.png)
![Step 3b](https://user-images.githubusercontent.com/6248819/47017509-71034b00-d120-11e8-82ee-01b17afb7ec0.png)
![Step 3c](https://user-images.githubusercontent.com/6248819/47017534-811b2a80-d120-11e8-9a8a-a625b8a74e02.png)
![Step 3d](https://user-images.githubusercontent.com/22043215/47195917-18080280-d32c-11e8-95a4-fca7b5d5f634.png)
4. Open a command prompt by pressing windows + r, then type `cmd` and hit ok, or search for command prompt
![Step 4](https://user-images.githubusercontent.com/6248819/47017557-9001dd00-d120-11e8-8258-651745338d78.png)
5. Run `go get github.com/wtg/shuttletracker`
6. Install `govendor` by running `go get -u github.com/kardianos/govendor`
7. Switch to the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`)
8. Run `govendor sync`
![Steps 5-8](https://user-images.githubusercontent.com/6248819/47017579-9db76280-d120-11e8-8de5-ab5cbe11e072.png)
9. Download and run the installer PostgreSQL [from here](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads). Select the latest version. When prompted to set a password, make it something simple, as you will be using this later, for example `shuttle`.  Use default options for everything else. Write down your password.
![Step 9](https://user-images.githubusercontent.com/6248819/47017593-ac057e80-d120-11e8-8637-18307ebeaf7e.png)
![Step 9b](https://user-images.githubusercontent.com/6248819/47017613-b6c01380-d120-11e8-85ee-35442f6ea737.png)
![Step 9c](https://user-images.githubusercontent.com/6248819/47017633-c2abd580-d120-11e8-95ab-50088fb38c40.png)
10. When complete, open pgAdmin from your search bar. If you do not see pgAdmin restart your computer. In the Object Browser, open Servers, then open your PostgreSQL server. You will need to enter your password from step 8. Once this is done, right click on Databases and select New Database. Name it `shuttletracker` and hit Ok.

![Step 10](https://user-images.githubusercontent.com/6248819/47017651-d1928800-d120-11e8-849f-535b48215923.png)
![Step 10b](https://user-images.githubusercontent.com/6248819/47017670-e0793a80-d120-11e8-9c26-473a3da6ddb5.png)
![Step 10c](https://user-images.githubusercontent.com/6248819/47017687-ecfd9300-d120-11e8-8f43-606e2405f236.png)

11. Navigate to your shuttle tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`) and rename `conf.json.sample` to `conf.json`
12. Edit `conf.json` with the following, if necessary:
   * `API.MapboxAPIKey`: Necessary for creating routes through the admin interface. [Create your own token](https://www.mapbox.com/help/how-access-tokens-work/) or ask a Shuttle Tracker developer to provide you with one.
   * `Postgres.URL`: URL where Postgres is located which will be the default with postgres:password@localhost added before the database name, where password is your password from step 8. For example: `"URL": "postgres://postgres:shuttle@localhost/shuttletracker?sslmode=disable"`
![Step 12](https://user-images.githubusercontent.com/6248819/47017722-f981eb80-d120-11e8-8ad2-4919c4052dc0.png)
13. Open pgAdmin and click on your shuttletracker database once to highlight it. Now select the query button above which has the letters SQL written on it. In the text field write `INSERT INTO users (id, username) VALUES(0, 'Your RCS ID');` and make sure to fill in Your RCS ID which is the same as your email without @rpi.edu. Then hit the execute button which is a green arrow. This will add you as an admin on your local shuttle tracker.
![Step 13](https://user-images.githubusercontent.com/6248819/47017734-03a3ea00-d121-11e8-831b-92864302ed93.png)
![Step 13b](https://user-images.githubusercontent.com/6248819/47017748-0ef71580-d121-11e8-894f-f90d31f877a7.png)
14. Start the app by running `go run cmd/shuttletracker/main.go` in the command prompt in the project root directory (`$GOPATH/src/github.com/wtg/shuttletracker`)
![Step 14](https://user-images.githubusercontent.com/6248819/47017759-1cac9b00-d121-11e8-8dcf-be6df1ff09fc.png)
16. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel
17. Copy the information from [vehicles](https://shuttles.rpi.edu/vehicles), [routes](https://shuttles.rpi.edu/routes), and [stops](https://shuttles.rpi.edu/stops) into the admin panel if you want to mimic the current shuttle tracker site.

