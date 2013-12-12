'use strict';

/* Directives */


angular.module('jiraApp.directives', []).
directive('chart', function() {
    return {
        restrict: 'A',
        link: function(scope, elem, attrs) {

            var dayMillies = 60 * 60 * 24 * 1000;

            function weekendAreas(axes) {

                var markings = [],
                    d = new Date(axes.xaxis.min);

                // go to the first Saturday

                d.setUTCDate(d.getUTCDate() - ((d.getUTCDay() + 1) % 7))
                d.setUTCSeconds(0);
                d.setUTCMinutes(0);
                d.setUTCHours(0);

                var i = d.getTime();

                // when we don't set yaxis, the rectangle automatically
                // extends to infinity upwards and downwards

                do {
                    markings.push({
                        color: '#ccddcc',
                        xaxis: {
                            from: i,
                            to: i + 2 * 24 * 60 * 60 * 1000
                        }
                    });
                    i += 7 * 24 * 60 * 60 * 1000;
                } while (i < axes.xaxis.max);

                return markings;
            }

            function isWeekend(d) {
                return (d.getDay() % 6) == 0;
            }

            function countWeekendDays(from, to) {
                var cnt = 0;
                for (var i = from; i < to; i += dayMillies) {
                    if (isWeekend(new Date(i))) {
                        cnt += 1;
                    }
                }
                return cnt;
            }


            console.log("elem" + elem);
            var chart = null,
                opts = {
                    xaxis: {
                        mode: "time",
                        tickLength: 10
                    },
                    grid: {
                        markings: weekendAreas
                    },
                    yaxis: {
                        min: 0
                    }
                };

            scope.$watch(attrs.ngModel, function(data) {
                if (!chart) {
                    chart = $.plot(elem, data, opts);
                    elem.show();
                } else {
                    chart.setData(data);
                    chart.setupGrid();
                    chart.draw();
                }
            }, true);
        }
    }


});