angular.module('ShuttleTracking.controllers',[]).

controller("indexCtrl", function($scope, $http){
}).

controller("vehiclesCtrl", function($scope, $http){
  $scope.vehicles = {};

  $http.get('/vehicles').
    success(function(data) {
      $scope.vehicles = data;
    }).
    error(function(data) {
      // handle error
    });

  $scope.addVehicle = function() {
    var vehicleObj = {
      vehicleId: $scope.vehicleId,
      vehicleName: $scope.vehicleName
    };
    var res = $http.post('/vehicles/create', vehicleObj);
    res.success(function(data) {
      $scope.message = data;
      $scope.vehicles.push(vehicleObj);
    });
    res.error(function(data) {
      alert( "failure message: " + JSON.stringify({data: data}));
    });   
    // Clear input fields
    $scope.vehicleId = '';
    $scope.vehicleName = '';
  };
}).

controller("routesCtrl", function($scope) {

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