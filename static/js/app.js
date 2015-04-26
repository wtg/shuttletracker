angular.module('ShuttleTracking', [
  'ShuttleTracking.controllers',
  'ngRoute'
]).
config(['$routeProvider', '$locationProvider', 
  function($routeProvider, $locationProvider) {
    $routeProvider.
      when('/', {
        templateUrl: 'static/partials/fullscreen_map.html',
        controller: 'indexCtrl'
      }).
      when('/admin', {
        templateUrl: 'static/partials/vehicles.html',
        controller: 'vehiclesCtrl'
      }).
      when('/admin/vehicles', {
        templateUrl: 'static/partials/vehicles.html',
        controller: 'vehiclesCtrl'
      }).
      when('/admin/tracking', {
        templateUrl: 'static/partials/updates.html',
        controller: 'updatesCtrl'
      }).
      otherwise({
        redirectTo: '/'
      });
      $locationProvider.html5Mode(true);
  }
]);