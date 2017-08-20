var App ={
  ShuttleMap: null,
  ShuttleRoutes: [],
  Stops: {},
  Shuttles: {},
  MapBoundPoints: [],

  initMap: function(){
    App.ShuttleMap = L.map('mapid').setView([42.728172, -73.678803], 15.3);
    L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
      attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://mapbox.com">Mapbox</a>',
      maxZoom: 20,
      id: 'mapbox.streets',
      accessToken: 'pk.eyJ1Ijoiamx5b24xIiwiYSI6ImNqNmR4ZTVmejAwaTEzM3FsMmU0d2RmYjIifQ._VUaEMHioVwJIf11PzIqAQ'
    }).addTo(App.ShuttleMap);
    App.grabRoutes();
  },

  grabRoutes: function(){
    $.get( "http://127.0.0.1:8080/routes", App.updateRoutes);
  },

  updateRoutes: function(data){
    console.log(data);
    for(var i = 0; i < data.length; i ++){
      var points = []
      for(var j = 0; j < data[i]['coords'].length; j ++){
        points.push(new L.LatLng(data[i]['coords'][j]['lat'],data[i]['coords'][j]['lng']));
      }

      var polylineOptions = {
        color: data[i]['color'],
        weight: data[i]['width'],
        opacity: 1
      };

      var polyline = new L.Polyline(points, polylineOptions);

      var r ={
        name: data[i]['name'],
        id: data[i]['id'],
        description: data[i]['description'],
        color: data[i]['color'],
        created: data[i]['created'],
        enabled: data[i]['enabled'],
        stops: data[i]['stopsid'],
        start_time: data[i]['startTime'],
        end_time: data[i]['endTime'],
        points: points,
        line: polyline
      };

      App.ShuttleRoutes.push(r);

    }
    App.drawRoutes();

  },

  drawRoutes: function(){
    for(var i = 0; i < App.ShuttleRoutes.length; i ++){
      for(var j = 0; j < App.ShuttleRoutes[i]['points'].length; j ++){
        App.MapBoundPoints.push(App.ShuttleRoutes[i]['points'][j]);
      }
    }
    var polylineOptions = {
      color: 'blue',
      weight: 1,
      opacity: 1
    };
    var polyline = new L.Polyline(App.MapBoundPoints, polylineOptions);
    App.ShuttleMap.fitBounds(polyline.getBounds());

    for(var i = 0; i < App.ShuttleRoutes.length; i ++){
      App.ShuttleMap.addLayer(App.ShuttleRoutes[i]['line']);
    }

  }

}

$(document).ready(function(){
  App.initMap();

});
