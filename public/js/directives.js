'use strict';

/* Directives */


angular.module('jiraApp.directives', []).
directive('chart', function() {
    return {
        restrict: 'E',
        link: function(scope, elem, attrs){

            var chart = null,
                opts  = {
                    xaxis: {
                        mode: "time",
                        tickLength: 5
                    }
                };


            scope.$watch(attrs.ngModel, function(data){
                if(!chart){
                    chart = $.plot(elem, data, opts);
                    elem.show();
                }else{
                    chart.setData(data);
                    chart.setupGrid();
                    chart.draw();
                }
            },true);
        }
    }


});