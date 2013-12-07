'use strict';

/* Controllers */

function ChartCtrl($scope, $http) {

    $scope.error = {msg:""};
    $scope.messages = [];

    $scope.boardId = 51;
    $scope.sprintId = 217;


    var dayMillies = 60 * 60 * 24 *1000;

    $scope.counter = 3;
    $scope.chart = {}
    $scope.chart.data = [{
        label: "total",
        data:[],
        bars: { show: true,
            barWidth: dayMillies/2,
            align: "left"
        }
    }, {
        label: "ideal",
        data:[],
        lines: { show: true }
    }, {
        label: "burndown",
        data:[],
        lines: { show: true }
    }];

    $http({method: 'GET', url: '/data/burndown?boardId='+$scope.boardId+'&sprintId='+ $scope.sprintId})
    .success(function(data, status, headers, config) {
        // var start = Date.parse(data.sprintstart);
        // var end = Date.parse(data.sprintend);
        var totaldata = [];
        for(var i in data.timeline){
            var elem = data.timeline[i];
            totaldata.push([Date.parse(elem.timestamp),elem.remainingWorkInHours]);
        }
        $scope.chart.data[2].data = totaldata;
    }).
    error(function(data, status, headers, config) {

    });

    $http({method: 'GET', url: '/data/boards'})
    .success(function(data, status, headers, config) {
        console.log("boards: ",data);
    }).
    error(function(data, status, headers, config) {

    });

    $http({method: 'GET', url: '/data/sprints?boardId='+$scope.boardId})
    .success(function(data, status, headers, config) {
        console.log("sprints: ",data);
    }).
    error(function(data, status, headers, config) {

    });

}
ChartCtrl.$inject = ['$scope', '$http'];
