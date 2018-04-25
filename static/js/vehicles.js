
Vue.component('vehicle-create',{
  template:`<div class="vehicle-card route-description-box">
    <b>Tracker ID</b>: <input type="textbox" v-model="trackerID" placeholder="1123454125"></input> (must be same as iTrak vehicle ID) <br>
    <b>name</b>:<input type="textbox" v-model="name" placeholder="Vehicle Name"></input> <br>
    <b>enabled</b>:<input type="checkbox" v-model="enabled"></input><br>
    <div class = "button" @click="send" style="width: 50px;">add</div>
    </div>`,
    data (){
      return{
        trackerID: "",
        name: "",
        enabled: true,
      };
    },
    methods: {
      send: function(){

        var pkg = {
          "tracker_id": this.trackerID,
          "name":this.name,
          "enabled":this.enabled
        };

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
        this.trackerID = "";
        this.name = "";
      }
    }
});

Vue.component('vehicle-card', {
  props: ['info'],
  template:
  `<div class="vehicle-card route-description-box">
    <b>Tracker ID</b>: {{info.tracker_id}}<br>
    <b>name</b>: <input type="textbox" v-model="info.name"></input> <br>
    <b>enabled</b>: <input type="checkbox" v-model="info.enabled"></input>{{info.enabled}}<br>
    <b>Created</b>: {{info.created}} <br>
    <div @click="editVehicle" class = "button" style="width: auto; float:left;">Change</div>
    <div @click="deleteVehicle" class = "button" style="width: auto; float:left;">Delete</div>
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
        url: '/vehicles?id=' + el.info.id,
        type: 'DELETE',
        success: function(result) {
          refresh = true;
        }
      });
    },
    editVehicle: function(){
      var el = this;
      var pkg = {
        "id": this.info.id,
        "name": this.info.name,
        "enabled": this.info.enabled,
        "tracker_id": this.info.tracker_id
      };
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
      shuttleData: []
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
