var Admin = {
  RoutesMap: null,
  ShuttleRoutes: [],
  drawnRoute: null,
  RoutingControl: null,
  RoutingWaypoints: [],
  RouteData: null,
  StopMap: null,

  updateRoutes: function(data){
    var updatedRoute = [];
    for(var i = 0; i < data.length; i ++){
      var points = [];
      for(var j = 0; j < data[i].coords.length; j ++){
        points.push(new L.LatLng(data[i].coords[j].lat,data[i].coords[j].lng));
      }
      var polylineOptions = {
        color: data[i].color,
        weight: 3,
        opacity: 1,
      };
      if(data[i].width === 0){
        polylineOptions.dashArray = '10,10';
      }

      var polyline = new L.Polyline(points, polylineOptions);

      var r ={
        name: data[i].name,
        id: data[i].id,
        description: data[i].description,
        color: data[i].color,
        created: data[i].created,
        enabled: data[i].enabled,
        stops: data[i].stopsid,
        start_time: data[i].startTime,
        end_time: data[i].endTime,
        points: points,
        line: polyline
      };

      updatedRoute.push(r);

    }

    Admin.ShuttleRoutes = updatedRoute;
    Admin.drawRoutes();

  },


  drawRoutes: function(){

    var polylineOptions = {
      color: 'blue',
      weight: 1,
      opacity: 1
    };
    for(i = 0; i < Admin.ShuttleRoutes.length; i ++){
      Admin.RoutesMap.addLayer(Admin.ShuttleRoutes[i].line);
    }

  },

  bindStopButtons: function(){
    $(".stopDelete").click(function(e){
      var routeId = $(this).parent().attr("id");
      $.ajax({
        url: '/stops/' + $(this).parent().attr("stopId"),
        type: 'DELETE',
        success: function(result) {
          $.get( "/stops", function(e){
            Admin.showStopsPanel();
            Admin.populateStopsForm(e, routeId);
            Admin.bindStopButtons();
          });
        }
      });
    });
  },

  populateRoutesPanel: function(data){
    Admin.RouteData = data;
    $(".routePanel").html("");

    //console.log(data);
    //Admin.updateRoutes(data);
    if(data == null){

    }else{
      for(var i = 0; i < data.length; i ++){
        //console.log(data[i]);
        Admin.buildRouteBox(data[i]);

      }
      $(".deleteroute").click(function(){
        if ($(this).html() == "sure?"){
          $.ajax({
            url: '/routes/' + $(this).attr("routeId"),
            type: 'DELETE',
            success: function(result) {
              //Admin.populateRoutesPanel(data);
              $.get( "/routes", Admin.populateRoutesPanel);
            }
          });
        }else{
          $(this).html("sure?")
        }
      });
      $(".stops").click(function(){
        var a = $(this).attr("routeId");
        $.get( "/stops", function(e){
          Admin.showStopsPanel();
          Admin.populateStopsForm(e,a);
          Admin.bindStopButtons();


        });

      });

    }
  },
  populateRouteForm: function(data){

  },
  initStopMap: function(){
    if(Admin.StopsMap == null){
      Admin.StopsMap = L.map('newStopMap', {
        zoomControl: false,
        attributionControl: false
      });

      Admin.StopsMap.setView([42.728172, -73.678803], 15.3);
      Admin.StopsMap.addControl(L.control.attribution({
        position: 'bottomright',
        prefix: ''
      }));

      L.tileLayer('http://tile.stamen.com/toner-lite/{z}/{x}/{y}{r}.png', {
        attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
        minZoom: 13
      }).addTo(Admin.StopsMap);
    }

  },

  populateStopsForm: function(data,routeId){
    $(".stopPanel").html("");
    Admin.StopsMap= null;
    for (var i = -1; i < data.length; i ++){
      if (i != -1 && data[i].routeId == routeId){
        var box = "";
        box += "<div id=" + data[i].routeId + " stopId='"+data[i].id+"' class = 'route-description-box'>";
        box += "<span class = 'emphasis'>Name:</span><input type='text' value='" + data[i].name + "'></input><br>";
        box += "<span class = 'emphasis'>Description:</span><input type='text' value='" + data[i].description + "'></input><br>";
        box += "<span class = 'emphasis'>Route:</span><select>";

        for (var j = 0 ; j < Admin.RouteData.length; j++){
          box += "<option value="+ Admin.RouteData[j].routeId + ">" + Admin.RouteData[j].name + "</option>"
          //console.log(box);
        }

        box += "</select><br>"
        box += "<span class = 'emphasis'>Enabled:</span><input type='textbox' value="+data[i].enabled+"></input>"
        box += "<span class='button stopEdit' style='float:right;'>submit</span><span class='button stopDelete' style='float:right;'>delete</span></div>"

        $(".stopPanel").append(box);

      }else if(i == -1){
        var box = "";
        box += "<div id='' stopId='' class = 'route-description-box'>";
        box += "<div id='newStopMap'style='height: 50%;position: inherit;width: 50%;background-color:black;z-index:0;border-style: solid; border-width:1px; border-color:black; float: inherit;'; background-color:black;z-index:0;'></div>";

        box += "<span class = 'emphasis'>Name:</span><input type='text' value='New stop' ></input><br>";
        box += "<span class = 'emphasis'>Description:</span><input type='text' value=></input><br>";
        box += "<span class = 'emphasis'>Route:</span><select>";

        for (var j = 0 ; j < Admin.RouteData.length; j++){
          box += "<option value="+ Admin.RouteData[j].routeId + ">" + Admin.RouteData[j].name + "</option>"
          //console.log(box);
        }

        box += "</select><br>"
        box += "<span class = 'emphasis'>Enabled:</span><input type='textbox' value=></input>"
        box += "<span class='button stopEdit' style='float:right;'>submit</span></div>"
        $(".stopPanel").append(box);
        Admin.initStopMap();
      }

    }
  },

  loadRoute: function(id){

  },

  buildRouteBox: function(routeInfo){
    var box = "";
    box += "<div id = " + routeInfo.id +" class = 'route-description-box'>";
    box += "<span class = 'emphasis'>name:</span><span class ='content'> " + routeInfo.name + "</span><br>";
    box += "<span class = 'emphasis'>description:</span><span class ='content'>" + routeInfo.description + "</span><br>"
    box += "<span class = 'emphasis'>enabled:</span><span class='content'>"+routeInfo.enabled + "</span><br>";
    box += "<span class = 'emphasis'>color:</span><span class='content'>" + routeInfo.color + "</span><br>";
    box += "<span class = 'emphasis'>time:</span><span class='content'>"+routeInfo.startTime + "-" + routeInfo.endTime + "</span><br>";
    box += "<span class = 'emphasis'>id:</span><span class='content'>"+ routeInfo.id + "</span><br>";
    box += "<div style='float: right;width:auto;'><button class='button cbutton stops' routeId="+routeInfo.id +">stops</button><button id='delete' routeId="+routeInfo.id +" class='button cbutton deleteroute'>delete</button></div><br></div>"
    $(".routePanel").append(box);

  },

  routeDeleteHandler: function(info){
    //console.log(info);
  },

  showMapPanel: function(){
    Admin.hideStopsPanel();
    $('.mapPanel').css('display','block');
    Admin.RoutesMap.invalidateSize();
  },
  hideMapPanel: function(){
    $('.mapPanel').css('display','none');
    Admin.RoutesMap.invalidateSize();
  },
  hideStopsPanel: function(){
    $('.stopPanel').css('display','none');

  },
  showStopsPanel: function(){
    Admin.hideMapPanel();
    $('.stopPanel').css('display','block');
    if(Admin.StopsMap != null){
      Admin.StopsMap.invalidateSize();
    }

  },

  initMap: function(){
    Admin.RoutesMap = L.map('mapid', {
        zoomControl: false,
        attributionControl: false
    });

    Admin.RoutesMap.setView([42.728172, -73.678803], 15.3);
    Admin.RoutesMap.addControl(L.control.attribution({
        position: 'bottomright',
        prefix: ''
    }));

    L.tileLayer('http://tile.stamen.com/toner-lite/{z}/{x}/{y}{r}.png', {
      attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
      minZoom: 13
    }).addTo(Admin.RoutesMap);

    Admin.RoutingControl = L.Routing.control({
      waypoints: [

      ],
      routeWhileDragging: true
    });

    Admin.RoutingControl.on('routeselected', function(e) {
      if (Admin.drawnRoute != null){
        Admin.RoutesMap.removeLayer(Admin.drawnRoute);
      }
      Admin.drawnRoute = L.polyline(e.route.coordinates, {color: 'blue'});
      Admin.drawnRoute.addTo(Admin.RoutesMap)

    });
    Admin.RoutingControl.addTo(Admin.RoutesMap);
    Admin.RoutingWaypoints =[

    ];

    Admin.RoutesMap.on('click', function(e) {
      Admin.RoutingWaypoints.push(e.latlng);
      Admin.RoutingControl.setWaypoints(Admin.RoutingWaypoints);
    });

  },

  removeLastPoint: function(){
    Admin.RoutingWaypoints = Admin.RoutingWaypoints.slice(0, -1);
    Admin.RoutingControl.setWaypoints(Admin.RoutingWaypoints);

  },
  pullForm: function(){
    var coords = [];
    coords = Admin.drawnRoute.getLatLngs();

    var toSend = {
      "name":$("#name").val(),
      "description":$("#desc").val(),
      "startTime":$("#startTimeDay").val() +";"+ $("#startTimeTime").val(),
      "endTime":$("#endTimeDay").val() +";"+ $("#endTimeTime").val(),
      "enabled":$("#en").val(),
      "color":$("#color").val(),
      "width":$("#width").val(),
      "coords":JSON.stringify(coords)};
    Admin.hideMapPanel();
    return toSend;

  },

  submitChange: function(){
      Admin.submitForm($('#jsonField').val());
      console.log($('#jsonField').val());
      $('.prompt').css('display','none');
  },

  getJson: function(){
    var toSend = Admin.pullForm()
    var wnd = window.open("about:blank", "", "_blank");
    wnd.document.write(JSON.stringify(toSend));

  },

  submitForm: function(toSend){
    $.ajax({
      url: "/routes/create",
      type: "POST",
      data: toSend,
      contentType: "application/json",
      complete: function(data){
        $.get( "/routes", Admin.populateRoutesPanel);
      }
    });

  }

};
$(document).ready(function(){
  Admin.initMap();
  $.get( "/routes", Admin.populateRoutesPanel);


});
