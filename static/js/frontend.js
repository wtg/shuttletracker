

Vue.component('shuttle-map',{
  template: `<div id="mapid" style="height: 100%; z-index:0; filter: invert(0)"></div>`,
  mounted(){
    this.initMap();
    this.grabStops();
    var a = setInterval(this.grabVehicles, 1000);
    var b = setInterval(this.grabRoutes, 1000);

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
      $.get( "/routes", this.updateRoutes);
    },

    updateRoutes: function(data){
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
      $.get( "/stops", this.updateStops);

    },

    updateStops: function(data){
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
      $.get( "/updates", this.updateVehicles);
    },

    updateVehicles: function(data){

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
            //console.log(data[j].color);
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
      $.get( "/vehicles", this.grabMessages);

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
      var nameToId = {};
      for(var i = 0; i < data.length; i ++){
        nameToId[data[i].vehicleName] = data[i].vehicleID;
      }
      var el = this;
      $.get( "/updates/message", function(data){
        for(var i = 0 ; i < data.length; i ++){

          var start_pos = data[i].indexOf('>') + 1;
          var end_pos = data[i].indexOf('<',start_pos);
          el.ShuttleMessages[nameToId[data[i].substring(start_pos,end_pos)]] = data[i];

        }
        el.updateMessages();
      });

    },

  }
});

Vue.component('dropdown-menu',{
  template: `
<div class="dropdown">
  <ul class="dropdown-main">
    <li class="dropdown-main-item">
      <a href="#" class="dropdown-icon">
        <img v-on:click="toggleDropdownMenuVisibility()" src="static/images/menu.svg"></button>
      </a>
      <ul class="dropdown-menu">
        <li class="dropdown-menu-item" id="dropdown-menu-item_shuttle-schedule">
          <p>Shuttle Schedules</p>
          <!-- http://www.rpi.edu/dept/parking/shuttle/ -->
          <ul class="dropdown-submenu" id="dropdown-submenu_shuttle-schedule">
            <li class="dropdown-submenu-item" id="dropdown-submenu-item_shuttle-schedule" v-for="item in list_data">
              <p><a target="_blank" rel="noopener noreferrer" :href="item.link">{{item.name}}</a></p>
            </li>
          </ul>
        </li>
        <li class="dropdown-menu-item" id="dropdown-menu-item_styling">
          <p>Styling</p>
          <!-- for changing the view of the page -->
          <ul class="dropdown-submenu" id="dropdown-submenu_styling">
            <li class="dropdown-submenu-item" id="dropdown-submenu-item_styling">
              <div id="darkmode-icon">
                <img src="static/images/moon.svg">
              </div>
            </li>
          </ul>
        </li>
      </ul>
    </li>
  </ul>
</div>
`,
  data (){
    return{
      list_data: [
          {name: "East: Monday-Thursday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleEastRoute.pdf"},
          {name: "East: Friday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyEastShuttleSchedule.pdf"},
          {name: "West: Monday-Thursday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleWestRoute.pdf"},
          {name: "West: Friday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyWestShuttleSchedule.pdf"},
          {name: "Weekend Late Night", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018Weekend-LateNightShuttleSchedule.pdf"}
      ],
        title: "RPI Shuttle Tracker",
        moonicon: "static/images/moon.svg",
        sunicon: "static/images/sun.svg",
        darkmodeOn: 0,
        dropdownMenuList: [],
        dropdownSubmenuList: [],

    };
  },
    methods: {
        toggleDarkmode: function () {
            if (this.darkmodeOn === 0) {
                this.darkmodeOn = 1;
                document.querySelector('div#darkmode-icon>img').src = this.sunicon;
                document.getElementById('mapid').style.filter = 'invert(1)';
            } else {
                this.darkmodeOn = 0;
                document.querySelector('div#darkmode-icon>img').src = this.moonicon;
                document.getElementById('mapid').style.filter = 'invert(0)';
            }
        },
        toggleDropdownMenuVisibility: function() {
          alert("TESTING");
            for (var i = 0; i < this.dropdownMenuList.length; i++) {
                if (this.dropdownMenuList[i].style.display === 'inline-block') {
                    this.dropdownMenuList[i].style.display = 'none';
                } else {
                    this.dropdownMenuList[i].style.display = 'inline-block';
                }
                if (this.dropdownSubmenuList[i].style.display === 'inline-block') {
                    this.dropdownSubmenuList[i].style.display = 'none';
                }
            }
        },
    }
});

Vue.component('title-bar', {
  template: `  
  <div class="titleBar">
    <!-- left side of tile bar -->
    <ul class="titleContent" id="titleContent-left">
      <li>
        <dropdown-menu></dropdown-menu>
      </li>
      <li>
        <p class="title">{{title}}</p>
      </li>
    </ul>
    <!-- right side of title bar -->
    <ul class="titleContent" id="titleContent-right">
      <li>
        <a href="https://webtech.union.rpi.edu" class="logo">
          <img src="static/images/wtg.svg">
        </a>
      </li>
    </ul>
  </div>`,
    mounted() {
        this.dropdownMenuList = document.getElementsByClassName('dropdown-menu');
        this.dropdownSubmenuList = document.getElementsByClassName('dropdown-submenu');
        console.log(document.getElementsByClassName('dropdown-menu').length);
        //document.querySelector('div#darkmode-icon').addEventListener('click', this.toggleDarkmode);
        //document.querySelector('a.dropdown-icon>img').addEventListener('click', this.toggleDropdownMenuVisibility);
        for (var i = 0; i < this.dropdownMenuList.length; i++) {
          for (var j = 0; j < this.dropdownMenuList[i].getElementsByClassName('dropdown-menu-item').length; j++) {
              this.dropdownMenuList[i].getElementsByClassName('dropdown-menu-item')[j].addEventListener('click', this.toggleDropdownSubmenuVisibility(i));
          }
        }
    },
    data (){
      return {
          title: "RPI Shuttle Tracker",
          moonicon: "static/images/moon.svg",
          sunicon: "static/images/sun.svg",
          darkmodeOn: 0,
          dropdownMenuList: [],
          dropdownSubmenuList: [],
      };
    },
    methods: {
        toggleDarkmode: function () {
          if (this.darkmodeOn === 0) {
            this.darkmodeOn = 1;
              document.querySelector('div#darkmode-icon>img').src = this.sunicon;
              document.getElementById('mapid').style.filter = 'invert(1)';
          } else {
              this.darkmodeOn = 0;
              document.querySelector('div#darkmode-icon>img').src = this.moonicon;
              document.getElementById('mapid').style.filter = 'invert(0)';
          }
        },

        toggleDropdownSubmenuVisibility: function(i) {
          var dropdownChildren = this.dropdownMenuList[i].getElementsByClassName('dropdown-submenu');
          for (var k = 0; k < dropdownChildren.length; k++) {
            if (dropdownChildren[k].style.display === 'inline-block') {
              dropdownChildren[k].style.display = 'none';
            } else {
              dropdownChildren[k].style.display = 'inline-block';
            }
          }
        },
    }

});

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {
  }

});
