var key = "";

function initializeAPIKeys(data){
  key = data;
  $.get( "/routes", Admin.populateRoutesPanel);
  Admin.initMap();

}

var Routes = {
  RoutesMap: null,
  ShuttleRoutes: [],
  drawnRoute: null,
  RoutingControl: null,
  RoutingWaypoints: [],
  RouteData: null,

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

    Routes.ShuttleRoutes = updatedRoute;
    Routes.drawRoutes();

  },
  drawRoutes: function(){
    var polylineOptions = {
      color: 'blue',
      weight: 1,
      opacity: 1
    };
    for(i = 0; i < Routes.ShuttleRoutes.length; i ++){
      Routes.RoutesMap.addLayer(Routes.ShuttleRoutes[i].line);
    }

  },

  populateRoutesPanel: function(data){
    Routes.RouteData = data;
    $(".routePanel").html("");

    //console.log(data);
    //Routes.updateRoutes(data);
    if(data === null){

    }else{
      for(var i = 0; i < data.length; i ++){
        //console.log(data[i]);
        Routes.buildRouteBox(data[i]);

      }
      Routes.buildSubmitBox();

      $(".addStopJson").click(function(){
        Routes.submitForm($(".json").val());
      });

      $('html').click(function() {
        $(".deleteroute").html("Delete");
      });

      $(".deleteroute").click(function(e){
        e.stopPropagation()
        if ($(this).html() == "Confirm deletion"){
          $.ajax({
            url: '/routes/' + $(this).attr("routeId"),
            type: 'DELETE',
            success: function(result) {
              //Routes.populateRoutesPanel(data);
              $.get( "/routes", Routes.populateRoutesPanel);
            }
          });
        }else{
          $(this).html("Confirm deletion");
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
  buildRouteBox: function(routeInfo){
    var box = "";
    box += "<div id = " + routeInfo.id +" class = 'route-description-box'>";
    box += "<span class = 'emphasis'>Name:</span><span class ='content'> " + routeInfo.name + "</span><br>";
    box += "<span class = 'emphasis'>Description: </span><span class ='content'>" + routeInfo.description + "</span><br>";
    box += "<span class = 'emphasis'>Enabled: </span><span class='content'>"+routeInfo.enabled + "</span><br>";
    box += "<span class = 'emphasis'>Color: </span><span class='content'>" + routeInfo.color + "</span><br>";
    box += "<span class = 'emphasis'>Time: </span><span class='content'>"+routeInfo.startTime + "-" + routeInfo.endTime + "</span><br>";
    box += "<span class = 'emphasis'>ID: </span><span class='content'>"+ routeInfo.id + "</span><br>";
    box += "<span class = 'emphasis'>Enabled:</span> <input class='enbox' routeId='"+routeInfo.id+"' id='" + routeInfo.id + "checkbox' type='checkbox'>enabled</input><br>";

    box += "<div style='float: right;width:auto;'><button class='button cbutton stops' routeId="+routeInfo.id +">Stops</button><button id='delete' routeId="+routeInfo.id +" class='button cbutton deleteroute'>Delete</button></div><br></div>";
    $(".routePanel").append(box);
    $("#" + routeInfo.id + "checkbox").prop("checked", routeInfo.enabled);

    $("#" + routeInfo.id + "checkbox").click(function(el){
      //el.target.checked;
      toSend = {id: el.target.getAttribute("routeid"), enabled: el.target.checked};

      $.ajax({
        url: "/routes/edit",
        type: "POST",
        data: JSON.stringify(toSend),
        contentType: "application/json",
        complete: function(data){
        }
      });
    });
  },
  buildSubmitBox: function(){
    var box = "";
    box += "<div style='padding-bottom: 30px;' class ='route-description-box'>";
    box += "<span class = 'emphasis'>Submit Route Json</span><br>";
    box += "<textarea class='json' style='width:100%; height: 100px;'></textarea>";

    box += "<button id='delete' style='float:right;' class='button cbutton addStopJson'>Add</button><br></div>";

    $(".routePanel").append(box);
  },

  initMap: function(){
    Routes.RoutesMap = L.map('mapid', {
      zoomControl: false,
      attributionControl: false
    });

    Routes.RoutesMap.setView([42.728172, -73.678803], 15.3);
    Routes.RoutesMap.addControl(L.control.attribution({
      position: 'bottomright',
      prefix: ''
    }));

    L.tileLayer('http://tile.stamen.com/toner-lite/{z}/{x}/{y}{r}.png', {
      attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
      minZoom: 13
    }).addTo(Routes.RoutesMap);

    Routes.RoutingControl = L.Routing.control({
      waypoints: [

      ],
      router: new L.Routing.mapbox(key),
      routeWhileDragging: true
    });

    Routes.RoutingControl.on('routeselected', function(e) {
      if (Routes.drawnRoute !== null){
        Routes.RoutesMap.removeLayer(Routes.drawnRoute);
      }
      Routes.drawnRoute = L.polyline(e.route.coordinates, {color: 'blue'});
      Routes.drawnRoute.addTo(Routes.RoutesMap);

    });
    Routes.RoutingControl.addTo(Routes.RoutesMap);
    Routes.RoutingWaypoints =[

    ];

    Routes.RoutesMap.on('click', function(e) {
      Routes.RoutingWaypoints.push(e.latlng);
      Routes.RoutingControl.setWaypoints(Routes.RoutingWaypoints);
    });

  },
  removeLastPoint: function(){
    Routes.RoutingWaypoints = Routes.RoutingWaypoints.slice(0, -1);
    Routes.RoutingControl.setWaypoints(Routes.RoutingWaypoints);
  },
  pullForm: function(){
    var coords = [];
    coords = Routes.drawnRoute.getLatLngs();

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
      Routes.submitForm($('#jsonField').val());
      console.log($('#jsonField').val());
      $('.prompt').css('display','none');
    },

    getJson: function(){
      var toSend = Routes.pullForm();
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
          $.get( "/routes", Routes.populateRoutesPanel);
        }
      });
    },

};

var Stops = {

};

var Admin = {
  StopMap: null,
  addStopMarker: null,
  StopsMap: null,

  toggleDisplay: function(){
    $(".vehicle-panel").toggle();
    $(".routePanel").toggle();
    $(".rtb").toggle();
  },
  bindStopButtons: function(){
    $(".stopDelete").click(function(e){
      var routeId = $(this).parent().attr("id");
      $.ajax({
        url: '/stops/' + $(this).parent().attr("stopid"),
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
    /*{"routeId":"asdf","name":"Test","description":"asdf","address":"aSDFASD","startTime":"asd","endTime":"asdf","enabled":"true","id":"","lat":"42.73074227479951","lng":"-73.67736339569092","toDelete":"false"}*/
    $(".stopSubmit").click(function(e){
      var toSend = {
        "id":"",
        "name":$(this).parent().find("#name").val(),
        "description":$(this).parent().find("#desc").val(),
        "startTime":"",
        "endTime":"",
        "routeId":$(this).parent().find("#route").find(":selected").val(),
        "enabled":"true",
        "toDelete":false,
        "lat":1,
        "lng":1,
      };
      console.log(toSend);
      $(this).parent().attr("lat", Admin.addStopMarker.getLatLng().lat);
      $(this).parent().attr("lng", Admin.addStopMarker.getLatLng().lng);
      if($(this).parent().attr("stopid") !== ""){
        toSend.toDelete = true;
      }
      toSend.id = $(this).parent().attr("stopid");
      toSend.lat = $(this).parent().attr("lat");
      toSend.lng = $(this).parent().attr("lng");

      Admin.submitStopForm(toSend);

    });
  },

  populateRoutesPanel: function(data){
    Routes.populateRoutesPanel(data);
  },

  initStopMap: function(){
    if(Admin.StopsMap === null){
      Admin.StopsMap = L.map('newStopMap', {
        zoomControl: false,
        attributionControl: false
      });

      Admin.StopsMap.setView([42.728172, -73.678803], 15.3);
      Admin.StopsMap.addControl(L.control.attribution({
        position: 'bottomright',
        prefix: ''
      }));

      L.tileLayer('https://stamen-tiles.a.ssl.fastly.net/toner-lite/{z}/{x}/{y}{r}.png', {
        attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
        minZoom: 13
      }).addTo(Admin.StopsMap);

      Admin.StopsMap.on( 'click', function(e){
        if(Admin.addStopMarker !== null){
          Admin.StopsMap.removeLayer(Admin.addStopMarker);
        }
        Admin.addStopMarker = L.marker(e.latlng,
          {
            draggable: true
          });
        Admin.addStopMarker.addTo(Admin.StopsMap);

      });

  }
},
  generateStopUIElement: function(data, i){
    tmp = "";
    box = "";
    box += "<div id=" + data[i].routeId + " lat="+data[i].lat+" lng="+data[i].lng+" stopid='"+data[i].id+"' class = 'route-description-box'>";
    box += "<span class = 'emphasis'>Name:</span><input id='name' type='text' value='" + data[i].name + "'></input><br>";
    box += "<span class = 'emphasis'>Description:</span><input id='desc' type='text' value='" + data[i].description + "'></input><br>";
    box += "<span class = 'emphasis'>Route:</span><select id='route'>";

    for (j = 0 ; j < Routes.RouteData.length; j++){
      if(Routes.RouteData[j].id == data[i].routeId){
        box += "<option value='"+ Routes.RouteData[j].id + "' selected>" + Routes.RouteData[j].name + "</option>";
      }else{
        tmp += "<option value='"+ Routes.RouteData[j].id + "'>" + Routes.RouteData[j].name + "</option>";
      }
      box += tmp;
    }

    box += "</select><br>";
    box += "<span class = 'emphasis'>Enabled:</span><input id='enabled' type='textbox' value="+data[i].enabled+"></input>";
    box += "<span class='button stopSubmit' style='float:right;'>submit</span><span class='button stopDelete' style='float:right;'>delete</span></div>";
    return box;
  },

  populateStopsForm: function(data,routeId){
    $(".stopPanel").html("");
    Admin.StopsMap= null;
    var tmp = "";
    var box = "";
    var j = 0;
    if (data === null){
      data = [];
    }
    for (var i = -1; i < data.length; i ++){
      if (i != -1 && data[i].routeId == routeId){
        $(".stopPanel").append(Admin.generateStopUIElement(data,i));
      }else if(i == -1){
        tmp = "";
        box = "";
        box += "<div id='' stopid='new' class = 'route-description-box'>";
        box += "<div id='newStopMap'style='height: 50%;position: inherit;width: 50%;background-color:black;z-index:0;border-style: solid; border-width:1px; border-color:black; float: inherit;'; background-color:black;z-index:0;'></div>";

        box += "<br><span class = 'emphasis'>Name:</span><input id='name' input type='text' value='New stop' ></input><br>";
        box += "<span class = 'emphasis'>Description:</span><input id='desc' type='text' value=></input><br>";
        box += "<span class = 'emphasis'>Route:</span><select id='route'>";

        for (j = 0 ; j < Routes.RouteData.length; j++){
          if(Routes.RouteData[j].id == routeId){
            box += "<option value='"+ Routes.RouteData[j].id + "' selected>" + Routes.RouteData[j].name + "</option>";
          }else{
            tmp += "<option value='"+ Routes.RouteData[j].id + "'>" + Routes.RouteData[j].name + "</option>";
          }
          box += tmp;
        }

        box += "</select><br>";
        box += "<span class = 'emphasis'>Enabled:</span><input id='enabled' type='textbox' value=></input>";
        box += "<span class='button stopSubmit' style='float:right;'>submit</span></div>";
        $(".stopPanel").append(box);
        Admin.initStopMap();
      }

    }
  },

  showMapPanel: function(){
    Admin.hideStopsPanel();
    $('.mapPanel').css('display','block');
    Routes.RoutesMap.invalidateSize();
  },

  hideMapPanel: function(){
    $('.mapPanel').css('display','none');
    Routes.RoutesMap.invalidateSize();
  },

  hideStopsPanel: function(){
    $('.stopPanel').css('display','none');
  },
  showStopsPanel: function(){
    Admin.hideMapPanel();
    $('.stopPanel').css('display','block');
    if(Admin.StopsMap !== null){
      Admin.StopsMap.invalidateSize();
    }

  },
  initMap: function(){
    Routes.initMap();
  },

  removeLastPoint: function(){
    Routes.removeLastPoint();
  },

  submitStopForm: function(toSend){
    $.ajax({
      url: "/stops/create",
      type: "POST",
      data: JSON.stringify(toSend),
      contentType: "application/json",
      complete: function(data){
        $.get( "/stops", function(e){
          Admin.populateStopsForm(e,toSend.routeId);
          Admin.bindStopButtons();

        });
      }
    });
  }

};
$(document).ready(function(){
  $.get("/getKey/", initializeAPIKeys);


});
