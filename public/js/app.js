'use strict';

var jiraApp = angular.module('jiraApp', [ 'jiraApp.services', 'jiraApp.directives', 'ngCookies', 'ngRoute']);

jiraApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/burndown', {
        templateUrl: 'partials/burndown.html',
        controller: 'BurndownCtrl'
      }).
      when('/delta', {
        templateUrl: 'partials/delta.html',
        //controller: 'ChartCtrl'
      }).
      otherwise({
        redirectTo: '/'
      });
  }]);