angular.module('ShuttleTracking.controllers',[]).

controller("indexCtrl", function($scope, $http){
}).

controller("vehiclesCtrl", function($scope, $http) {
  $scope.vehicles = {};

  $http.get('/vehicles').
    success(function(data) {
      $scope.vehicles = data;
    }).
    error(function(data) {
      console.log(data);
    });

  $scope.addVehicle = function() {
    var vehicleObj = {
      vehicleID: $scope.vehicleID,
      vehicleName: $scope.vehicleName
    };
    var res = $http.post('/vehicles/create', vehicleObj);
    res.success(function(data) {
      $scope.message = data;
      $scope.vehicles.push(vehicleObj);
    });
    res.error(function(data) {
      console.log("failure message: " + JSON.stringify({data: data}));
    });   
    $scope.vehicleID = '';
    $scope.vehicleName = '';
  };
}).

controller("routesCtrl", function($scope, $http) {
  $scope.routes = {};

  $http.get('/routes').
    success(function(data) {
      $scope.routes = data;
    }).
    error(function(data) {
      console.log(data);
    });

  $scope.addRoute = function() {
    var routeObj = {
      name: $scope.name,
      description: $scope.description,
      startTime: $scope.startTime,
      endTime: $scope.endTime,
      enabled: $scope.enabled,
      color: $scope.color,
      width: $scope.width
    };
    var res = $http.post('/routes/create', routeObj);
    res.success(function(data) { 
      $scope.message = data;
      $scope.routes.push(routeObj);
    });
    res.error(function(data) {
      console.log("failure message: " + JSON.stringify({data: data}));
    });
    $scope.name = '';
    $scope.description = '';
    $scope.startTime = '';
    $scope.endTime = '';
    $scope.enabled = '';
    $scope.color = '';
    $scope.width = '';
  };
}).

controller("stopsCtrl", function($scope) {
  
}).

controller("updatesCtrl", function($scope) {

}).

controller("scheduleCtrl", function($scope) {

}).

controller("usersCtrl", function($scope) {

}).

controller("updatesCtrl", function($scope, $http) {
}).

controller("stopsCtrl", function($scope, $http) {
})