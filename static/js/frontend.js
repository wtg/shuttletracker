
var live = true;
var partial = false;
var routeSuccess = true;
var stopsSuccess = true;
var vehicleUpdateSuccess = true;
var vehicleMessageSuccess = true;
var lastUpdateTime = "";

var d = new Date();

function checkTime(i) {
    if (i < 10) {
        i = "0" + i;
    }
    return i;
}

Vue.component('live-indicator',{
  template: `<div v-bind:style="liveStyle">{{lv}} <div v-if="live" class="pulsate" style="position:absolute;float:right; width:10px;height:10px;background-color:blue;border-radius:50%;top:9px;right:3px;"></div>{{text}}</div>`,
  data (){
    return{
      liveStyle: {color:"black",width: "40px", height:"18px",padding:"5px",borderRadius:"5px",fontSize:"15px", backgroundColor:"rgba(255, 255, 255, 0.88)", boxShadow: "0 1px 1px rgba(0, 0, 0, 0.8)", display:"none", position: "absolute", right:"10px",top:"42px"},
      text: "",
      lv:"Live",
      live: false
      };
  },
  methods: {
    update: function(){
      live = routeSuccess && vehicleUpdateSuccess && vehicleMessageSuccess;
      partial = routeSuccess || vehicleUpdateSuccess;
      if(live === false){
        this.text = window.lastUpdateTime;
        this.liveStyle.width="auto";
        this.live=false;
        this.lv = "Last Updated";
      }else{
        //this.text="Live";
        this.live = true;
        this.lv = "Live";
        this.text = "";
        this.liveStyle.width = "40px";
        this.liveStyle.display = "inline-block";


      }
    }
  },
  mounted (){
    setInterval(this.update,1000);
  }

});

Vue.component('shuttle-map',{
  template: `<div id="mapid" style="height: 100%; z-index:0;"></div>`,
  mounted(){
    this.initMap();
    this.grabStops();
    var a = setInterval(this.grabVehicles, 3000);
    var b = setInterval(this.grabRoutes, 3000);

  },
  data (){
    return{
      ShuttlesArray: {},
      ShuttleMessages: {},
      ShuttleMap: null,
      ShuttleRoutes: [],
      Stops: [],
      Shuttles: {},
      MapBoundPoints: [],
      ShuttleUpdateCounter: 0,
      first: true,

      ShuttleSVG: `<?xml version="1.0" encoding="UTF-8"?>
          <svg width="52px" height="52px" viewBox="0 0 52 52" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
              <title>shuttle</title>
              <defs></defs>
              <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
                  <g id="shuttle">
                      <path d="M51.353,0.914 C51.648,1.218 51.72,1.675 51.532,2.054 L27.532,50.469 C27.362,50.814 27.011,51.025 26.636,51.025 C26.58,51.025 26.524,51.02 26.467,51.01 C26.032,50.936 25.697,50.583 25.643,50.145 L23.098,29.107 L0.835,25.376 C0.402,25.304 0.067,24.958 0.009,24.522 C-0.049,24.086 0.184,23.665 0.583,23.481 L50.218,0.701 C50.603,0.524 51.058,0.609 51.353,0.914 Z" id="Background" fill="COLOR"></path>
                      <path d="M51.353,0.914 C51.058,0.609 50.603,0.524 50.218,0.701 L0.583,23.481 C0.184,23.665 -0.049,24.086 0.009,24.522 C0.067,24.958 0.402,25.304 0.835,25.376 L23.098,29.107 L25.643,50.145 C25.697,50.583 26.032,50.936 26.467,51.01 C26.524,51.02 26.58,51.025 26.636,51.025 C27.011,51.025 27.362,50.814 27.532,50.469 L51.532,2.054 C51.72,1.675 51.648,1.218 51.353,0.914 Z M27.226,46.582 L24.994,28.125 C24.94,27.685 24.603,27.332 24.166,27.259 L4.374,23.941 L48.485,3.697 L27.226,46.582 Z" id="Shape" fill="#000"></path>
                  </g>
              </g>
          </svg>
          `,
    };
  },
  methods:{
    getShuttleIcon: function(color){

      var url = "data:image/svg+xml;base64," + btoa(this.ShuttleSVG.replace("COLOR", color));
      return url;
    },

    initMap: function(){
      this.ShuttleMap = L.map('mapid', {
          zoomControl: false,
          attributionControl: false // hide Leaflet
      });
      this.ShuttleMap.setView([42.728172, -73.678803], 15.3);
      // show attribution without Leaflet
      this.ShuttleMap.addControl(L.control.attribution({
          position: 'bottomright',
          prefix: ''
      }));
      L.tileLayer('https://stamen-tiles.a.ssl.fastly.net/toner-lite/{z}/{x}/{y}{r}.png', {
        attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
        maxZoom: 17,
        minZoom: 14
      }).addTo(this.ShuttleMap);
      this.grabRoutes();
      this.showUserLocation(this.ShuttleMap);

    },

    grabRoutes: function(){
      $.get( "/routes", this.updateRoutes).fail(function(){routeSuccess = false;});
    },

    updateRoutes: function(data){
      routeSuccess = true;
      var updatedRoute = [];
      for(var i = 0; i < data.length; i ++){
        if(data[i].enabled === false){
          continue;
        }
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
      for(i = 0; i < this.ShuttleRoutes.length; i ++){
        this.ShuttleMap.removeLayer(this.ShuttleRoutes[i].line);
      }
      this.ShuttleRoutes = updatedRoute;
      this.drawRoutes();

    },


    drawRoutes: function(){
      for(i = 0; i < this.ShuttleRoutes.length; i ++){
        this.ShuttleMap.removeLayer(this.ShuttleRoutes[i].line);
      }
      if(this.first){
        for(i = 0; i < this.ShuttleRoutes.length; i ++){
          for(var j = 0; j < this.ShuttleRoutes[i].points.length; j ++){
            this.MapBoundPoints.push(this.ShuttleRoutes[i].points[j]);
          }
        }

        var polylineOptions = {
          color: 'blue',
          weight: 1,
          opacity: 1
        };
        var polyline = new L.Polyline(this.MapBoundPoints, polylineOptions);
        this.ShuttleMap.fitBounds(polyline.getBounds());
        this.first = false;
      }
      for(i = 0; i < this.ShuttleRoutes.length; i ++){
        this.ShuttleMap.addLayer(this.ShuttleRoutes[i].line);
      }

    },

    grabStops: function(){
      $.get( "/stops", this.updateStops).fail(function(){stopsSuccess = false;});

    },

    updateStops: function(data){
      stopsSuccess = true;
      var stopIcon = L.icon({
        iconUrl: 'static/images/circle.svg',

        iconSize:     [12, 12], // size of the icon
        iconAnchor:   [6, 6], // point of the icon which will correspond to marker's location
        shadowAnchor: [6, 6],  // the same for the shadow
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
        stop.marker.bindPopup(stop.name);
        stop.marker.addTo(this.ShuttleMap).on('click', this.stopClicked);
      }

    },

    grabVehicles: function(){
      $.get( "/updates", this.updateVehicles).fail(function(){vehicleUpdateSuccess = false;});
    },

    updateVehicles: function(data){
      window.d = new Date();
      var hours = window.d.getHours();
      hours = checkTime(hours);
      var mins = window.d.getMinutes();
      mins = checkTime(mins);
      var secs = window.d.getSeconds();
      secs = checkTime(secs);

      window.lastUpdateTime = (hours) + ":" + (mins) + ":" + (secs);
      vehicleUpdateSuccess = true;
      var shuttleIcon = L.icon({
        iconUrl: this.getShuttleIcon("#FFF"),

        iconSize:     [32, 32], // size of the icon
        iconAnchor:   [16, 16], // point of the icon which will correspond to marker's location
        popupAnchor:  [0, 0] // point from which the popup should open relative to the iconAnchor
      });

      if(this.ShuttleUpdateCounter >= 15){
        for (var key in this.ShuttlesArray){
          var good = false;
          for(var i = 0; i < data.length; i ++){
            if(key == data[i].vehicleID){
              good = true;
            }
          }
          if(good === false && this.ShuttlesArray[key] !== null) {
            this.ShuttleMap.removeLayer(this.ShuttlesArray[key].marker);
            this.ShuttlesArray[key] = null;
          }

        }

        this.ShuttleUpdateCounter = 0;
      }

      if(data !== null){
        for(var j = 0; j < data.length; j ++){
          for(var k = 0; k < this.ShuttleRoutes.length; k ++){
            if (this.ShuttleRoutes[k].id === data[j].RouteID){
              data[j].color = this.ShuttleRoutes[k].color;
              break;
            }
          }
          if(data[j].color === undefined){
            data[j].color = "#FFF";
          }

          if(this.ShuttlesArray[data[j].vehicleID] === undefined){
            shuttleIcon.options.iconUrl = this.getShuttleIcon(data[j].color);
            this.ShuttlesArray[data[j].vehicleID] = {
              data: data[j],
              marker: L.marker([data[j].lat,data[j].lng], {
                  icon: shuttleIcon,
                  rotationAngle: parseInt(data[j].heading)-45,rotationOrigin: 'center',
                  zIndexOffset: 1000
              }),
              message: ""
            };
            this.ShuttlesArray[data[j].vehicleID].marker.addTo(this.ShuttleMap);
          }else{
            shuttleIcon.options.iconUrl = this.getShuttleIcon(data[j].color);
            this.ShuttlesArray[data[j].vehicleID].marker.setIcon(shuttleIcon);
            this.ShuttlesArray[data[j].vehicleID].marker.setLatLng([data[j].lat,data[j].lng]);
            this.ShuttlesArray[data[j].vehicleID].marker.setRotationAngle(parseInt(data[j].heading)-45);
          }
        }
      }
      this.ShuttleUpdateCounter ++;
      this.grabVehicleInfo();

    },

    showUserLocation: function(map){
      var userIcon = L.icon({
        iconUrl: 'static/images/user.svg',

        iconSize:     [12, 12], // size of the icon
        iconAnchor:   [6, 6], // point of the icon which will correspond to marker's location
        shadowAnchor: [6, 6],  // the same for the shadow
        popupAnchor:  [0, 0] // point from which the popup should open relative to the iconAnchor
      });

       if (navigator.geolocation) {
          navigator.geolocation.getCurrentPosition(showPosition);
        } else {
          console.log("Geolocation is either not supported by this browser, or geolocation permissions were not given by the user.");
        }

      function showPosition (position) {
        var locationMarker = {
              name: "You are here",
              marker: L.marker([position.coords.latitude, position.coords.longitude], {
                  icon: userIcon,
                  zIndexOffset: 1000
              }),
        };
        locationMarker.marker.bindPopup(locationMarker.name);

        locationMarker.marker.addTo(map);
      }
    },

    stopClicked: function(e){
    },

    grabVehicleInfo: function(){
      $.get( "/vehicles", this.grabMessages).fail(function(){vehicleUpdateSuccess = false;});

    },
    updateMessages: function(){
      for(var key in this.ShuttlesArray){
        for(var messageKey in this.ShuttleMessages){
          if(key == messageKey && this.ShuttlesArray[key] !== null){
            this.ShuttlesArray[key].marker.bindPopup(this.ShuttleMessages[messageKey]);
          }
        }
      }

    },

    grabMessages: function(data){
      vehicleUpdateSuccess = true;
      var nameToId = {};
      for(var i = 0; i < data.length; i ++){
        nameToId[data[i].vehicleName] = data[i].vehicleID;
      }
      var el = this;
      $.get( "/updates/message", function(data){
        vehicleMessageSuccess = true;
        for(var i = 0 ; i < data.length; i ++){

          var start_pos = data[i].indexOf('>') + 1;
          var end_pos = data[i].indexOf('<',start_pos);
          el.ShuttleMessages[nameToId[data[i].substring(start_pos,end_pos)]] = data[i];

        }
        el.updateMessages();
      }).fail(function(){vehicleMessageSuccess = false;});

    },

  }
});

Vue.component('dropdown-menu',{
  template: `<div>
      <ul class="dropdown">
          <li>
              <a href="#" class="dropdown-schedule">
                  <img src="static/images/menu.svg">
              </a>
              <ul class="dropdown-menu">
                  <li>
                      <p id="schedule-menu">Shuttle Schedules</p>
                      <!-- http://www.rpi.edu/dept/parking/shuttle/ -->
                  </li>
                  <li v-for="item in list_data">
                      <p><a target="_blank" rel="noopener noreferrer" :href="item.link">{{item.name}}</a></p>
                  </li>
              </ul>
          </li>
      </ul>
  </div>`,
  data (){
    return{
        list_data: [
          {name: "East: Monday-Thursday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleEastRoute.pdf"},
          {name: "East: Friday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyEastShuttleSchedule.pdf"},
          {name: "West: Monday-Thursday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleWestRoute.pdf"},
          {name: "West: Friday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyWestShuttleSchedule.pdf"},
          {name: "Weekend Late Night", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018Weekend-LateNightShuttleSchedule.pdf"}
        ],

    };
  }
});

Vue.component('title-bar', {
  template:
  `  <div class="titleBar">
        <div class="titleContent">
            <dropdown-menu></dropdown-menu>
            <p class="title">{{title}}</p>
            <live-indicator></live-indicator>
            <a href="https://webtech.union.rpi.edu" class="logo">
                <img src="static/images/wtg.svg">
            </a>
        </div>
    </div>`,
    data (){
      return {
        title: "RPI Shuttle Tracker"
      };
    }

});

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {
  }

});
