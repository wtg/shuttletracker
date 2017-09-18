var Admin = {
  RoutesMap: null,
  ShuttleRoutes: [],
  drawnRoute: new L.featureGroup(),
  /*
  <div class = "route-description-box">
    <span class = "emphasis">name:</span><span class ="content"> West Campus</span><br>
    <span class = "emphasis">description:</span><span class ="content"> The west campus Red Hawk shuttle will make stops at the Student Union Horseshoe, Sage Ave & Troy Building crosswalks, Blitman Residence Commons, 6th Avenue and City Station, Polytechnic Residence Commons, and 15th Street and College Avenue. This shuttle route will have a yellow indicator labeled "WEST" on the front of the shuttle.</span><br>
    <span class = "emphasis">enabled:</span><span class="content">true</span><br>
    <span class = "emphasis">color:</span><span class="content">#fff</span><br>
    <span class = "emphasis">time:</span><span class="content">idek</span><br>
    <span class = "emphasis">id:</span><span class="content">13445213a</span><br>
    <div style="float: right;"><a class="button">edit</a><a class="button">delete</a></div><br>

  </div>
  */

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

  populateRoutesPanel: function(data){
    $(".routePanel").html("");

    //console.log(data);
    //Admin.updateRoutes(data);
    if(data == null){

    }else{
      for(var i = 0; i < data.length; i ++){
        console.log(data[i]);
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
      $(".changeroute").click(function(){
        Admin.loadRoute($(this).attr("routeId"));
      });

    }
  },
  populateRouteForm: function(data){
    var polylineOptions = {
      color: data.color,
      weight: 3,
      opacity: 1,
    };
    var points = [];
    for(var j = 0; j < data.coords.length; j ++){
      points.push(new L.LatLng(data.coords[j].lat,data.coords[j].lng));
    }

    var polyline = new L.Polyline(points, polylineOptions);
    polyline.addTo(Admin.drawnRoute);

    $("#name").val(data.name);
    $("#desc").val(data.description);
    $("#en").val(data.enabled);
    $("#color").val(data.color);
    $("#time").val("lel");
    $("#width").val(data.width);





  },

  loadRoute: function(id){
    $.get( "/routes", function(data){
      for(var i = 0; i < data.length; i ++){
        if(data[i].id == id){
          Admin.populateRouteForm(data[i])
        }
      }
    });

  },

  buildRouteBox: function(routeInfo){
    var box = "";
    box += "<div id = " + routeInfo.id +" class = 'route-description-box'>";
    box += "<span class = 'emphasis'>name:</span><span class ='content'> " + routeInfo.name + "</span><br>";
    box += "<span class = 'emphasis'>description:</span><span class ='content'>" + routeInfo.description + "</span><br>"
    box += "<span class = 'emphasis'>enabled:</span><span class='content'>"+routeInfo.enabled + "</span><br>";
    box += "<span class = 'emphasis'>color:</span><span class='content'>" + routeInfo.color + "</span><br>";
    box += "<span class = 'emphasis'>time:</span><span class='content'>lel</span><br>";
    box += "<span class = 'emphasis'>id:</span><span class='content'>"+ routeInfo.id + "</span><br>";
    box += "<div style='float: right;width:auto;'><button class='button cbutton changeroute' routeId="+routeInfo.id +">change</button><button id='delete' routeId="+routeInfo.id +" class='button cbutton deleteroute'>delete</button></div><br></div>"
    $(".routePanel").append(box);

  },

  routeDeleteHandler: function(info){
    //console.log(info);
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

    Admin.RoutesMap.addLayer(Admin.drawnRoute);

    var drawControl = new L.Control.Draw({
      edit: {
        featureGroup: Admin.drawnRoute
      },

      draw: {
        polygon: false,
        marker: false,
        circle: false,
        rectangle: false,
      },
    });

    Admin.RoutesMap.addControl(drawControl);
    Admin.RoutesMap.on(L.Draw.Event.CREATED, function (e) {
      var type = e.layerType,
      layer = e.layer;
      Admin.drawnRoute.addLayer(layer);
    });

    var tmp = L.Routing.control({
      waypoints: [
        L.latLng(42.728172, -73.678803),
        L.latLng(42.728372, -73.678803)
      ],
      routeWhileDragging: true
    });
    tmp.addTo(Admin.RoutesMap);
    var waypoints =[
      L.latLng(42.728172, -73.678803),
      L.latLng(42.728372, -73.678803)
    ];

    Admin.RoutesMap.on('click', function(e) {
      waypoints.push(e.latlng);
      tmp.setWaypoints(waypoints);
    });

  },

  submitForm: function(){
    var coords = [];
    if(Admin.drawnRoute.toGeoJSON().features.length != 0){
      var data = Admin.drawnRoute.toGeoJSON().features[0].geometry.coordinates;
      for(var i = 0; i < data.length; i ++){
        coords.push({
          "lat": data[i][1],
          "lng": data[i][0]
        });
      }
    }
    var toSend = {
      "name":$("#name").val(),
      "description":$("#desc").val(),
      "startTime":"",
      "endTime":"",
      "enabled":$("#en").val(),
      "color":$("#color").val(),
      "width":$("#width").val(),
      "coords":JSON.stringify(coords)};
    $.ajax({
      url: "/routes/create",
      type: "POST",
      data: JSON.stringify(toSend),
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
