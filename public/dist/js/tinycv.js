$(document).ready(function(){
    mark_active_link();
});

function mark_active_link() {
    var anchors = document.getElementsByTagName('a');
    var href = window.location.href;
    var hostname = window.location.protocol + '//' + window.location.hostname;
    var ma_app_url = app_url;
    ma_app_url = ma_app_url.substring(0, ma_app_url.length - 1);

    if (window.location.port == '8000')
        hostname += ':8000'
    else if (window.location.port == '8080')
        hostname += ':8080'
    else if (window.location.port == '5000')
        hostname += ':5000'

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
        }
    }
}