google.load('visualization', '1.1', {packages: ['corechart', 'line', 'imagelinechart']});

function drawGraph(graphData, title, containerId) {
    if (graphData == null || graphData.length <= 0) {
        document.getElementById(containerId).innerHTML = "<i>No Data</i>";
        return false
    }

    // Graph
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn('date', 'Timeline');
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

        dataTable.setValue(i, 0, new Date(date*1000));
        dataTable.setValue(i, 1, data);
    }

    var chart = new google.charts.Line(document.getElementById(containerId));
    chart.draw(dataTable, google.charts.Line.convertOptions(
        {
            //width: 1200, 
            height: 500, 
            min: 0, 
            interpolateNulls: true,
            curveType: 'none',
            chartArea: {left:0,top:0,width:'100%',height:'100%'}
        }
    ));
}

function drawTimerGraph(graphData, containerId) {
    // Graph
    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Timeline');

    data.addColumn('number', 'Min');
    data.addColumn('number', 'Avg');
    data.addColumn('number', 'Max');

    data.addRows(graphData.Avg.length);

    for (var i = 0; i < graphData.Avg.length; i++) {
        var min = graphData.Min[i]
        var avg = graphData.Avg[i]
        var max = graphData.Max[i]


        data.setValue(i, 0, new Date(avg[0]*1000));
        data.setCell(i, 1, min[1], formatTimer(min[1]));
        data.setCell(i, 2, avg[1], formatTimer(avg[1]));
        data.setCell(i, 3, max[1], formatTimer(max[1]));
    }

    var chart = new google.charts.Line(document.getElementById(containerId));
    chart.draw(data, google.charts.Line.convertOptions({
        //width: 1200, 
        height: 500, 
        min: 0, 
        curveType: 'none',
        chartArea: {left:0,top:0,width:'100%',height:'100%'}, 
        interpolateNulls: true,
        colors: ['blue', 'green', 'red'], 
        series: { colors: ['blue', 'green', 'red'] }
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
