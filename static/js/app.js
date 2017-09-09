var ShuttlesArray = {};
var ShuttleMessages = {};

var App ={
  ShuttleMap: null,
  ShuttleRoutes: [],
  Stops: [],
  Shuttles: {},
  MapBoundPoints: [],
  ShuttleUpdateCounter: 0,

  initMap: function(){
    App.ShuttleMap = L.map('mapid', {
        zoomControl: false,
        attributionControl: false // hide Leaflet
    })
    App.ShuttleMap.setView([42.728172, -73.678803], 15.3);
    // show attribution without Leaflet
    App.ShuttleMap.addControl(L.control.attribution({
        position: 'bottomright',
        prefix: ''
    }));
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
      maxZoom: 19
    }).addTo(App.ShuttleMap);
    App.grabRoutes();
  },

  grabRoutes: function(){
    $.get( "/routes", App.updateRoutes);
  },

  updateRoutes: function(data){
    var updatedRoute = [];
    for(var i = 0; i < data.length; i ++){
      var points = [];
      for(var j = 0; j < data[i].coords.length; j ++){
        points.push(new L.LatLng(data[i].coords[j].lat,data[i].coords[j].lng));
      }
      var polylineOptions = {
        color: data[i].color,
        weight: 3,//data[i]['width'],
        opacity: 1
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

    App.ShuttleRoutes = updatedRoute;
    App.drawRoutes();

  },


  drawRoutes: function(){
    for(i = 0; i < App.ShuttleRoutes.length; i ++){
      App.ShuttleMap.removeLayer(App.ShuttleRoutes[i].line)
    }
    for(i = 0; i < App.ShuttleRoutes.length; i ++){
      for(var j = 0; j < App.ShuttleRoutes[i].points.length; j ++){
        App.MapBoundPoints.push(App.ShuttleRoutes[i].points[j]);
      }
    }

    var polylineOptions = {
      color: 'blue',
      weight: 1,
      opacity: 1
    };
    var polyline = new L.Polyline(App.MapBoundPoints, polylineOptions);
    App.ShuttleMap.fitBounds(polyline.getBounds());
    for(i = 0; i < App.ShuttleRoutes.length; i ++){
      App.ShuttleMap.addLayer(App.ShuttleRoutes[i].line);
    }

  },

  grabStops: function(){
    $.get( "/stops", App.updateStops);

  },

  updateStops: function(data){
    var stopIcon = L.icon({
      iconUrl: 'static/images/stop.png',

      iconSize:     [8, 8], // size of the icon
      iconAnchor:   [4, 4], // point of the icon which will correspond to marker's location
      shadowAnchor: [4, 62],  // the same for the shadow
      popupAnchor:  [0, 0] // point from which the popup should open relative to the iconAnchor
    });
    for(var i = 0; i < data.length; i ++){
      var stop = {
        name: data[i].name,
        description: data[i].description,
        id: data[i].id,
        latlng: [data[i].lat, data[i].lng],
        marker: L.marker([data[i].lat,data[i].lng], {icon: stopIcon})
      };
      stop['marker'].bindPopup(stop.name);
      stop['marker'].addTo(App.ShuttleMap).on('click', App.stopClicked);
    }

  },

  grabVehicles: function(){
    $.get( "/updates", App.updateVehicles);
  },

  updateVehicles: function(data){
    //console.log(data.length + " shuttles updated");
    var shuttleIcon = L.icon({
      iconUrl: 'static/images/shuttle_arrow.png',

      iconSize:     [14, 14], // size of the icon
      iconAnchor:   [7, 7], // point of the icon which will correspond to marker's location
      popupAnchor:  [0, 0] // point from which the popup should open relative to the iconAnchor
    });

    if(App.ShuttleUpdateCounter >= 15){
      for (var key in ShuttlesArray){
        var good = false;
        for(var i = 0; i < data.length; i ++){
          if(key == data[i].vehicleID){
            good = true;
          }
        }
        if(good === false && ShuttlesArray[key] !== null) {
          App.ShuttleMap.removeLayer(ShuttlesArray[key].marker);
          ShuttlesArray[key] = null;
        }

      }

      App.ShuttleUpdateCounter = 0;
    }
    if(data !== null){
      for(var j = 0; j < data.length; j ++){
        //console.log(parseInt(data[i]['heading']));
        if(ShuttlesArray[data[j].vehicleID] == null){
          ShuttlesArray[data[j].vehicleID] = {
            data: data[j],
            marker: L.marker([data[j].lat,data[j].lng], {icon: shuttleIcon, rotationAngle: parseInt(data[j].heading)-45,rotationOrigin: 'center'}),
            message: ""
          };
          ShuttlesArray[data[j].vehicleID].marker.addTo(App.ShuttleMap);
        }else{
          ShuttlesArray[data[j].vehicleID].marker.setLatLng([data[j].lat,data[j].lng]);
          ShuttlesArray[data[j].vehicleID].marker.setRotationAngle(parseInt(data[j].heading)-45);
        }
      }
    }
    App.ShuttleUpdateCounter ++;
    App.grabVehicleInfo();

  },

  stopClicked: function(e){
    App.ShuttleMap.setView(e.target.getLatLng(),16);
  },

  grabVehicleInfo: function(){
    $.get( "/vehicles", App.grabMessages);

  },

  grabMessages: function(data){
    var nameToId = {};
    for(var i = 0; i < data.length; i ++){
      nameToId[data[i].vehicleName] = data[i].vehicleID;
    }
    $.get( "/updates/message", function(data){
      for(var i = 0 ; i < data.length; i ++){
        ShuttleMessages[nameToId[data[i].substring(3,9)]] = data[i];
      }
      App.updateMessages();
    });

  },

  updateMessages: function(){
    for(var key in ShuttlesArray){
      for(var messageKey in ShuttleMessages){
        if(key == messageKey && ShuttlesArray[key] !== null){
          ShuttlesArray[key].marker.bindPopup(ShuttleMessages[messageKey]);
        }
      }
    }

  }

};

$(document).ready(function(){
  App.initMap();
  App.grabStops();
  var a = setInterval(App.grabVehicles, 1000);
});
