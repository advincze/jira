'use strict';

/* Controllers */

function ChartCtrl($scope, $http) {

    $scope.error = {
        msg: ""
    };
    $scope.messages = [];

    $scope.boardId = 51;
    $scope.sprintId = 217;


    var dayMillies = 60 * 60 * 24 * 1000;

    $scope.counter = 3;
    $scope.chart = {}
    $scope.chart.data = [{
        label: "total",
        data: [],
        bars: {
            show: true,
            barWidth: dayMillies / 2,
            align: "left"
        }
    }, {
        label: "ideal",
        data: [],
        lines: {
            show: true
        }
    }, {
        label: "burndown",
        data: [],
        lines: {
            show: true
        }
    }];

    $scope.filters = ["MagicWombats"]

    $scope.boards = [];
    $scope.board = {};
    $http({
        method: 'GET',
        url: '/data/boards'
    }).success(function(data, status, headers, config) {
        $scope.boards = data;
    });

    $scope.sprints = [];
    $scope.sprint = {};
    $scope.refreshSprints = function() {
        $http({
            method: 'GET',
            url: '/data/sprints?boardId=' + $scope.board.Id
        }).success(function(data, status, headers, config) {
            $scope.sprints = data.Sprints;
        });
    };

    $scope.refreshChart = function() {
        $http({
            method: 'GET',
            url: '/data/burndown?boardId=' + $scope.board.Id + '&sprintId=' + $scope.sprint.Id
        }).success(function(data) {
            var totaldata = [];
            for (var i in data.timeline) {
                var elem = data.timeline[i];
                totaldata.push([Date.parse(elem.timestamp), elem.remainingWorkInHours]);
            }
            $scope.chart.data[2].data = totaldata;
        });
    }

    $scope.selects = [{
        "id": "1",
        "name": "<i class=\"icon-star\"></i>&nbsp;foo"
    }, {
        "id": "2",
        "name": "<i class=\"icon-heart\"></i>&nbsp;bar"
    }, {
        "id": "3",
        "name": "<i class=\"icon-fire\"></i>&nbsp;baz"
    }];
    $scope.selectedItem = "1";
}
ChartCtrl.$inject = ['$scope', '$http'];