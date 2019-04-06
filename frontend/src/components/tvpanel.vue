<template>
  <div class="tvPanel">
    
    <!-- <TimesQueue/> -->
    <Map class="map"/>

    <!-- <eta></eta> -->
 
    <!-- Bottom Banner --> 
    <div class="sidebar has-text-centered">
      <!-- <img src="~../assets/icon.svg" id="shuttlepic"> -->
      <TimeDisplay /> 
      <eta></eta>
      <img class="logo" src="~../assets/icon.svg">
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
  height: calc( 100% + 44px );
  width: 100%;
  overflow: hidden;
}
#title-RPI {
  text-align: center;
  position: relative;
  top: 100px;
  font-size: 75px;
  color: black;
}
#WTG {
  position: absolute;
  right: 20px;
}

#wtglogo{
  right: 10px;
}

.map{
  width: calc( 100% - 400px );
  height: 100% !important;
}

.bot {
  position: absolute;
  bottom: 0px;
  height: 60px;
}

.tvPanel{
  width: 100%;
  height: calc( 100% + 44px ) !important;
  display: flex;
  flex-flow: row;
  justify-content: space-between;
}

.sidebar {
  width: 400px;
  border-bottom: 0.5px solid #eee;
  box-shadow: -3px 0px 8px 0 #ddd;
  height: calc( 100% + 44px );
  z-index: 3;
  padding-top:20px;
}

.titleBar {
  height: 40px;
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
  height:60px;
  width:60px;
  position:fixed;
  left: calc( 100% - 230px );
  bottom:30px;
}

</style>
