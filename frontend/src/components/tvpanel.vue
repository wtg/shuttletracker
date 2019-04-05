<template>
  <div>
    <p id="title-RPI">Shuttle Tracker 🚐💨</p>
    
    <!-- <TimesQueue/> -->
    <Map/>
    <!-- <TimeDisplay/>  -->
    <eta></eta>

    <!-- Bottom Banner --> 
    <div class="titleBar bot">
      <img src="~../assets/icon.svg" id="shuttlepic">
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

.bot {
  position: absolute;
  bottom: 0px;
  height: 60px;
}

.titleBar {
  height: 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  border-bottom: 0.5px solid #eee;
  box-shadow: 0 3px 8px 0 #ddd;
  user-select: none;
  background: white;
  z-index: 1;

}
</style>
