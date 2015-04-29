angular.module('ShuttleTracking', [
  'ShuttleTracking.controllers',
  'ngRoute'
]).
config(['$routeProvider', '$locationProvider', 
  function($routeProvider, $locationProvider) {
    $routeProvider.
      when('/', {
        title: 'Shuttle Tracking',
        templateUrl: 'static/partials/fullscreen_map.html',
        controller: 'indexCtrl'
      }).
      when('/admin', {
        title: 'Vehicles',
        templateUrl: 'static/partials/vehicles.html',
        controller: 'vehiclesCtrl'
      }).
      when('/admin/vehicles', {
        title: 'Vehicles', 
        templateUrl: 'static/partials/vehicles.html',
        controller: 'vehiclesCtrl'
      }).
      when('/admin/routes', {
        title: 'Routes',
        templateUrl: 'static/partials/routes.html',
        controller: 'routesCtrl'
      }).
      when('/admin/stops', {
        title: 'Stops',
        templateUrl: 'static/partials/stops.html',
        controller: 'stopsCtrl'
      }).
      when('/admin/tracking', {
        title: 'Tracking Updates',
        templateUrl: 'static/partials/updates.html',
        controller: 'updatesCtrl'
      }).
      when('/admin/schedule', {
        title: 'Schedule',
        templateUrl: 'static/partials/schedule.html',
        controller: 'scheduleCtrl'
      }).
      when('/admin/users', {
        title: 'Users',
        templateUrl: 'static/partials/users.html',
        controller: 'usersCtrl'
      }).
      otherwise({
        redirectTo: '/admin'
      });
      $locationProvider.html5Mode(true);
  }
])
.run(['$rootScope', '$route', function($rootScope, $route) {
  $rootScope.$on('$routeChangeSuccess', function() {
    $rootScope.title = $route.current.title;
  });
}]);