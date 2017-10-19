$(document).ready(function(){
  console.log("asdf");
});

refresh = true;

Vue.component('vehicle-card', {
  props: ['data'],
  template:
  `<div class ="vehicle-card route-description-box">
    {{data.vehicleID}}
  </div>`,
  data (){
    return{
      shuttleData: [{shuttleID:22}]
    }
  },
  mounted(){
    var el = this;

  }

})

Vue.component('vehicle-panel', {
  template:
  `<div class ="vehicle-panel">
    <div v-for="vehicle in shuttleData" class="vehicle-info">
      <vehicle-card v-bind:data="vehicle"></vehicle-card>
    </div>
  </div>`,
  data (){
    return{
      shuttleData: [{shuttleID:22}]
    }
  },
  mounted(){
    var el = this;
    setInterval(function(){
      if(1 == 1){
        $.get("/vehicles",function(data){
          console.log("wowhee");
          el.shuttleData = data;
          refresh = false;
        })

      }
    },1000)

  }

})

var ShuttleTracker = new Vue({
  el: '#document-vue',
  data: {
  }
})
