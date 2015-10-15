$(document).ready(function(){
    mark_active_link();
});

function mark_active_link() {
    var anchors = document.getElementsByTagName('a');
    var href = window.location.href;
    var hostname = window.location.protocol + '//' + window.location.hostname;
    var ma_app_url = app_url;
    ma_app_url = ma_app_url.substring(0, ma_app_url.length - 1);

    // if (window.location.port == '8000')
    //     hostname += ':8000'
    // else if (window.location.port == '8080')
    //     hostname += ':8080'
    // else if (window.location.port == '5000')
    //     hostname += ':5000'

    for (var i = 0; i < anchors.length; i++)
    {
        if (anchors[i].href == href)
        {
            var node = anchors[i];
            if (node.className.indexOf("skip-active") == -1) {
                node.className = 'active';
            }
            while (node.parentNode)
            {
                if (node.tagName.toLowerCase() == "li" && node.className.indexOf("skip-active") == -1)
                {
                    node.className = 'active';
                }
                node = node.parentNode;
            }
        }
        else if (href.search(anchors[i].href) == 0)
        {
            if (hostname != anchors[i].href && hostname + '/' != anchors[i].href && anchors[i].href != ma_app_url && anchors[i].href != ma_app_url + '/')
            {
                console.log(hostname + ' ' + anchors[i].href + " " + href)
                var node = anchors[i];
                if (node.className.indexOf("skip-active") == -1) {
                    // node.className = 'active';
                }
                while (node.parentNode)
                {
                    if (node.tagName.toLowerCase() == "li" && node.className.indexOf("skip-active") == -1)
                    {
                        node.className = 'active';
                    }
                    node = node.parentNode;
                }
            }
        }
    }
}

function delete_account() {
    window.location.href = app_url + 'account-delete';
}

function get_months() {
    var months = ["Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec"];
    var newMonths = [];
    var d = new Date();
    var m = d.getMonth();

    for (var i = m + 1; i < 12; i++) {
        newMonths.push(months[i]);
    }

    for (var i = 0; i <= m; i++) {
        newMonths.push(months[i]);
    }

    return newMonths;
}

function draw_chart() {
    /* ChartJS
     * -------
     * Here we will create a few charts using ChartJS
     */

    //--------------
    //- AREA CHART -
    //--------------

    // Get context with jQuery - using jQuery's .get() method.
    var areaChartCanvas = $("#areaChart").get(0).getContext("2d");
    // This will get the first returned node in the jQuery collection.
    // var areaChart = new Chart(areaChartCanvas);

    var areaChartData = {
        labels: get_months(),
        datasets: [
            {
                label: "Digital Goods",
                fillColor: "rgba(60,141,188,0.9)",
                strokeColor: "rgba(60,141,188,0.8)",
                pointColor: "#3b8bba",
                pointStrokeColor: "rgba(60,141,188,1)",
                pointHighlightFill: "#fff",
                pointHighlightStroke: "rgba(60,141,188,1)",
                data: __barecv_chart_data
            }
        ]
    };

    var areaChartOptions = {
        //Boolean - If we should show the scale at all
        showScale: true,
        //Boolean - Whether grid lines are shown across the chart
        scaleShowGridLines: false,
        //String - Colour of the grid lines
        scaleGridLineColor: "rgba(0,0,0,.05)",
        //Number - Width of the grid lines
        scaleGridLineWidth: 1,
        //Boolean - Whether to show horizontal lines (except X axis)
        scaleShowHorizontalLines: true,
        //Boolean - Whether to show vertical lines (except Y axis)
        scaleShowVerticalLines: true,
        //Boolean - Whether the line is curved between points
        bezierCurve: true,
        //Number - Tension of the bezier curve between points
        bezierCurveTension: 0.3,
        //Boolean - Whether to show a dot for each point
        pointDot: false,
        //Number - Radius of each point dot in pixels
        pointDotRadius: 4,
        //Number - Pixel width of point dot stroke
        pointDotStrokeWidth: 1,
        //Number - amount extra to add to the radius to cater for hit detection outside the drawn point
        pointHitDetectionRadius: 20,
        //Boolean - Whether to show a stroke for datasets
        datasetStroke: true,
        //Number - Pixel width of dataset stroke
        datasetStrokeWidth: 2,
        //Boolean - Whether to fill the dataset with a color
        datasetFill: true,
        //String - A legend template
        legendTemplate: "<ul class=\"<%=name.toLowerCase()%>-legend\"><% for (var i=0; i<datasets.length; i++){%><li><span style=\"background-color:<%=datasets[i].lineColor%>\"></span><%if(datasets[i].label){%><%=datasets[i].label%><%}%></li><%}%></ul>",
        //Boolean - whether to maintain the starting aspect ratio or not when responsive, if set to false, will take up entire container
        maintainAspectRatio: true,
        //Boolean - whether to make the chart responsive to window resizing
        responsive: true
    };

    //Create the line chart
    areaChart = new Chart(areaChartCanvas).Line(areaChartData, areaChartOptions);
}

if (typeof __barecv_draw_chart != 'undefined' && __barecv_draw_chart) {
    draw_chart();
}

if (typeof __barecv_create_editor != 'undefined' && __barecv_create_editor) {
    $("textarea").wysihtml5({
        toolbar: {
            "font-styles": false,
            "emphasis": true,
            "lists": true,
            "html": false,
            "link": true,
            "image": false,
            "color": false,
            "blockquote": false,
            "size": 'sm' //default: none, other options are xs, sm, lg
        }
    });
}

if (typeof __barecv_sortable_sections != 'undefined' && __barecv_sortable_sections) {
    //jQuery UI sortable for the todo list
    $(".todo-list").sortable({
       placeholder: "sort-highlight",
        handle: ".handle",
        forcePlaceholderSize: true,
        zIndex: 999999
    });
}