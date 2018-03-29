
var live = true;
var partial = false;
var routeSuccess = true;
var stopsSuccess = true;
var a = true;
var lastUpdateTime = "";

var d = new Date();

function checkTime(i) {
  if (i < 10) {
    i = "0" + i;
  }
  return i;
}

Vue.component('live-indicator',{
  template: `<div id="live" v-bind:style="liveStyle">
    <p>{{lv + " " + text}}</p>
    <div v-if="live" class="pulsate" style=""></div></div>`,
  data (){
    return{
      liveStyle: {},
      text: "",
      lv:"Live",
      live: false
    };
  },
  methods: {
    update: function(){
      live = routeSuccess && vehicleSuccess;
      partial = routeSuccess || vehicleSuccess;
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

Vue.component('message-modal',{
  props: ['dim'],
  template:
  `<div id ="messagebox" @click="modalStyle.display = 'none'" v-if="msg != ''" v-bind:style="modalStyle">
  <div style="width: 100%;float:left;" v-html="msg"></div>
  <div style="position:absolute;right:10px;top:6px;color:#333;font-size:20px;">&times;</div>
  </div>`,
  data (){
    return {
      modalStyle: {display:"block"},
      msg: "",
    };
  },
  mounted (){
    let el = this;

    $.get("/adminMessage",function(data){
      if(data.Display === true){
        el.msg = data.Message;
      }
    });

  }
});

Vue.component('shuttle-map',{
  template: `
  <span>
  <div id="mapid" style="height: 100%; z-index:0; filter: invert(0)"></div>
  </div>
  <message-modal dim=100px></message-modal>
  </span>`,
  mounted(){
    this.initMap();
    this.grabStops();
    this.grabVehicles();
    var a = setInterval(this.grabVehicles, 3000);
    var b = setInterval(this.grabRoutes, 15000);

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
      legend: L.control({position: 'bottomleft'}),


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
      CircleSVG: `<?xml version="1.0"?>
      <svg height="600" width="600" xmlns="http://www.w3.org/2000/svg">
      <circle cx="50%" cy="50%" r="50%" fill="COLOR" />
      </svg>`
    };
  },
  methods:{
    getShuttleIcon: function(color){

      var url = "data:image/svg+xml;base64," + btoa(this.ShuttleSVG.replace("COLOR", color));
      return url;
    },

    getLegendIcon: function(color) {

		var url = "data:image/svg+xml;base64," + btoa(this.ShuttleSVG.replace("COLOR",color));
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

	updateLegend () {
	  let app = this;
	  app.legend.onAdd = function(map) {
		  var div = L.DomUtil.create('div','info legend');
		  var legendstring = "";
		  var darkModeVal = (document.querySelector('div.titleBar').style.filter === 'invert(0)') ? 0 : 1;
		    for (i = 0; i < app.ShuttleRoutes.length; i++){
			  let route = app.ShuttleRoutes[i];
			  legendstring += `<li><img class="legend-icon" src=` + app.getLegendIcon(route.color)+`
			  width="12" height="12" style="filter: invert(`+darkModeVal+`)"> `+
			  route.name;
		  }

		  div.innerHTML = `<ul style="list-style:none">
					<li><img class="legend-icon" src="static/images/user.svg" width="12" height="12" style="filter: invert(`+darkModeVal+`)"> You</li>`+
					legendstring +
					`<li><img class="legend-icon" src="static/images/circle.svg" width="12" height="12" style="filter: invert(`+darkModeVal+`)"> Shuttle Stop</li>
				</ul>`;
		return div;

		};
	  app.legend.addTo(app.ShuttleMap);
	},

    updateRoutes: function(data){
      routeSuccess = true;
      var updatedRoute = [];
      for(var i = 0; i < data.length; i ++){
        if(data[i].enabled === false || data[i].active === false){
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
      this.updateLegend();
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

    grabVehicles: function () {
      $.get( "/updates", this.updateVehicles).fail(function(){vehicleSuccess = false;});
    },

    updateVehicles: function (data) {
      window.d = new Date();
      var hours = window.d.getHours();
      hours = checkTime(hours);
      var mins = window.d.getMinutes();
      mins = checkTime(mins);
      var secs = window.d.getSeconds();
      secs = checkTime(secs);

      window.lastUpdateTime = (hours) + ":" + (mins) + ":" + (secs);
      vehicleSuccess = true;
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

          if (this.ShuttlesArray[data[j].vehicleID] === undefined){
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
          } else {
            shuttleIcon.options.iconUrl = this.getShuttleIcon(data[j].color);
            this.ShuttlesArray[data[j].vehicleID].data = data[j];
            this.ShuttlesArray[data[j].vehicleID].marker.setIcon(shuttleIcon);
            this.ShuttlesArray[data[j].vehicleID].marker.setLatLng([data[j].lat,data[j].lng]);
            this.ShuttlesArray[data[j].vehicleID].marker.setRotationAngle(parseInt(data[j].heading)-45);
          }
        }
      }

      this.setVehicleMessages();

      this.ShuttleUpdateCounter++;
    },

    cardinalDirection: function(heading) {
      if (heading >= 22.5 && heading < 67.5) {
        return "northeast";
      } else if (heading >= 67.5 && heading < 112.5) {
        return "east";
      } else if (heading >= 112.5 && heading < 157.5) {
        return "southeast";
      } else if (heading >= 157.5 && heading < 202.5) {
        return "south";
      } else if (heading >= 202.5 && heading < 247.5) {
        return "southwest";
      } else if (heading >= 247.5 && heading < 292.5) {
        return "west";
      } else if (heading >= 292.5 && heading < 337.5) {
        return "northwest";
      }
      return "north";
    },

    setVehicleMessages: function () {
      const v = this;
      $.get("/vehicles", function (data) {
        vehicleSuccess = true;
        let idToName = {};
        for (let i = 0; i < data.length; i++) {
          idToName[data[i].vehicleID] = data[i].vehicleName;
        }

        for (const key in v.ShuttlesArray) {

          const update = v.ShuttlesArray[key].data;
          const speed = Math.round(update.speed * 100) / 100;
          const date = new Date(update.created);
          const direction = v.cardinalDirection(parseFloat(update.heading));

          const message = `<b>${idToName[key]}</b><br>` +
          `Traveling ${direction} at ${speed} mph<br>` +
          `as of ${date.toLocaleTimeString()}`;
          v.ShuttlesArray[key].marker.bindPopup(message);
        }
      }).fail(function() {
        vehicleSuccess = false;
      });
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
  }
});

Vue.component('about-modal',{
  template:`
  <div class="modal">
  <div class="modalContent">
    <span class="close" @click="$emit('hidemodal')">&times;</span>
    <h2>About the Shuttle Tracker</h2>
    <p>
      The Shuttle Tracker is <a href="https://github.com/wtg/shuttletracker">open source</a> and maintained by the <a href="https://webtech.union.rpi.edu">Web Technologies Group</a> of the Rensselaer Union Student Senate for the benefit of the student body.
      It is also an active <a href="https://rcos.io">RCOS</a> project, and dozens of students have contributed to it through RCOS in previous semesters.
    </p>

    <h3>Why is there no app?</h3>
    <p>We're focused on making the Shuttle Tracker website as good as it can be. While most of our users are on mobile, we think the best way to deliver a great experience to the most people is to focus our efforts on this website.</p>

    <h3>How can I contribute?</h3>
    <p>We have an active <a href="https://github.com/wtg/shuttletracker">repository</a> on GitHub where we track issues and accept pull requests. If you want to get involved in another way, email <a href="mailto:webtech@union.lists.rpi.edu">webtech@union.lists.rpi.edu</a> for more information.</p>

    <!-- <h3>How are the shuttles tracked?</h3>
    <h3>What is the Web Technologies Group?</h3> -->

    </div>
  </div>
  `,
  methods: {
      toggleModal: function(){
        $(".modal").toggle();
      }
    }
});

Vue.component('dropdown-menu',{
  template: `
<div class="dropdown" >
  <ul class="dropdown-main">
    <li class="dropdown-main-item">
      <a v-on:click="toggleDropdownMenuVisibility()" class="dropdown-icon">
        <img class="dropdown-icon" src="static/images/menu.svg">
      </a>
      <ul class="dropdown-menu">
        <li class="dropdown-submenu-item" id="dropdown-submenu-item_styling">
        <p class="dropdown-menu-item_p">Information</p>
          <div class="dropdown-submenu-item_div" id="darkmode-icon">
            <p @click="$emit('displaymodal')" class="dropdown-submenu-subitem">About</p>
          </div>
        </li>
        <li class="dropdown-menu-item" id="dropdown-menu-item_shuttle-schedule">
          <p class="dropdown-menu-item_p">Shuttle Schedules</p>
          <!-- http://www.rpi.edu/dept/parking/shuttle/ -->
          <ul class="dropdown-submenu" id="dropdown-submenu_shuttle-schedule">
            <li class="dropdown-submenu-item" id="dropdown-submenu-item_shuttle-schedule" v-for="item in list_data">
              <p class="dropdown-submenu-item_p"><a class="dropdown-submenu-item_link" target="_blank" rel="noopener noreferrer" :href="item.link">{{item.name}}</a></p>
            </li>
          </ul>
        </li>
        <li class="dropdown-menu-item" id="dropdown-menu-item_styling">
          <p class="dropdown-menu-item_p">Styling</p>
          <!-- for changing the view of the page -->
          <ul class="dropdown-submenu" id="dropdown-submenu_styling">
            <li class="dropdown-submenu-item" id="dropdown-submenu-item_styling">
              <div class="dropdown-submenu-item_div" id="darkmode-icon">
                <img v-on:click="toggleDarkmode" class="dropdown-submenu-subitem" :src="moonicon">
                <p v-on:click="toggleDarkmode" class="dropdown-submenu-subitem">Darkmode</p>
              </div>
            </li>
          </ul>
        </li>
      </ul>
    </li>
  </ul>
</div>
`,
    mounted() {
        var vm = this;

        window.addEventListener('touchstart', function (event) {vm.dropdownWindowclick(event);});
        window.addEventListener('mousedown', function (event) {vm.dropdownWindowclick(event);});
        window.addEventListener('load', function () {vm.loadDarkmode();});

    },
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
        menuVisibility: 0,
        dropdownMenuList: document.getElementsByClassName('dropdown-menu'),
        dropdownSubmenuList: document.getElementsByClassName('dropdown-submenu'),

    };
  },
    methods: {
        // following functions involve changing the dark mode state
        toggleDarkmode: function () {
            // toggles dark mode for the map portion of the site by applying the 'filter: invert' property
            // For browsers that implement localStorage, read/write the dark mode toggle from localStorage
            // LocalStorage only stores keys and values as strings (Thanks JS!)
            if (typeof(Storage) !== "undefined") {
                if (localStorage.darkmodeOn === "0") {
                    this.enableDarkmode();
                } else if (localStorage.darkmodeOn === "1") {
                    this.disableDarkmode();
                } else {
                    localStorage.darkmodeOn = 1;
                    this.enableDarkmode();
                }
            } else {
                if (this.darkmodeOn === 0) {
                    this.enableDarkmode();
                } else {
                    this.disableDarkmode();
                }
            }
        },
        loadDarkmode: function () {
            // if localStorage is supported, read the value of darkmodeOn (if it exists),
            // then enable/disable dark mode as needed
            if (typeof(Storage) !== "undefined") {
                if (localStorage.darkmodeOn === "0") {
                    this.disableDarkmode();
                } else if (localStorage.darkmodeOn === "1") {
                    this.enableDarkmode();
                } else {
                    localStorage.darkmodeOn = 0;
                    this.disableDarkmode();
                }
            }
        },
        enableDarkmode: function () {
            if (typeof(Storage) !== "undefined") {
                localStorage.darkmodeOn = 1;
            } else {
                this.darkmodeOn = 1;
            }
            document.querySelector('div#darkmode-icon>img').src = this.sunicon;
            document.querySelector('div.leaflet-tile-pane').style.filter = 'invert(1)';
            document.querySelector('div.leaflet-bottom.leaflet-left').style.filter = 'invert(1)';
            document.querySelector('div.leaflet-bottom.leaflet-right').style.filter = 'invert(1)';
            document.querySelector('div.titleBar').style.filter = 'invert(1)';
            // invert specific colors twice to make them normal
            document.querySelector('ul#titleContent-right').style.filter = 'invert(1)';
            document.querySelector('ul#titleContent-right>li>div').style.filter = 'invert(1)';
            document.querySelector('div.pulsate').style.filter = 'invert(1)';
            var leafletControlLinks = document.querySelectorAll('div.leaflet-control>a');
            for (var i = 0; i < leafletControlLinks.length; i++) {
                leafletControlLinks[i].style.filter = 'invert(1)';
            }
            var legendIcons = document.querySelectorAll('img.legend-icon');
            for (var i = 0; i < legendIcons.length; i++) {
                legendIcons[i].style.filter = 'invert(1)';
            }
        },
        disableDarkmode: function () {
            if (typeof(Storage) !== "undefined") {
                localStorage.darkmodeOn = 0;
            } else {
                this.darkmodeOn = 0;
            }
            document.querySelector('div#darkmode-icon>img').src = this.moonicon;
            document.querySelector('div.leaflet-tile-pane').style.filter = 'invert(0)';
            document.querySelector('div.leaflet-bottom.leaflet-left').style.filter = 'invert(0)';
            document.querySelector('div.leaflet-bottom.leaflet-right').style.filter = 'invert(0)';
            document.querySelector('div.titleBar').style.filter = 'invert(0)';
            // reset specific colors to make normal
            document.querySelector('ul#titleContent-right').style.filter = 'invert(0)';
            document.querySelector('ul#titleContent-right>li>div').style.filter = 'invert(0)';
            document.querySelector('div.pulsate').style.filter = 'invert(0)';
            var leafletControlLinks = document.querySelectorAll('div.leaflet-control>a');
            for (var i = 0; i < leafletControlLinks.length; i++) {
                leafletControlLinks[i].style.filter = 'invert(0)';
            }
            var legendIcons = document.querySelectorAll('img.legend-icon');
            for (var i = 0; i < legendIcons.length; i++) {
                legendIcons[i].style.filter = 'invert(0)';
            }
        },
        // ====================================================================


        // following functions involve manipulating the dropdown menu
        dropdownWindowclick: function (event) {
            // check two parents up in case some don't have a class starting with dropdown
            if ((typeof event.target.className !== 'undefined'
                && event.target.className.substr(0, 8) === "dropdown")
                || (typeof event.target.parentElement.className !== 'undefined'
                && event.target.parentElement.className.substr(0, 8) === "dropdown")
                || (typeof event.target.parentElement.parentElement.className !== 'undefined'
                && event.target.parentElement.parentElement.className.substr(0, 8) === "dropdown"))
            {
                // menu was clicked so do nothing
            } else {
                this.closeDropdownMenu();
            }
        },
        toggleDropdownMenuVisibility: function() {
            // 0 is closed, 1 is open
            if (this.menuVisibility === 0) {
                this.openDropdownMenu();
            } else {
                this.closeDropdownMenu();
            }
        },
        toggleDropdownSubmenuVisibility: function() {
        },
        closeDropdownMenu: function () {
            this.menuVisibility = 0;
            this.disableDropdownMenuVisibility();
            this.disableDropdownSubmenuVisibility();
        },
        openDropdownMenu: function () {
            this.menuVisibility = 1;
            this.enableDropdownMenuVisibility();
            this.enableDropdownSubmenuVisibility();
        },
        disableDropdownMenuVisibility: function() {
            // close the menu by changing the visibility of all main dropdown menu entries
            for (var i = 0; i < this.dropdownMenuList.length; i++) {
                this.dropdownMenuList[i].style.display = 'none';
            }
        },
        disableDropdownSubmenuVisibility: function () {
            // close the menu by changing the visibility of all the submenu entries for each entry in the main menu
            for (var i = 0; i < this.dropdownSubmenuList.length; i++) {
                this.dropdownSubmenuList[i].style.display = 'none';
            }
        },
        enableDropdownMenuVisibility: function () {
            // close the menu by changing the visibility of all main dropdown menu entries
            for (var i = 0; i < this.dropdownMenuList.length; i++) {
                this.dropdownMenuList[i].style.display = 'inline-block';
            }
        },
        enableDropdownSubmenuVisibility: function () {
            // close the menu by changing the visibility of all the submenu entries for each entry in the main menu
            for (var i = 0; i < this.dropdownSubmenuList.length; i++) {
                this.dropdownSubmenuList[i].style.display = 'inline-block';
            }
        },
        // ====================================================================
    }
});

Vue.component('title-bar', {

  template: `
  <div class="titleBar" style="filter: invert(0)" >
    <!-- left side of tile bar -->
    <ul class="titleContent" id="titleContent-left">
      <li>
        <dropdown-menu @displaymodal="showModal=true"></dropdown-menu>
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
      <li>
        <live-indicator></live-indicator>
      </li>
    </ul>
    <about-modal v-if="showModal" @hidemodal="showModal=false"/>
  </div>`,
    mounted() {},
    data (){
      return {
          title: "RPI Shuttle Tracker",
          showModal: false,
      };
    },
    methods: {
    },

});

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {
  }

});
