'use strict';

/* Controllers */

function ChartCtrl($scope, $http, $cookieStore) {

    $scope.filters = [
        {name:"", value:""},
        {name:"MagicWombats", value:"&team=MagicWombats"},
        {name:"ATeam", value:"&team=ATeam"}, 
    ];
    $scope.filter = $scope.filters[0];
    if ($cookieStore.get("filter")!=undefined) {
        $scope.filter =$cookieStore.get("filter");
    }

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
    if ($cookieStore.get("board")!=undefined) {
        $scope.board =$cookieStore.get("board");
    }
    $http({
        method: 'GET',
        url: '/data/boards'
    }).success(function(data, status, headers, config) {
        $scope.boards = data;
    });

    $scope.boardChanged = function(){
        $cookieStore.put("board", $scope.board);
        $scope.refreshSprints();
    }

    $scope.sprintChanged = function(){
        $cookieStore.put("sprint", $scope.sprint);
        $scope.refreshChart();
    }

    $scope.sprints = [];
    $scope.sprint = {};
    if ($cookieStore.get("sprint")!=undefined) {
        $scope.sprint =$cookieStore.get("sprint");
    }
    $scope.refreshSprints = function() {
        console.log("refreshSprints");
        $http({
            method: 'GET',
            url: '/data/sprints?board=' + $scope.board.Id
        }).success(function(data, status, headers, config) {
            console.log("sprints done:"+$(data));
            $scope.sprints = data;

        });
    };

    $scope.refreshChart = function() {
        console.log("refreshChart");
        $cookieStore.put("filter", $scope.filter);
        $http({
            method: 'GET',
            url: '/data/burndown?board=' + $scope.board.Id + '&sprint=' + $scope.sprint.Id +$scope.filter.value
        }).success(function(data) {
            var sprintstart = Date.parse(data.sprintstart);
            var sprintend = Date.parse(data.sprintend);
            // console.log("sprint "+sprintstart+" "+sprintend);
            var totaldata = [[sprintstart, data.timeline[0].totalWorkInHours]];
            var lastelem;
            

            for (var i in data.timeline) {
                var elem = data.timeline[i];
                var timestamp = Date.parse(elem.timestamp);
                // console.log("timest"+timestamp);
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

    $scope.refreshChart();
}
ChartCtrl.$inject = ['$scope', '$http', '$cookieStore'];