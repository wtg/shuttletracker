Shuttle Tracking v2
===================

Remaking the original [Shuttle Tracker](https://github.com/wtg/shuttle_tracking) using [Go](https://golang.org/), [AngularJS](https://angularjs.org/), [Polymer Web Components](https://www.polymer-project.org/0.5/), and [MongoDB](https://www.mongodb.org/).

Development Notes
-----------------
1. Clone this repository using `git clone https://github.com/wtg/shuttle_tracking_2`
2. Run `bower install` to install dependencies listed in bower.json
3. Rename conf.json.sample to conf.json
4. Edit conf.json with the following:
  * Data Feed: API with tracking information (iTrak in our case)
  * UpdateInterval: Number of seconds between each request to the data feed
  * MongoUrl: Url where MongoDB is located
  * MongoPort: Port where MongoDB is bound (default is 27017)
5. Run the app using `go run main.go` in the project root directory
6. Visit http://localhost:8080/ to view the tracking application and http://localhost:8080/admin to view the admin panel 