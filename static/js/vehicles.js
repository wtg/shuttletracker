$(document).ready(function(){
  console.log("asdf");
});

refresh = true;

Vue.component('vehicle-create',{
  template:`<div class="vehicle-card route-description-box">

    <b>id</b>: <input type="textbox" v-model="ID" placeholder="1123454125"></input> (must be same as itrak vehicle ID) <br>
    <b>name</b>:<input type="textbox" v-model="name" placeholder="Vehicle Name"></input> <br>

    <b>active</b>:<input type="checkbox" v-model="active"></input><br>
    <div class = "button" @click="send" style="width: 50px;">add</div>

    </div>`,
    data (){
      return{

        ID: "",
        name: "",
        active: true
      };
    },
    methods: {
      send: function(){

        var pkg = {"vehicleID":this.ID, "vehicleName":this.name, "active":this.active};

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
      }
    }
});

Vue.component('vehicle-card', {
  props: ['info'],
  template:
  `<div class="vehicle-card route-description-box">
    <b>id</b>: {{info.vehicleID}}<br>
    <b>name</b>: <input type="textbox" v-model="info.vehicleName"></input> <br>
    <b>active</b>: <input type="checkbox" v-model="info.active"></input>{{info.active}}<br>
    <b>Created</b>: {{info.Created}} <br>
    <div @click="editVehicle" class = "button" style="width: auto; float:left;">change</div>
    <div @click="deleteVehicle" class = "button" style="width: auto; float:left;">delete</div>
    <br>
  </div>`,
  data (){
    return{
      myData: {}
    };
  },
  methods: {
    deleteVehicle: function(){
      var el = this;
      $.ajax({
        url: '/vehicles/' + el.info.vehicleID,
        type: 'DELETE',
        success: function(result) {
          refresh = true;

        }
      });
    },
    editVehicle: function(){
      var el = this;

      var pkg = {"vehicleID":this.info.vehicleID, "vehicleName":this.info.vehicleName, "active":this.info.active};
      $.ajax({
        url: "/vehicles/edit",
        type: "POST",
        data: JSON.stringify(pkg),
        contentType: "application/json",
        complete: function(data){
          refresh = true;

        }
      });
    }
  }

});

Vue.component('vehicle-panel', {
  template:
  `<div class ="vehicle-panel">
    <vehicle-create></vehicle-create>
    <div v-for="vehicle in shuttleData" class="vehicle-info">
      <vehicle-card v-bind:info="vehicle"></vehicle-card>
    </div>
  </div>`,
  data (){
    return{
      shuttleData: [{shuttleID:22}]
    };
  },
  mounted(){
    var el = this;
    $.get("/vehicles",function(data){
      console.log(data);
      el.shuttleData = data;
      refresh = false;
    });
    setInterval(function(){
      if(refresh){
        $.get("/vehicles",function(data){
          el.shuttleData = data;
          refresh = false;
        });
      }
    },100);

  }

});

var ShuttleTracker = new Vue({
  el: '#document-vue',
  data: {
  }

});
