<template>
  <div>
    <p id="title-RPI">Shuttle Tracker</p>
    
    <TimesQueue/>
    <Map/>
    <TimeDisplay/> 

    <!-- Bottom Banner --> 
    <div class="titleBar bot">
      <img src="~../assets/icon.svg" id="shuttlepic">
      <h1 id="WTG">Web Technologies Group</h1>
    </div>

  </div>
</template>

<script lang="ts">
// This component handles the main TV Panel Application
import Vue from "vue";
import Map from "./tvpanel/Map.vue";
import TimesQueue from "./tvpanel/TimesQueue.vue";
import TimeDisplay from "./tvpanel/TimeDisplay.vue";
export default Vue.extend({
  name: "tvpanel",
  components: {
    TimesQueue,
    Map,
    TimeDisplay,
  },
});
</script>

<style>
html,body {
  height: 100%;
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
.bot {
  position: absolute;
  bottom: 0px;
  height: 60px;
}
</style>
