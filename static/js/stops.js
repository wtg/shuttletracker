var refresh = false;


Vue.component('stop-card', {
  props: ['info','routeids'],
  template:
  `<div class="vehicle-card route-description-box">
    <b>id</b>: {{info.id}} <b>currentRoute</b>: {{info.routeId}}<br>
    <b>name</b>: {{info.name}}<br>
    <b>Description</b>: {{info.description}}<br>
    <b>enabled</b>: <input type="checkbox" v-model="info.enabled"></input>{{info.enabled}}<br>
    <b>Route</b>: <select v-model="myRouteId" v-bind:id="'sel' + info.routeId"><option v-for="routeId in routeids" v-bind:value="routeId.id" v-bind:selected="info.routeId == routeId.id">{{routeId.name}}</option></select><br>
    <br>
    <button class=" button delete" @click="shouldDelete(info.id)">{{deleteText}}</button>
    <button class=" button delete" @click="update(info.id)">{{updateText}}</button>

  </div>`,
  data (){
    return{
      myData: {},
      deleteCount: 0,
      deleteText: "Delete",
      updateText: "Update",
      myRouteId: "58d9b6af43162227cf364cce",
    };
  },
  mounted (){
    this.myRouteId = this.info.routeId;
  },
  methods: {
    shouldDelete: function(id){
      if(this.deleteCount == 1){
        this.deleteStop(id);
      }else{
        this.deleteText = "Are you Sure?";
        this.deleteCount++;
      }
    },
    deleteStop: function(id){
      $.ajax({
           url: '/stops/' + id,
           type: 'DELETE',
           success: function(result) {
             refresh = true;
           }
         });
    },
    update: function(id){
      var el = this;
      var toSend = {
        "id": id.toString(),
        "name": this.info.name,
        "description":this.info.description,
        "startTime":"",
        "endTime":"",
        "routeId":this.myRouteId.toString(),
        "enabled":this.info.enabled.toString(),
        "toDelete":true,
        "lat":this.info.lat,
        "lng":this.info.lng,
      };
      $.ajax({
        url: "/stops/create",
        type: "POST",
        data: JSON.stringify(toSend),
        contentType: "application/json",
        complete: function(data){
          refresh=true;

        }
      });
    },

  }

});

Vue.component('stops-panel', {
  template:
  `<div class ="vehicle-panel">
    <stop-card v-for="stop in stopData" v-bind:info="stop" v-bind:routeids="routeIDS"></stop-card>
    <stop-create></stop-create>
  </div>`,
  data (){
    return{
      stopData: [{shuttleID:22}],
      routeIDS: [],
    };
  },
  mounted(){
    var el = this;
    $.get("/stops",function(data){
      el.stopData = data;
      refresh = false;
    });
    $.get("/routes",function(data){
      for (let i = 0; i < data.length; i ++){
        el.routeIDS.push({id: data[i].id, name: data[i].name});

      }
      refresh = false;
    });
    setInterval(function(){
      if(refresh){
        $.get("/stops",function(data){
          el.stopData = data;
          refresh = false;
        });
      }
    },100);

  }

});
Vue.component('stop-create',{
  template:`
  <div><div class="route-description-box" style="height: 800px; padding-bottom: 10px;">
              <div id="stopsmap" style="height: 650px;float: left; width: 100%; background-color:black;z-index:0;"></div>
              <b>name</b>: <input v-model="name" placeholder="Stop Name"></input><br>
              <b>Description</b>: <input v-model="description" placeholder="Stop Description"></input><br>
              <b>Enabled</b>:  <input type="checkbox" v-model="enabled"></input>{{enabled}}<br>
              <b>Route</b>: <select v-model="myRouteId"><option v-for="routeId in routeIDS" v-bind:value="routeId.id">{{routeId.name}}</option></select><br>
            <button class="button" @click="send(JSON.stringify(getJSON()))">Submit</button>
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
        StopsMap: null,
        APIKey: null,
        RoutingControl: null,
        drawnRoute: null,
        RoutingWaypoints: [],
        addStopMarker: null,
        myRouteId: "",
        routeIDS: [],

      };
    },
    methods: {
      send: function(data){
        $.ajax({
          url: "/stops/create",
          type: "POST",
          data: data,
          contentType: "application/json",
          complete: function(data){
            refresh=true;

          }
        });
      },
      getJSON: function(){
        var el = this;
        var toSend = {
          "name":this.name,
          "description":this.description,
          "startTime":"",
          "endTime":"",
          "enabled":this.enabled.toString(),
          "routeId":this.myRouteId.toString(),
          "lat": this.addStopMarker.getLatLng().lat.toString(),
          "lng": this.addStopMarker.getLatLng().lng.toString()};
          return toSend;
      },
      showJSON: function(data){
        var wnd = window.open("about:blank", "", "_blank");
        wnd.document.write(data);

      },

      initMap: function(){
        let el = this;
        this.StopsMap = L.map('stopsmap', {
          zoomControl: false,
          attributionControl: false
        });

        this.StopsMap.setView([42.728172, -73.678803], 15.3);
        this.StopsMap.addControl(L.control.attribution({
          position: 'bottomright',
          prefix: ''
        }));
        L.tileLayer('http://tile.stamen.com/toner-lite/{z}/{x}/{y}{r}.png', {
          attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
          minZoom: 13
        }).addTo(this.StopsMap);
        this.StopsMap.on( 'click', function(e){
          if(el.addStopMarker !== null){
            el.StopsMap.removeLayer(el.addStopMarker);
          }
          el.addStopMarker = L.marker(e.latlng,
            {
              draggable: true
            });
          el.addStopMarker.addTo(el.StopsMap);

        });
      }
    },
    mounted (){
      let el = this;
      this.initMap();
      $.get("/routes",function(data){
        for (let i = 0; i < data.length; i ++){
          el.routeIDS.push({id: data[i].id,name: data[i].name});

        }
        refresh = false;
      });
    }
});
