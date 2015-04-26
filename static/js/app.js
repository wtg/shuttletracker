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
      otherwise({
        redirectTo: '/'
      });
      $locationProvider.html5Mode(true);
  }
]);