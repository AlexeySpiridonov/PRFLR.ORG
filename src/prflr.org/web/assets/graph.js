google.load('visualization', '1.1', {packages: ['corechart', 'line', 'imagelinechart']});

function drawGraph(graphData, title, containerId) {
    if (graphData == null || graphData.length <= 0) {
        document.getElementById(containerId).innerHTML = "<i>No Data</i>";
        return false
    }

    // Graph
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn('string', 'Date');
    dataTable.addColumn('number', title);

    dataTable.addRows(graphData.length);

    for (var i = 0; i < graphData.length; i++) {
        var elem = graphData[i]

        var date = elem[0]
        var data = elem[1]

        // check if it's Float
        if (Math.round(data) != data) {
            data = data.toFixed(2)
        }

        dataTable.setValue(i, 0, formatGraphDate(new Date(date*1000)));
        dataTable.setValue(i, 1, data);
    }

    var chart = new google.visualization.LineChart(document.getElementById(containerId));
    chart.draw(dataTable, {width: 900, height: 500, min: 0, interpolateNulls: true});
}

function drawTimerGraph(graphData, containerId) {
    // Graph
    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Date');

    data.addColumn('number', 'Min');
    data.addColumn('number', 'Avg');
    data.addColumn('number', 'Max');

    data.addRows(graphData.Avg.length);

    for (var i = 0; i < graphData.Avg.length; i++) {
        var min = graphData.Min[i]
        var avg = graphData.Avg[i]
        var max = graphData.Max[i]

        data.setValue(i, 0, formatGraphDate(new Date(avg[0]*1000)));
        //data.setValue(i, 1, formatTimer(min[1]));
        //data.setValue(i, 2, formatTimer(avg[1]));
        //data.setValue(i, 3, formatTimer(max[1]));
        data.setCell(i, 1, min[1], formatTimer(min[1]));
        data.setCell(i, 2, avg[1], formatTimer(avg[1]));
        data.setCell(i, 3, max[1], formatTimer(max[1]));
    }

    var chart = new google.charts.Line(document.getElementById(containerId));
    chart.draw(data, google.charts.Line.convertOptions({
        width: 900, 
        height: 500, 
        min: 0, 
        curveType: 'function',
        legend: { position: 'bottom' },
        chartArea: {left:0,top:0,width:'100%',height:'100%'}, 
        interpolateNulls: true, 
        series: {
            0: { color: '#0000FF' },
            1: { color: '#00FF00' },
            2: { axis: 'Temps', color: '#FF0000' }
        },
        axes: {
          y: {
            'Temps': {label: 'Time, ms'}
          }
        }
    })
    );
}

function formatGraphDate(date) {
    year  = date.getUTCFullYear() 
    month = (date.getUTCMonth()+1) > 10 ? (date.getUTCMonth()+1) : "0"+(date.getUTCMonth()+1)
    day   = date.getUTCDate()  > 10 ? date.getUTCDate()  : "0"+date.getUTCDate()
    hours   = date.getUTCHours()   > 10 ? date.getUTCHours()   : "0"+date.getUTCHours()
    minutes = date.getUTCMinutes() > 10 ? date.getUTCMinutes() : "0"+date.getUTCMinutes()
    seconds = date.getUTCSeconds() > 10 ? date.getUTCSeconds() : "0"+date.getUTCSeconds()
    return day+"/"+month+"/"+year + " " + hours + ":" + minutes + ":" + seconds
}

function formatTimer(timer) {
    if (timer > 10000) {
        timer = Math.round(timer/1000) + " sec"
    } else if (timer > 1000) {
        if (timer % 1000 > 500) {
            timer = Math.round(timer/1000) + ".5 sec"
        } else {
            timer = Math.round(timer/1000) + " sec"
        }
    } else {
        timer = timer.toFixed(2) + " ms"
    }
    return timer
}
