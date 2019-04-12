<template>
  <div class="tvPanel">
    
    <!-- Map Component --> 
    <Map class="map"/>

    <!-- Side Bar --> 
    <div class="sidebar has-text-centered">  
      <TimeDisplay />
      <div id="queue">
        <TimesQueue /> 
        <img class="logo" src="~../assets/icon.svg">
        <div id="eta"> 
        <eta > </eta>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
// This component handles the main TV Panel Application
import Vue from "vue";
import Map from "./tvpanel/Map.vue";
import TimesQueue from "./tvpanel/TimesQueue.vue";
import TimeDisplay from "./tvpanel/TimeDisplay.vue";
import ETA from "./tvpanel/eta.vue";
export default Vue.extend({
  name: "tvpanel",
  components: {
    TimesQueue,
    Map,
    TimeDisplay,
    eta: ETA,
  },
});
</script>

<style lang="scss" scoped>
html,body {
  height: 100%;
  width: 100%;
  overflow:hidden;
}
#title-RPI {
  text-align: center;
  position: relative;
  top: 100px;
  font-size: 75px;
  color: black;
}

.map{
  width: calc( 100% - 400px );
}
.tvPanel{
  width: 100%;
  height: 100%;
  display: flex;
  flex-flow: row;
  justify-content: space-between;
}
.sidebar {
  width: 400px;
  border-bottom: 0.5px solid #eee;
  box-shadow: -3px 0px 8px 0 #ddd;
  height:100%;
  z-index: 3;
  padding-top:20px;
}
.titleBar {
  height: 100%;
  justify-content: center;
  align-items: center;
  width: 100%;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  border-bottom: 0.5px solid #eee;
  box-shadow: -3px 0 8px 0 #ddd;
  user-select: none;
  background: white;
  z-index: 1;
}
.logo {
  height:20px;
  width:20px;
  position:relative;
  top:600px;
  // left: calc( 100% - 230px );
}
#queue {
  height:100%;
  width:100%;
  position:relative;
  top:0px;
}
#eta {
  height:200px;
  width:100%;
  position:relative;
  top:55%;
}
</style>