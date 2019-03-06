# Shuttle Tracker [![Build Status](https://travis-ci.org/wtg/shuttletracker.svg?branch=master)](https://travis-ci.org/wtg/shuttletracker)&nbsp;[![codecov](https://codecov.io/gh/wtg/shuttletracker/branch/master/graph/badge.svg)](https://codecov.io/gh/wtg/shuttletracker)&nbsp;[![GoDoc](https://godoc.org/github.com/wtg/shuttletracker?status.svg)](https://godoc.org/github.com/wtg/shuttletracker)

Tracks and maps RPI's shuttles with [Go](https://golang.org/), [Postgres](https://www.postgresql.org), [Vue.js](https://vuejs.org/), [TypeScript](https://www.typescriptlang.org), and [Leaflet](https://leafletjs.com). Check it out in action at [shuttles.rpi.edu](https://shuttles.rpi.edu).

## Setting up

1. [Install Go](https://golang.org/doc/install). Shuttle Tracker requires Go 1.11 or newer, and we recommend using the latest stable Go release.
2. Clone the repository to your computer. On macOS, Linux, or WSL, this can be done with `git clone git@github.com:wtg/shuttletracker.git`. If you receive a "permission denied" error, ensure you have [added your SSH key to your GitHub account](https://help.github.com/articles/connecting-to-github-with-ssh/).
3. Ensure you have [Postgres downloaded](https://www.postgresql.org/download/), installed, and running. On macOS, Homebrew makes this easy.
4. Run `createdb shuttletracker` to create a Postgres database.
5. Switch to the Shuttle Tracker directory (`cd shuttletracker`)
9. Rename `conf.json.sample` to `conf.json`
10. Edit `conf.json` with the following, if necessary:
    - `API.MapboxAPIKey`: Necessary for creating routes through the admin interface. [Create your own token](https://www.mapbox.com/help/how-access-tokens-work/) or ask a Shuttle Tracker developer to provide you with one.
    - `Postgres.URL`: URL where Postgres is located. The provided default typically won't need to be modified.
11. Install Node.js and npm.
12. Switch to the `./frontend` directory.
13. Run `npm install`
14. Build the frontend using `npx vue-cli-service build --mode development`
    - _Note: if you are working on the frontend, you may instead use `npx vue-cli-service build --mode development --watch` in another terminal to continuously watch for changes and rebuild._
15. Go back up to the project root directory and build Shuttle Tracker by running `go build -o shuttletracker ./cmd/shuttletracker`
16. Start the app by running `./shuttletracker` in the project root directory.
17. Add yourself as an administrator by using `./shuttletracker admins --add RCS_ID`, replacing `RCS_ID` with your RCS ID. See the "Administrators" section below for more information.
18. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel.

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

## Setting up (Windows)

1. Uninstall any older version of Go that may be on your computer. [Click here](https://golang.org/doc/install#uninstall) for a guide on uninstalling Go.  
2. [Download Go](https://golang.org/dl/). Shuttle Tracker targets Go version 1.7 and newer, but we recommend using the latest stable release of Go.  

![Step 2](https://user-images.githubusercontent.com/6248819/47017483-5e891180-d120-11e8-9157-6c5823ecb13a.png)  
3. Open your System Properties by searching `Edit the system environment variables` then press `Environment Variables...`.  
 * Ensure your `GOPATH` variable is set correctly in the `User variables for (Username)`.  
 * Select `Path` under `User variables for (Username)` and make sure `%GOPATH%\bin` is on the list.  
 * Make sure `GOROOT` is set correctly under `System variables`.   
 * Select `Path` under `System variables` Make sure `GOROOT\bin` is on the list.  
 
 This step should be done for you already.  

![Step 3](https://user-images.githubusercontent.com/6248819/47713951-e4718600-dc11-11e8-8ebc-73425eea8384.png)  
![Step 3b](https://user-images.githubusercontent.com/6248819/47017509-71034b00-d120-11e8-82ee-01b17afb7ec0.png)  
![Step 3c](https://user-images.githubusercontent.com/6248819/47017534-811b2a80-d120-11e8-9a8a-a625b8a74e02.png)  
![Step 3d](https://user-images.githubusercontent.com/22043215/47195917-18080280-d32c-11e8-95a4-fca7b5d5f634.png)  
4. Open a command prompt by pressing windows + r, then type `cmd` and hit ok, or search for command prompt.  

![Step 4](https://user-images.githubusercontent.com/6248819/47017557-9001dd00-d120-11e8-8258-651745338d78.png)  
5. Run `go get github.com/wtg/shuttletracker`.  
6. Install `govendor` by running `go get -u github.com/kardianos/govendor`.  
7. Switch to the Shuttle Tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker`).  
8. Run `govendor sync`.  
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
   * `Postgres.URL`: URL where Postgres is located which will be the default with postgres:password@localhost added before the database name, where password is your password from step 8. For example: `"URL": "postgres://postgres:shuttle@localhost/shuttletracker?sslmode=disable"`.  
   
![Step 12](https://user-images.githubusercontent.com/6248819/47017722-f981eb80-d120-11e8-8ad2-4919c4052dc0.png)  

13. Install [Node.js and npm](https://nodejs.org/en/). Download the recommeneded version and install. Restart your command prompt.  
![Step 13](https://user-images.githubusercontent.com/6248819/48438877-f47d7f80-e752-11e8-9584-e5dd79ec92d2.png)  
![Step 13b](https://user-images.githubusercontent.com/6248819/48438809-d1eb6680-e752-11e8-99b1-257742e7c559.png)  
14. Navigate to the frontend directory within the shuttle tracker directory (`$GOPATH/src/github.com/wtg/shuttletracker/frontend`).  
15. Run `npm install`  
![Step 14+15](https://user-images.githubusercontent.com/6248819/48438927-1bd44c80-e753-11e8-9eb5-fc4c16795f57.png)  
16. Build the frontend using `npx vue-cli-service build --mode development`
    - _Note: if you are working on the frontend, you may instead use `npx vue-cli-service build --mode development --watch` in another terminal to continuously watch for changes and rebuild._  
![Step 16](https://user-images.githubusercontent.com/6248819/48438998-432b1980-e753-11e8-871c-3e6639e6383f.png)  
17. Go back up to the project root directory (using `cd ..`) and build Shuttle Tracker by running `go build -o shuttletracker.exe cmd/shuttletracker/main.go`  
![Step 17](https://user-images.githubusercontent.com/6248819/48439209-ba60ad80-e753-11e8-8e4a-740989a8ca55.png)  
18. Start the app by running `shuttletracker.exe` in the project root directory.  
19. Add yourself as an administrator by using `shuttletracker.exe admins --add RCS_ID`, replacing `RCS_ID` with your RCS ID. See the "Administrators" section below for more information.  
20. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the administration panel.  
21. Copy the information from [vehicles](https://shuttles.rpi.edu/vehicles), [routes](https://shuttles.rpi.edu/routes), and [stops](https://shuttles.rpi.edu/stops) into the admin panel if you want to mimic the current shuttle tracker site.  
