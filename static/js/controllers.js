angular.module('ShuttleTracking.controllers',[]).

controller("indexCtrl", function($scope){
  $scope.heading = "Shuttle Tracker";
}).

controller("vehiclesCtrl", function($scope, $http){
  $scope.vehicles = {};

  $http.get('/vehicles').
    success(function(data, status, headers, config) {
      $scope.vehicles = data;
    }).
    error(function(data, status, headers, config) {
      // handle error
    });

  $scope.addVehicle = function() {
    var vehicleObj = {
      vehicleId: $scope.vehicleId,
      vehicleName: $scope.vehicleName
    };
    var res = $http.post('/vehicles/create', vehicleObj);
    res.success(function(data, status, headers, config) {
      $scope.message = data;
      $scope.vehicles.push(vehicleObj);
    });
    res.error(function(data, status, headers, config) {
      alert( "failure message: " + JSON.stringify({data: data}));
    });   
    // Clear input fields
    $scope.vehicleId = '';
    $scope.vehicleName = '';
  };
}).

controller("updatesCtrl", function($scope, $http) {
  $scope.pageTitle = "Tracking Updates";
})