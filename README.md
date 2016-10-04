Shuttle Tracking v2
===================

Remaking the original [Shuttle Tracker](https://github.com/wtg/shuttle_tracking) using [Go](https://golang.org/), [Polymer Web Components](https://www.polymer-project.org/), and [MongoDB](https://www.mongodb.org/).

Setting Up
-----------------
1. Clone this repository using `git clone https://github.com/wtg/shuttle_tracking_2`
2. Make sure you have npm, bower, golang and mongodb installed
  * On Debian-based linux, run `sudo apt-get install nodejs npm golang mongodb` to install npm and go language packages
  * *On CentOs run `sudo yum install nodejs npm golang mongodb` instead*
  * Run `sudo npm install -g bower` to install bower
3. Run `bower install` inside shuttle tracking directory to install dependencies listed in bower.json
4. Rename conf.json.sample to conf.json
5. Edit conf.json with the following:
  * Data Feed: API with tracking information (iTrak in our case)
  * UpdateInterval: Number of seconds between each request to the data feed
  * MongoUrl: Url where MongoDB is located
  * MongoPort: Port where MongoDB is bound (default is 27017)
6. Run the app using `go run main.go` in the project root directory
7. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the admin panel 
