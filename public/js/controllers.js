'use strict';

/* Controllers */

function ChartCtrl($scope, $http) {

    $scope.error = {
        msg: ""
    };
    $scope.messages = [];

    $scope.boardId = 51;
    $scope.sprintId = 217;

    $scope.filters = [
        {name:"MagicWombats", value:"&filter=MagicWombats"},
        {name:"ATeam", value:"&filter=ATeam"}, 
        {name:"none", value:""}
        ];
    $scope.filter = $scope.filters[0];

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
            url: '/data/burndown?boardId=' + $scope.board.Id + '&sprintId=' + $scope.sprint.Id +$scope.filter.value
        }).success(function(data) {
            var sprintstart = Date.parse(data.sprintstart);
            var sprintend = Date.parse(data.sprintend);
            console.log("sprint "+sprintstart+" "+sprintend);
            var totaldata = [[sprintstart, data.timeline[0].totalWorkInHours]];
            var lastelem;
            

            for (var i in data.timeline) {
                var elem = data.timeline[i];
                var timestamp = Date.parse(elem.timestamp);
                console.log("timest"+timestamp);
                if(timestamp>sprintstart && timestamp<sprintend){
                    totaldata.push([timestamp, elem.remainingWorkInHours]);
                    lastelem = elem.remainingWorkInHours;
                }

            }
            var now  = new Date();
            // totaldata.push([now.getTime(), lastelem]);

            $scope.chart.data[2].data = totaldata;

            var idealdata = [[sprintstart, data.timeline[0].totalWorkInHours],[sprintend, 1]];
            
            $scope.chart.data[1].data = idealdata;            

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