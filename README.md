# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)

Tracks and maps RPI's shuttles with [Go](https://golang.org/), [Postgres](https://www.postgresql.org), [Vue.js](https://vuejs.org/), [TypeScript](https://www.typescriptlang.org), and [Leaflet](https://leafletjs.com). Check it out in action at [shuttles.rpi.edu](https://shuttles.rpi.edu).

Looking to contribute? Please review our [Contribution Guidelines](CONTRIBUTING.md).

More project documentation can be found on the [Repository Wiki](https://github.com/wtg/shuttletracker/wiki).

## Setting up (Linux / WSL)

1. [Install Go](https://golang.org/doc/install). Shuttle Tracker requires Go 1.11 or newer, and we recommend using the latest stable Go release.
2. Clone the repository to your computer. This can be done with `git clone https://github.com/wtg/shuttletracker`.
3. Ensure you have [Postgres downloaded](https://www.postgresql.org/download/), installed.
    - This can be done with `sudo apt-get install postgresql`
4. Create a postgres databased titled "shuttletracker".
    - This can be done by running `createdb shuttletracker`. If there is an error `createdb: could not connect to database postgres: FATAL: role "<username>" does not exist`, it is due that adminstrator has not created a PostgreSQL user account for you. It can be fixed by `psql -U postgres` in terminal, `CREATE USER <username>;`, `ALTER USER <username> SUPERUSER CREATEDB;` in `postgres=#`. And you can confirm the success by running `\du` in `postgres=#`.
    - Alternatively, download the Postgres management software [pgAdmin](https://www.pgadmin.org/download/). Click on the Servers drop down, then the Databases drop down, and add a database titled "shuttletracker" by right-clicking on Databases and selecting create.
5. Switch to the Shuttle Tracker directory (`cd shuttletracker`)
9. Rename `conf.json.sample` to `conf.json`
10. Edit `conf.json` with the following, if necessary:
    - `Postgres.URL`: URL where Postgres is located. This will generally look like `postgres://user:password@localhost/shuttletracker?sslmode=disable` where `user` and `password` are replaced
11. Install Node.js and npm
    - This can be done with `sudo apt-get install nodejs`.
    - Make sure you are on version 10.16.0 or higher by running `nodejs -v`.
12. Switch to the `./frontend` directory.
13. Run `npm install`
14. Build the frontend using `npx vue-cli-service build --mode development`
    - _Note: if you are working on the frontend, you may instead use `npx vue-cli-service build --mode development --watch` in another terminal to continuously watch for changes and rebuild._
15. Go back up to the project root directory and build Shuttle Tracker by running `go build -o shuttletracker ./cmd/shuttletracker`
16. Start the app by running `./shuttletracker`
17. Add yourself as an administrator by using `./shuttletracker admins --add RCS_ID`, replacing `RCS_ID` with your RCS ID. See the "Administrators" section below for more information.
18. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel.

## Setting up (macOS)

1. [Install Go](https://golang.org/doc/install). Shuttle Tracker requires Go 1.11 or newer, and we recommend using the latest stable Go release.
2. Clone the repository to your computer. This can be done with `git clone git@github.com:wtg/shuttletracker.git`. If you receive a "permission denied" error, ensure you have [added your SSH key to your GitHub account](https://help.github.com/articles/connecting-to-github-with-ssh/).
3. If you are already familiar with Postgres and prefer to work from Terminal, run `brew install postgresql` to install Postgres to your machine.  Then run `createdb shuttletracker` to create a Postgres database.  Then skip to step 7.
4. If you are unfamiliar with Postgres and prefer to work with a graphical interface, download [Postgres.app](https://postgresapp.com) and [Postico](https://eggerapps.at/postico/).  Postgres.app allows a Postgres server to be started using a graphical interface, while Postico allows a PostgreSQL database to be managed using a graphical interface.
5. Open Postgres.app and create a new Postgres server by pressing the + in the sidebar.  Name the server "shuttletracker" and specify its port as 5432.  Then press start to run the server.
6. Open Postico and create a new database named "shuttletracker" by pressing "+ Database".
7. Switch to the Shuttle Tracker directory (`cd shuttletracker`)
8. Rename `conf.json.sample` to `conf.json`
9. Edit `conf.json` with the following, if necessary:
    - `Postgres.URL`: URL where Postgres is located. This will generally look like `postgres://user:password@localhost/shuttletracker?sslmode=disable` where `user` and `password` are replaced
10. Install [Node.js and npm](https://nodejs.org/en/download/). Be sure to download the latest LTS version. Do not download the current version or errors will occur when building the project
11. Switch to the `./frontend` directory.
12. Run `npm install`
13. Build the frontend using `npx vue-cli-service build --mode development`
    - _Note: if you are working on the frontend, you may instead use `npx vue-cli-service build --mode development --watch` in another terminal to continuously watch for changes and rebuild._
14. Go back up to the project root directory and build Shuttle Tracker by running `go build -o shuttletracker ./cmd/shuttletracker`
15. Start the app by running `./shuttletracker`
16. Add yourself as an administrator by using `./shuttletracker admins --add RCS_ID`, replacing `RCS_ID` with your RCS ID. See the "Administrators" section below for more information.
17. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel.

## Configuration

Shuttle Tracker needs configuration to run properly. The preferred method during development is to create a `conf.json` file. See `conf.json.sample` for an example of what it should contain.

`Updater.DataFeed`: API with tracking information from iTrak. For RPI, this is a unique API URL that we can get data from. It's private, and a Shuttle Tracker developer can provide it to you if necessary. However, by default, Shuttle Tracker will reach out to the instance running at shuttles.rpi.edu to piggyback off of its data feed. This means that most developers will not have to configure this key.

### Environment variables

Most keys can be overridden with environment variables. The variables names usually take the format `PACKAGE_KEY`. For example, overriding the iTRAK updater's update interval could be done with a variable named `UPDATER_UPDATEINTERVAL`.

#### Database URL

The database URL is a special case. Following the above convention, it can be set with `POSTGRES_URL`. However, for ease of deployment on Dokku, it can also be set with `DATABASE_URL`.

## Administrators

The admin interface (at `/admin`) is only accessible to users who have been added as administrators. There is a command-line utility to do this: `shuttletracker admins`. It has two flags: `--add RCS_ID` and `--remove RCS_ID`. Replace `RCS_ID` with a valid RCS ID.

### Example usage

```
> ./shuttletracker admins
No Shuttle Tracker administrators.
> ./shuttletracker admins --add naraya5
Added naraya5.
> ./shuttletracker admins --add lazare2
Added lazare2.
> ./shuttletracker admins
naraya5
lazare2
> ./shuttletracker admins --remove lazare2
Removed lazare2.
> ./shuttletracker admins
naraya5
```

## Setting up (Windows)

1. [Download Go](https://golang.org/dl/). Shuttle Tracker targets Go version 1.11 and newer, but we recommend using the latest stable release of Go.  
2. Open your System Properties by searching `Edit the system environment variables` then press `Environment Variables...`.  
 * Ensure your `GOPATH` variable is set correctly in the `User variables for (Username)`.  
 * Select `Path` under `User variables for (Username)` and make sure `%GOPATH%\bin` is on the list.  
 * Make sure `GOROOT` is set correctly under `System variables`.   
 * Select `Path` under `System variables` Make sure `GOROOT\bin` is on the list.  

 This step should be done for you already.  

<p style="text-align:center;"><img src="https://user-images.githubusercontent.com/6248819/47713951-e4718600-dc11-11e8-8ebc-73425eea8384.png" alt="3a" width="70%"></p>  
<p style="text-align:center;"><img src="https://user-images.githubusercontent.com/6248819/47017534-811b2a80-d120-11e8-9a8a-a625b8a74e02.png" alt="3b" width="70%"></p>  
<p style="text-align:center;"><img src="https://user-images.githubusercontent.com/22043215/47195917-18080280-d32c-11e8-95a4-fca7b5d5f634.png" alt="3c" width="70%"></p>  

3. Open a command prompt by pressing windows + r, then type `cmd` and hit ok, or search for command prompt.
4. Run `go get github.com/wtg/shuttletracker`.   
5. Switch to the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`).  
6. Download and run the installer PostgreSQL [from here](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads). Select the latest version. When prompted to set a password, make it something simple, as you will be using this later, for example `shuttle`.  Use default options for everything else. Remember your password.  
7. When complete, open pgAdmin from your search bar. If you do not see pgAdmin restart your computer. In the Object Browser, open Servers, then open your PostgreSQL server. You will need to enter your password from step 8. Once this is done, right click on Databases and select New Database. Name it `shuttletracker` and hit Ok.  

<p style="text-align:center;"><img src="https://user-images.githubusercontent.com/6248819/47017651-d1928800-d120-11e8-849f-535b48215923.png" alt="11" width="60%"></p>  
<p style="text-align:center;"><img src="https://user-images.githubusercontent.com/6248819/47017670-e0793a80-d120-11e8-9c26-473a3da6ddb5.png" alt="11" width="60%"></p>  

8. Navigate to your shuttle tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`) and rename `conf.json.sample` to `conf.json`  
9. Edit `conf.json` with the following, if necessary:   
   * `Postgres.URL`: URL where Postgres is located which will be the default with postgres:password@localhost added before the database name, where password is your password from step 8. For example: `"URL": "postgres://postgres:shuttle@localhost/shuttletracker?sslmode=disable"`.  

<p style="text-align:center;"><img src="https://user-images.githubusercontent.com/6248819/47017722-f981eb80-d120-11e8-8ad2-4919c4052dc0.png" alt="11" width="70%"></p>  

10. Install [Node.js and npm](https://nodejs.org/en/). Download the recommended version and install. Restart your command prompt.  
11. Navigate to the frontend directory within the shuttle tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker/frontend`).  
12. Run `npm install`   
13. Build the frontend using `npx vue-cli-service build --mode development`
    - _Note: if you are working on the frontend, you may instead use `npx vue-cli-service build --mode development --watch` in another terminal to continuously watch for changes and rebuild._

<p style="text-align:center;"><img src="https://user-images.githubusercontent.com/6248819/48438998-432b1980-e753-11e8-871c-3e6639e6383f.png" alt="11" width="70%"></p>   

14. Go back up to the project root directory (using `cd ..`) and build Shuttle Tracker by running `go build -o shuttletracker.exe cmd/shuttletracker/main.go`  
15. Start the app by running `shuttletracker.exe` in the project root directory.  
16. Add yourself as an administrator by using `shuttletracker.exe admins --add RCS_ID`, replacing `RCS_ID` with your RCS ID. See the "Administrators" section above for more information.  
17. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel.  
18. Copy the information from [vehicles](https://shuttles.rpi.edu/vehicles), [routes](https://shuttles.rpi.edu/routes), and [stops](https://shuttles.rpi.edu/stops) into the admin panel if you want to mimic the current shuttle tracker site.  