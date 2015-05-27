google.load('visualization', '1.1', {packages: ['corechart', 'imagelinechart']});

function drawGraph() {
    // Graph
    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Date');
    //data.addColumn('number', 'Users');
    data.addColumn('number', 'RPS');
    data.addColumn('number', 'Median');
    data.addColumn('number', 'Avg');

    data.addRows(tsMax - tsMin + 1);

    var j = 0
    var previousUsersAmount = 0
    for (var i = tsMin; i <= tsMax; i++) {
        var iDate = new Date(i*1000)

        data.setValue(j, 0, iDate.getHours() + ":" + iDate.getMinutes() + ":" + iDate.getSeconds()); // ts

        key = "key_" + i
        /*if (typeof(usersStats) != "undefined" && typeof(usersStats[key]) != "undefined") {
            previousUsersAmount = usersStats[key]
            data.setValue(j, 1, usersStats[key]);
        } else {
            data.setValue(j, 1, previousUsersAmount);
        }*/
        if (typeof(rpsStats) != "undefined" && typeof(rpsStats[key]) != "undefined") {
            data.setValue(j, 1, rpsStats[key]);
        }
        if (typeof(medianStats) != "undefined" && typeof(medianStats[key]) != "undefined") {
            data.setValue(j, 2, medianStats[key]);
        }
        if (typeof(avgStats) != "undefined" && typeof(avgStats[key]) != "undefined") {
            data.setValue(j, 3, avgStats[key]);
        }
        j++
    }

    var chart = new google.visualization.LineChart(document.getElementById('graph_container'));
    chart.draw(data, {width: 1500, height: 800, min: 0, interpolateNulls: true});
  }

  // Start it when the page is ready
  google.setOnLoadCallback(drawGraph);