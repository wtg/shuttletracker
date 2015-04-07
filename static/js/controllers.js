angular.module('ShuttleTracking.controllers',[]).

controller("indexCtrl", function($scope){
  $scope.heading = "Shuttle Tracker";
}).

controller("adminCtrl", function($scope){
  $scope.heading = "Admin Panel";
}).

controller('vehiclesCtrl', function($scope){
  $scope.heading = "Vehicles";
});