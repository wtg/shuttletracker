angular.module('ShuttleTracking', [
  'ShuttleTracking.controllers',
  'ngRoute'
]).
config(['$routeProvider', '$locationProvider', 
  function($routeProvider, $locationProvider) {
    $routeProvider.
      when('/', {
        templateUrl: 'static/partials/map.html',
        controller: 'indexCtrl'
      }).
      when('/admin', {
        templateUrl: 'static/partials/dashboard.html',
        controller: 'adminCtrl'
      }).
      when('/vehicles', {
        templateURL: 'static/partials/vehicles',
        controller: 'vehiclesCtrl'
      }).
      otherwise({
        redirectTo: '/'
      });
      $locationProvider.html5Mode(true);
  }
]);