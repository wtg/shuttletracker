var state = 0;

Vue.component('titlebar', {
  template:
  `<div>
  <div v-bind:style="titlebarStyle" class ="title-bar">
  <p class="title"><span class = "red" >Shuttle</span>Tracker</p>
  <a v-bind:style="logoutStyle" href="/admin/logout/">Logout</a>
  </div>

  </div>`,

  data (){
    return{
      titlebarStyle: {
        backgroundColor:"white",
        width: "100%",
        height:"50px",
        float:"left",
        position: "absolute",
        lineHeight:"50px",
        verticalAlign:"center",
        boxShadow: '0 3px 7px rgba(0,0,0,0.25)'},
      logoutStyle: {
        float: "right",
        textDecoration: "none",
        paddingRight: "10px"
      }
    };
  },
  mounted (){
  }

});

Vue.component('sidebar',{
  template: `
    <div v-bind:style="titlebarStyle" class="sidebar">
      <ul class ="nav-list">
        <li v-for="elem in elements" @click="setState(elem.id)" class="nav-item" v-bind:class="{ selected: (keepState == elem.id) }">{{elem.text}}</li>

      </ul>
    </div>
  `,
  data (){
    return{
      elements: [{text: "Routes",id: 0},{text: "Stops",id: 1},{text: "Vehicles",id: 2},{text: "Users",id: 3},{text: "Messages",id: 4},{text: "Logout",id:5}],

      titlebarStyle: {
        backgroundColor:"white",
        width: "10%",
        height:"auto",
        bottom: "0",
        top:"40",
        float:"left",
        position: "absolute",
        lineHeight:"50px",
        boxShadow: '0 3px 7px rgba(0,0,0,0.25)'},
      keepState: 0,

    };
  },
  methods:{
    setState(id){
      state = id;
      keepState = id;
    },

  },
  mounted(){
  },
});

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {

  }

});
$(document).ready(function(){

});
