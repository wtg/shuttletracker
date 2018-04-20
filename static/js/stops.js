var refresh = false;


Vue.component('stop-card', {
  props: ['info'],
  template:
  `<div class="tile is-child box">
    <b>ID</b>: {{info.id}} <br>
    <b>Name</b>: {{info.name}}<br>
    <b>Description</b>: {{info.description}}<br>
    <b>Enabled</b>: <input type="checkbox" v-model="info.enabled"></input>{{info.enabled}}<br>
    <b>Routes</b>:<br>
    <div v-for="route in myRouteIds">
      <label>{{route.name}}</label>
      <input type="checkbox" v-model="route.enabled" v-bind:value="route.enabled"></input>{{route.enabled}}
      <label>Order:</label>
      <input v-model="route.order" placeholder=route.order></input>
    </div>
    <br>
    <button class=" button" @click="shouldDelete(info.id)">{{deleteText}}</button>
    <button class=" button" @click="update(info.id)">{{updateText}}</button>

  </div>`,
  // v-bind:selected="info.routeId == routeId.id"
  // <b>currentRoute</b>: {{info.routeId}} -- previously on line 7
  //     <b>Order</b>: <input v-model="info.order"></input><br>
  // <input v-model="order" placeholder=info.order></input>
  //     <b>Route</b>:
  //        <select v-model="myRouteId" v-bind:id="'sel' + info.routeId">
  //        <option v-for="routeId in routeids" v-bind:value="routeId.id" v-bind:selected="info.routeId == routeId.id">{{routeId.name}}</option>
  //        </select><br>


  data (){
    return{
      //myData: {},
      deleteCount: 0,
      deleteText: "Delete",
      updateText: "Update",
      myRouteIds: [],
      //routeIDS: [],
      //myRouteId: "58d9b6af43162227cf364cce",
    };
  },
  mounted (){
    // add to myRouteIds from routeIds and from all routes
    var el = this;
    if (el.info.routeIds == null){
      console.log("null");
      $.get("/routes",function(data){
        console.log("data function");
        console.log(data);
        for (let i = 0; i < data.length; i ++){
          //el.routeIDS.push({id: data[i].id, name: data[i].name});
          el.myRouteIds.push({enabled: false, id: data[i].id, name: data[i].name, order: ""});
        }
      });
    } else {
      $.get("/routes",function(data){
        console.log("data function");
        console.log(data);
        for (let i = 0; i < data.length; i ++){
          var onRoute = false;
          for (let j = 0; j < el.info.routeIds.length; j++){
            if (data[i].id == el.info.routeIds[j]){
              el.myRouteIds.push({enabled: true, id: data[i].id, name: data[i].name, order: el.info.order[j]});
              onRoute = true;
              break;
            } 
          }
          if (!onRoute){
            el.myRouteIds.push({enabled: false, id: data[i].id, name: data[i].name, order:""});
          }
        }
      });    
    }
    console.log("card: " + el.myRouteIds.length);

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
        "routeIds":this.info.myRouteIds.toString(),
        "enabled":this.info.enabled.toString(),
        "toDelete":true,
        "lat":this.info.lat,
        "lng":this.info.lng,
        "order": this.info.order,
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
  `<div class ="column">
    <div v-for="stop in stopData" class="tile is-parent">
    <stop-card  v-bind:info="stop"></stop-card>
    </div>
    <stop-create></stop-create>
  </div>`,
  // v-bind:routeids="routeIDS" - previously on line 96
  data (){
    return{
      stopData: [{shuttleID:22}],
      routeIDS: [],
      myRouteIds: [],
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
        //el.myRouteIds.push({enabled: false, id: data[i].id, name: data[i].name, order: ""});

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

Vue.component('stop-create', {
  template:
  `<div><div class="route-description-box" style="height: 800px; padding-bottom: 10px;">
              <div id="stopsmap" style="height: 650px;float: left; width: 100%; background-color:black;z-index:0;"></div>
              <b>Name</b>: <input v-model="name" placeholder="Stop Name"></input><br>
              <b>Description</b>: <input v-model="description" placeholder="Stop Description"></input><br>
              <b>Enabled</b>:  <input type="checkbox" v-model="enabled"></input>{{enabled}}<br>
              <b>Route</b>:
              <div v-for="route in myRouteIds">
                <label>{{route.name}}</label>
                <input type="checkbox" v-model="route.enabled">{{route.enabled}}</input>
                <input v-model="route.order" placeholder= "Stop Order"></input>
                <br>
              </div>
            <button class="button" @click="send(JSON.stringify(getJSON()))">Submit</button>
            <button class="button" @click="showJSON(JSON.stringify(getJSON()))">GetJSON</button>
    </div></div>`,
    // v-bind:value="routeIDS[i]"
    //              <b>Order</b>: <input v-model="order" placeholder="Stop Order"></input><br>
    //              <b>Route</b>: <select v-model="myRouteId"><option v-for="routeId in routeIDS" v-bind:value="routeId.id">{{routeId.name}}</option></select><br>

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
        myRouteIds: [],
        routeIDS: [],
        order: [],
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
        var ids = [];
        var routeOrders = [];
        for (var i = 0; i < this.myRouteIds.length; i++){
          if (this.myRouteIds[i].enabled){
            ids.push(this.myRouteIds[i].id);
            routeOrders.push(this.myRouteIds[i].order);
          }
        }
        var toSend = {
          "name":this.name,
          "description":this.description,
          "startTime":"",
          "endTime":"",
          "enabled":this.enabled.toString(),
          "routeId": ids.toString(),
          "routeIds": ids,
          "lat": this.addStopMarker.getLatLng().lat.toString(),
          "lng": this.addStopMarker.getLatLng().lng.toString(),
          "order": routeOrders};
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
        L.tileLayer('https://stamen-tiles.a.ssl.fastly.net/toner-lite/{z}/{x}/{y}{r}.png', {
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
          el.routeIDS.push({id: data[i].id, name: data[i].name});
          el.myRouteIds.push({enabled: false, id: data[i].id, name: data[i].name, order: ""});

        }
        //console.log(el.myRouteIds.length);
        //var routes = routeIDS.length;
        refresh = false;
      });
    }
});