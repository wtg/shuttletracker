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
        zIndex: "20",
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
    <div v-bind:style="sidebarStyle" class="sidebar">
      <ul class ="nav-list">
        <li v-for="elem in elements" @click="setState(elem.id)" class="nav-item" v-bind:class="{ selected: (keepState == elem.id) }">{{elem.text}}</li>

      </ul>
    </div>
  `,
  data (){
    return{
      elements: [{text: "Routes",id: 0},{text: "Stops",id: 1},{text: "Vehicles",id: 2}],

      sidebarStyle: {
        backgroundColor:"white",
        height:"auto",
        bottom: "0",
        fontSize: "16px",
        zIndex: "10",
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
      this.keepState = id;
    },

  },
  mounted(){
  },
});

Vue.component("main-pane",{
  template:`
  <div v-bind:style="mainStyle">
    <transition name="slide-fade">
    <route-panel v-if="state == 0"></route-panel>
    <stops-panel v-if="state == 1"></stops-panel>
    <vehicle-panel v-if="state == 2"></vehicle-panel>
    </transition>
  </div>
  `,
  data (){
    return {
      state: 0,
      mainStyle: {position: "fixed",width: "auto",right:"0px",top:"50px",height: "auto", overflow: "scroll", bottom: "0",left: "150px"}
    }
  },
  mounted (){
    let el = this;
    setInterval(function(){
      if(el.state != state){
        el.state = state;
      }
    },10);
  },
});

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {

  }

});
$(document).ready(function(){

});
