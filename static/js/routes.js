var refresh = false;

Vue.component('active-selector',{
  props: ['routeId'],
  template: `<div>
  <select class="dropdown" v-model=selected>
    <option v-for="day in days" v-bind:value=day.val>{{day.Day}}</option>
  </select>
  <input class="field" type="time" v-model=tod></input>
  <input class="button" type="button" @click="on = !on" v-model="on"></input>
  <input class="button" type="button" value="+" @click="add"></input>
  <input class="button" type="button" value="Submit" @click="submit"></input>
  <input class="button" type="button" value="Refresh" @click="get"></input>
  <br>
  <div class = "tile is-parent columns" v-for="item in times">
  <div class = "column notification is-child is-light"><div class="column has-text-centered">Day: {{item.day}} Time: {{item.time.getHours()}}:<span v-if="(item.time.getMinutes() < 10)">0</span>{{item.time.getMinutes()}} Active: {{item.state}}</div><input class="button is-fullwidth" type="button" value="delete" @click="del(item.id)"></input></div>
  </div>
  </div>
  <br>
  </div>`,
  data (){
    return {
      selected: 1,
      times: [],
      on: false,
      tod: "20:00",
      days: {Monday: {Day: "Monday",val: 1},Tuesday: {Day: "Tuesday",val: 2},Wednesday: {Day: "Wednesday",val: 3},Thursday: {Day: "Thursday",val: 4},Friday: {Day: "Friday",val: 5},Saturday: {Day: "Saturday",val: 6},Sunday: {Day: "Sunday",val: 0}},

    };
  },
  mounted (){
    this.get();
  },
  methods: {
    add: function(){
      let o = (this.on) ? 1 : 0;
      split = this.tod.split(":");
      ints = [];
      for (let i = 0; i < split.length; i ++){
        ints[i] = parseInt(split[i]);
      }
      let i = this.times.length;
      let obj = {id: i, state: o, day: this.selected, time: (new Date(1,1,1,ints[0],ints[1],0,0))};
      this.times.push(obj);
    },
    del: function(id){
      if(this.times !== undefined){
        for (let i =0; i < this.times.length; i ++){
          if(this.times[i].id == id){
             this.times.splice(i, 1);
            break;
          }
        }
      }
    },
    submit: function(){
      let el = this;
      var send = {id: el.routeId, times: []};
      send.times = el.times;
      $.post("/routes/schedule",JSON.stringify(send),function(data){});
    },
    get: function(){
      let el = this;
      el.times = [];
      $.get("/routes",function(data){
        for(var i = 0; i < data.length; i ++){
          if(data[i].id == el.routeId){
            for (var j = 0; j < data[i].intervals.length; j ++){
              var obj = {id: j, state: data[i].intervals[j].State, day: data[i].intervals[j].Day, time: new Date(data[i].intervals[j].Time)};
              el.times.push(obj);
            }

          }
        }
      });
    }
  }
});

Vue.component('route-panel', {
  template:
  `<div class ="column">
    <div v-for="route in routeData" class="tile is-parent">
      <div class="tile is-parent">
        <route-card v-bind:info="route"></route-card>
      </div>
    </div>
    <route-create></route-create>
    <route-json></route-json>
  </div>`,
  data (){
    return{
      routeData: [{shuttleID:22}]
    };
  },
  mounted(){
    var el = this;
    $.get("/routes",function(data){
      el.routeData = data;
      refresh = false;
    });
    setInterval(function(){
      if(refresh){
        $.get("/routes",function(data){
          el.routeData = data;
          refresh = false;
        });
      }
    },100);

  }

});
Vue.component('route-card', {
  props: ['info'],
  template:
  `
  <div class="box container">
  <div class="tabs">
  <ul>
  <li v-bind:class="{'is-active':state==0}"><a @click="state = 0">Route</a></li>
  <li v-bind:class="{'is-active':state==1}"><a @click="state = 1">Schedule</a></li>
  </ul>
  </div>
  <div >
    <div class="columns">
    <div v-if="state==0" class="column">
    <b>id</b>: {{info.id}}<br>
    <b>name</b>: {{info.name}}<br>
    <b>Description</b>: {{info.description}}<br>
    <b>enabled</b>: <input type="checkbox" v-model="info.enabled" @click="toggleRoute(info.id)"></input>{{info.enabled}}<br>
    <b>Color</b>: {{info.color}}<br>
    <b>Created</b>: {{info.created}} <br>
    <br>
    <button class="button" @click="shouldDelete(info.id)">{{buttonText}}</button>
    </div>
    <div v-if="state==1" class="column">
    <active-selector :routeId=info.id></active-selector>
    </div>
    </div>
  </div></div>`,
  data (){
    return{
      myData: {},
      deleteCount: 0,
      buttonText: "Delete",
      state: 0,
    };
  },
  methods: {
    shouldDelete: function(id){
      if(this.deleteCount == 1){
        this.deleteRoute(id);
      }else{
        this.buttonText = "Are you sure?";
        this.deleteCount++;
      }
    },
    toggleRoute: function(id){
      data = {id: id, enabled: !this.info.enabled};
      $.ajax({
        url: "/routes/edit",
        type: "POST",
        data: JSON.stringify(data),
        contentType: "application/json",
        complete: function(data){
        }
      });
    },

    deleteRoute: function(id){
      $.ajax({
           url: '/routes?id=' + id,
           type: 'DELETE',
           success: function(result) {
             refresh = true;
           }
         });
    }

  }

});

Vue.component('route-json',{
  template:`
  <div class="tile box">
      <span class = 'emphasis'>Submit Route Json</span><br>
      <textarea id="jsontxt" class='json' style='width:100%; height: 100px;'></textarea>
      <button id='submitRouteJson' @click="submitForm" style='float:right;' class='button cbutton addStopJson'>Add</button><br></div>
    </div>`,
    data (){
      return{

      };
    },
    methods: {
      submitForm: function(){
        toSend = $("#jsontxt").val();
        $.ajax({
          url: "/routes/create",
          type: "POST",
          data: toSend,
          contentType: "application/json",
          complete: function(data){
            $.get( "/routes", this.populateRoutesPanel);
          }
        });
      },
    }
});


Vue.component('route-create',{
  template:`
  <div><div class="route-description-box box" style="padding-bottom: 20px;" >
              <div id="mapid" style="height: 650px;float: left; width: 100%; background-color:black;z-index:0;"></div>
              <div class="mapcontrols"><button class="button" @click="removeLastPoint">undo</button></div>
              <b>name</b>: <input v-model="name" placeholder="Route Name"></input><br>
              <b>Description</b>: <input v-model="description" placeholder="Route Description"></input><br>
              <b>Color</b>: <input v-bind:style="{ backgroundColor: color}" v-model="color" placeholder="#FFFF00"></input><br>
              <b>Enabled</b>:  <input type="checkbox" v-model="enabled"></input>{{enabled}}<br>
              <b>Width</b>: <input v-model="width" placeholder="4"></input><br>
            <button class="button" @click="submitForm(JSON.stringify(getJSON()))">Submit</button>
            <button class="button" @click="showJSON(JSON.stringify(getJSON()))">GetJSON</button>
    </div></div>`,
    data (){
      return{
        ID: "",
        name: "",
        description: "",
        color: "",
        width: 4,
        active: true,
        enabled: true,
        RoutesMap: null,
        APIKey: null,
        RoutingControl: null,
        drawnRoute: null,
        RoutingWaypoints: [],

      };
    },
    methods: {
      createPoints: function (route) {
        points = [];
        for (const p of route.getLatLngs()) {
          points.push({latitude: p.lat, longitude: p.lng});
        }
        return points;
      },
      getJSON: function(){
        coords = this.createPoints(this.drawnRoute);
        var toSend = {
          "name":this.name,
          "description":this.description,
          "enabled":this.enabled,
          "color":this.color,
          "width":this.width,
          "points": coords
        };
          return toSend;
      },
      showJSON: function(data){
        var wnd = window.open("about:blank", "", "_blank");
        wnd.document.write(data);

      },
      submitForm: function(toSend){
        $.ajax({
          url: "/routes/create",
          type: "POST",
          data: toSend,
          contentType: "application/json",
          complete: function(data){
            $.get( "/routes", this.populateRoutesPanel);
          }
        });
      },

      removeLastPoint: function(){
        this.RoutingWaypoints = this.RoutingWaypoints.slice(0, -1);
        this.RoutingControl.setWaypoints(this.RoutingWaypoints);
      },
      send: function(){

        var pkg = {"vehicleID":this.ID, "vehicleName":this.name, "enabled":this.enabled};

        pkg = JSON.stringify(pkg);
        $.ajax({
          url: "/vehicles/create",
          type: "POST",
          data: pkg,
          contentType: "application/json",
          success: function(data){
            refresh = true;
          }
        });
        this.id = "";
        this.name = "";
      },
      initMap: function(){
        let el = this;
        this.RoutesMap = L.map('mapid', {
          zoomControl: false,
          attributionControl: false
        });

        this.RoutesMap.setView([42.728172, -73.678803], 15.3);
        this.RoutesMap.addControl(L.control.attribution({
          position: 'bottomright',
          prefix: ''
        }));
        L.tileLayer('https://stamen-tiles.a.ssl.fastly.net/toner-lite/{z}/{x}/{y}{r}.png', {
          attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
          minZoom: 13
        }).addTo(this.RoutesMap);

        this.RoutingControl = L.Routing.control({
          waypoints: [

          ],
          router: new L.Routing.mapbox(this.APIKey),
          routeWhileDragging: true
        });

        this.RoutingControl.on('routeselected', function(e) {
          if (el.drawnRoute !== null){
            el.RoutesMap.removeLayer(el.drawnRoute);
          }
          el.drawnRoute = L.polyline(e.route.coordinates, {color: 'blue'});
          el.drawnRoute.addTo(el.RoutesMap);

        });
        this.RoutingControl.addTo(this.RoutesMap);
        this.RoutingWaypoints =[

        ];

        this.RoutesMap.on('click', function(e) {
          el.RoutingWaypoints.push(e.latlng);
          el.RoutingControl.setWaypoints(el.RoutingWaypoints);
        });

      }
    },
    mounted (){
      let el = this;
      $.get("/getKey/", function(data){
        el.APIKey = data;
        el.initMap();
      });
    }
});
