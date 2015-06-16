google.load('visualization', '1.1', {packages: ['corechart', 'imagelinechart']});

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

        dataTable.setValue(i, 0, formatGraphDate(new Date(date*1000)));
        dataTable.setValue(i, 1, data);
    }

    var chart = new google.visualization.LineChart(document.getElementById(containerId));
    chart.draw(dataTable, {width: 1500, height: 800, min: 0, interpolateNulls: true});
}

function drawAllGraph(graphData, containerId) {
    // Graph
    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Date');

    if (typeof(graphData.TPS) != "undefined" && graphData.TPS != null) {
        data.addColumn('number', 'TPS');
    }
    if (typeof(graphData.Median) != "undefined" && graphData.Median != null) {
        data.addColumn('number', 'Median');
    }
    if (typeof(graphData.Avg) != "undefined" && graphData.Avg != null) {
        data.addColumn('number', 'Avg');
    }

    data.addRows(graphData.Max - graphData.Min + 1);

    var j = 0
    var previousUsersAmount = 0
    for (var i = graphData.Min; i <= graphData.Max; i++) {
        var iDate = new Date(i*1000)
        var hasData = false
        var g = 0

        key = "key_" + i
        if (typeof(graphData.TPS) != "undefined" && graphData.TPS != null && typeof(graphData.TPS[key]) != "undefined") {
            g++
            data.setValue(j, g, graphData.TPS[key]);
            hasData = true;
        }
        if (typeof(graphData.Median) != "undefined" && graphData.Median != null && typeof(graphData.Median[key]) != "undefined") {
            g++
            data.setValue(j, g, graphData.Median[key]);
            hasData = true;
        }
        if (typeof(graphData.Avg) != "undefined" && graphData.Avg != null && typeof(graphData.Avg[key]) != "undefined") {
            g++
            data.setValue(j, g, graphData.Avg[key]);
            hasData = true;
        }

        //if (hasData) {
        data.setValue(j, 0, formatGraphDate(iDate)); // ts
        //}

        j++
    }

    var chart = new google.visualization.LineChart(document.getElementById(containerId));
    chart.draw(data, {width: 1500, height: 800, min: 0, interpolateNulls: true});
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