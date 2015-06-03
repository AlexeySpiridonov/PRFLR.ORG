google.load('visualization', '1.1', {packages: ['corechart', 'imagelinechart']});

function drawGraph(graphData) {
    // Graph
    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Date');
    //data.addColumn('number', 'Users');
    data.addColumn('number', 'RPS');
    data.addColumn('number', 'Median');
    data.addColumn('number', 'Avg');

    data.addRows(graphData.Max - graphData.Min + 1);

    var j = 0
    var previousUsersAmount = 0
    for (var i = graphData.Min; i <= graphData.Max; i++) {
        var iDate = new Date(i*1000)
        var hasData = false

        key = "key_" + i
        /*if (typeof(usersStats) != "undefined" && typeof(usersStats[key]) != "undefined") {
            previousUsersAmount = usersStats[key]
            data.setValue(j, 1, usersStats[key]);
        } else {
            data.setValue(j, 1, previousUsersAmount);
        }*/
        if (typeof(graphData.RPS) != "undefined" && typeof(graphData.RPS[key]) != "undefined") {
            data.setValue(j, 1, graphData.RPS[key]);
            hasData = true;
        }
        if (typeof(graphData.Median) != "undefined" && typeof(graphData.Median[key]) != "undefined") {
            data.setValue(j, 2, graphData.Median[key]);
            hasData = true;
        }
        if (typeof(graphData.Avg) != "undefined" && typeof(graphData.Avg[key]) != "undefined") {
            data.setValue(j, 3, graphData.Avg[key]);
            hasData = true;
        }

        //if (hasData) {
        data.setValue(j, 0, formatGraphDate(iDate)); // ts
        //}

        j++
    }

    var chart = new google.visualization.LineChart(document.getElementById('graph_container'));
    chart.draw(data, {width: 1500, height: 800, min: 0, interpolateNulls: true});
}

function formatGraphDate(date) {
    year  = date.getFullYear() 
    month = (date.getMonth()+1) > 10 ? (date.getMonth()+1) : "0"+(date.getMonth()+1)
    day   = date.getDate()  > 10 ? date.getDate()  : "0"+date.getDate()
    hours   = date.getHours()   > 10 ? date.getHours()   : "0"+date.getHours()
    minutes = date.getMinutes() > 10 ? date.getMinutes() : "0"+date.getMinutes()
    seconds = date.getSeconds() > 10 ? date.getSeconds() : "0"+date.getSeconds()
    return day+"/"+month+"/"+year + " " + hours + ":" + minutes + ":" + seconds
}