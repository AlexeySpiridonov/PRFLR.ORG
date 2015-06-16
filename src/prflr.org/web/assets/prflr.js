function start(){
	$.ajaxSetup({cache: false}); // turn off ajax cache

	// Menu Item Handlers
	$('#tab_menu a').click(function(){
	    if (inProgress) {
	        return false;
	    }

	    $('.profiler_block :input').unbind();
	    $('.profiler_block').hide();
	    var selector = $(this).attr('href');
	    $(selector).show();

        // flurry
        if (typeof(FlurryAgent) != 'undefined') {
            FlurryAgent.logEvent('trackMenuItemClick', {'item': selector});
        }

		var filter    = $('.profiler_block:visible').find('input[name=filter]').val();
	    var filterVal = typeof(filter) != 'undefined' && filter.length > 0 ? filter : '*/*/*/*';

        $(selector).find('input[name=filter]').val(filterVal);
        $(selector+' :input').not('input[name="filter"]').change(function(e){
            renderDataGrid(selector);
        });
        $(selector+' .refresh_button').click(function(){
            renderDataGrid(selector);
        });

        $('#tab_menu a').removeClass('tabselected');
        $(this).addClass('tabselected');

        renderDataGrid(selector, false);

	    return false;
	});

    // Special for Graph
    $('#tab_menu a[href="#graph"]').click(function(){
        $('#refresh_graph_button').click(function(){
            renderGraph('#graph')
        }).click()
    })

	// URL anchor handlers
	var hash = window.location.hash.replace('#', '');
	if (hash.length > 0) {
		hash = hash.split("|");

		if (typeof(hash[0]) == 'undefined') {
			hash[0] = 'aggregate';
		}

		$('input[name=filter]').val(hash[1]);
	    $('#tab_menu a[href="#'+hash[0]+'"]').click();
	} else {
	    $('#tab_menu a[href="#aggregate"]').click();
	}

	$('.prflrItemHeader').click(function(){
		assignFilterChunkValue($(this).attr('item'), '*');
		renderDataGrid(getCurrentMenuSelector());
	});

	$('.resetFilter').click(function(){
		assignFilterChunkValue('*', '*');
		renderDataGrid(getCurrentMenuSelector());
	});

    // Settings
    //alert($('#removeDataButton'))
    $('#removeDataButton').click(function(){
        if (!confirm("It will remove all stored data. Are you sure?")) {
            return false;
        }

        $(this).html("Removing...")

        $.get('/removeData', function(){
            $('#removeDataButton').html("Remove data")
            alert('Successfully!')
        });
    })
}

function clickMenuItem(item)
{
    item = item.replace("#", "");
    $('#tab_menu a[href=#'+item+']').click();
}

// Row items click handlers
function initProfilerItemsClickHandler()
{
	$(".prfrlItem").click(function(event){
		var item  = $(this).attr("item");
		var value = $(this).text();

		assignFilterChunkValue(item, value);

		var selector = getCurrentMenuSelector();

		event.stopPropagation();

		renderDataGrid(selector);
	});
}

function assignFilterChunkValue(chunk, value)
{
	var filter = $('input[name=filter]');

	if (chunk == '*') {
		filter.val(value+'/'+value+'/'+value+'/'+value);
		return true;
	}

	var chunkToSlot = {
		"src":   0,
		"timer": 1,
		"info":  2,
		"thrd":  3
	};

	if (typeof(chunkToSlot[chunk]) == 'undefined') {
		return false;
	}

	var chunks = $('input[name=filter]:visible').val().split('/');

	chunks[chunkToSlot[chunk]] = value;

	filter.val(chunks.join('/'));
}

function getCurrentMenuSelector()
{
	return $(".tabselected").attr('href');
}

function round(value)
{
	return Math.round(value*100)/100;
}

function formatNumber(number)
{
    var label = 'ms';
    if (number > 1000) {
        label = 'sec';
        number = number/1000;
    } else if (number > 10) {
        number = Math.floor(number);
    }
    return round(number) + label;
}

function renderGraph(selector)
{
    var elem = $(selector);
    var container = $('#graph_container')
    var button = elem.find('#refresh_graph_button');
    var query  = "/" + elem.attr('id') + "/?" + elem.find(':input').serialize();

    var filter = $('input[name=filter]:visible').val();
    filter     = typeof(filter) != 'undefined' && filter.length > 0 ? filter : '*/*/*/*';
    window.location.hash = "#"+getCurrentMenuSelector().replace("#", '')+"|"+filter;

    inProgress = true; // set to true (C.O.)  (^_^)
    elem.find(':input').attr('disabled', 'disabled');

    container.html("Loading...")
    button.css('color', 'grey').html('Loading...');
    $.getJSON(query, function(data){
        if (data != null) {
            //drawGraph({"Min": data.Min, "Max": data.Max, "Avg": data.Avg}, 'graph_container_avg')
            //drawGraph({"Min": data.Min, "Max": data.Max, "RPS": data.TPS}, 'graph_container_tps')
            drawGraph(data.Avg, 'Avg', 'graph_container_avg')
            drawGraph(data.TPS, 'TPS', 'graph_container_tps')
        }
    }).complete(function(){
        //grid.css('opacity', 1);
        button.css('color', 'black').html('Refresh');
        elem.find(':input').attr('disabled', null);

        initProfilerItemsClickHandler();

        inProgress = false;
    });
}

function renderDataGrid(selector, checkEmpty)
{
    var elem = $(selector);
    var grid = elem.find('table.profiler_grid');
    var button = elem.find('.refresh_button');
    var query  = "/" + elem.attr('id') + "/?" + elem.find(':input').serialize();

	var filter = $('input[name=filter]:visible').val();
    filter     = typeof(filter) != 'undefined' && filter.length > 0 ? filter : '*/*/*/*';
	//window.location.hash = "#"+filter+"|"+getCurrentMenuSelector().replace("#", '');
    window.location.hash = "#"+getCurrentMenuSelector().replace("#", '')+"|"+filter;

    if (grid.length == 0) {
        return false;
    }

    if (typeof(checkEmpty) == 'undefined') {
        checkEmpty = false;
    }
    if (checkEmpty && grid.html().length > 0) {
        return false;
    }

    inProgress = true; // set to true (C.O.)  (^_^)
    elem.find(':input').attr('disabled', 'disabled');

    grid.css('opacity', 0.3);
    button.css('color', 'grey').html('Loading...');
    $.getJSON(query, function(data){
        grid.empty().append('<tr class="b1"><td colspan="5">&nbsp;</td></tr>');

        if (data == null) return false;
        // first calculate line bars scale
        // we should get the biggest max value and divide lineBarLength on this value
        var maxMax = 0.000001;
        $.each(data, function(i, item){
            if (typeof(item.Max) == 'undefined') return false;

            if (item.Max > maxMax) {
                maxMax = item.Max;
            }
        });

        var scale = lineBarLength / maxMax;
        $.each(data, function(i, item){
            var dd = [];
            if (typeof(item.Src) != 'undefined' && item.Src != '') {
                dd.push('<span class="f18 prfrlItem" item="src">'+item.Src+'</span>')
            }
            if (typeof(item.Timer) != 'undefined' && item.Timer != '') {
                dd.push('<span class="f25 prfrlItem" item="timer">'+item.Timer+'</span>')
            }
            if (typeof(item.Info) != 'undefined' && item.Info != '') {
                dd.push('<span class="f15 prfrlItem" item="info">'+item.Info+'</span>')
            }
            if (typeof(item.Thrd)  != 'undefined' && item.Thrd != '') {
                dd.push('<span class="f12 prfrlItem" item="thrd">'+item.Thrd+'</span>')
            }

			var min = item.Min;
            var avg = item.Total / item.Count;
            var max = item.Max;

            grid.append(''+
                '<tr class="b1 prflrRow">'+
                '    <td class="r1">' + dd.join(' / ')+'</td>'+
            (typeof(item.Time) != 'undefined' ?
                '    <td class="r2"></td><td class="r3 f12">&nbsp;<br>&nbsp;<br>&nbsp;</td><td align="right" class="r4 f15">'+formatNumber(item.Time)+'</td>' 
            :
                '    <td class="r2">'+
                '        <div class="bln" style="width:'+(min > 0 ? round(min*scale) : 1)+'px;"/>'+
                '        <div class="gln" style="width:'+(avg > 0 ? round(avg*scale) : 1)+'px;"/>'+
                '        <div class="rln" style="width:'+(max > 0 ? round(max*scale) : 1)+'px;"/>'+
                '    </td>'+
                '    <td class="r3 f12">'+formatNumber(min)+'<br>'+formatNumber(avg)+'<br>'+formatNumber(max)+'</td>'+
                '    <td align="right" class="r4 f15">'+
                '        '+formatNumber(item.Total)+'<br/>'+
                '        '+item.Count+
                '    </td>')+
                '</tr>'+
                '');
        });
    }).complete(function(){
        grid.css('opacity', 1);
        button.css('color', 'black').html('Refresh');
        elem.find(':input').attr('disabled', null);

		initProfilerItemsClickHandler();

        inProgress = false;
    });
}

function Logout()
{
    return confirm('Are you sure you want to Logout?');
}

function resetApiKey()
{
    return confirm('Are you sure you want to RESET your Api Key?');
}

function changePassword()
{
    var password = prompt('Insert your new password');
    return false;
}