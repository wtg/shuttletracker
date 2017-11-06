var main = Vue.component('mainui', {
  template:
  `<div class ="main">
  </div>`,
  data (){
    return{

    };
  },


});

var sidebar = Vue.component('sidebar', {
  props: ['elems'],
  template:
  `<div class ="sidebar">
  <ul>
   <li v-for="item in elems" :id=item.id @click=select(item.id)>{{item.text}}</li>
  </ul>
  </div> `,
  data (){
    return{

    };
  },
  methods: {
    select: function(data){
    }
  }


});

Vue.component('titlebar', {
  template:
  `<div>
  <div class ="title-bar">
  <p class="title"><span class = "red" >Shuttle</span>Tracker</p>
  </div>
  <sidebar :elems=elements ></sidebar>
  </div>`,
  data (){
    return{
      elements: [{text: "Routes",id: 0},{text: "Stops",id: 1},{text: "Vehicles",id: 2},{text: "Users",id: 3},{text: "Messages",id: 4}]
    };
  },
  mounted (){
  }

});

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {

  }

});
$(document).ready(function(){

});
