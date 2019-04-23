<template>
  <div class="tvPanel">
    
    <!-- Map Component --> 
    <Map class="map"/>

    <!-- Side Bar --> 
    <div class="sidebar has-text-centered">  
      <TimeDisplay /> 
          <div id="eta"> 
            <eta > </eta>
          </div>

          <div id="queue">
            <TimesQueue />
          </div>
          <div id="bottom">
              <img src="~../assets/icon.svg" id="shuttle-logo">
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
import Fusion from '@/fusion';

import ETA from "./tvpanel/eta.vue";

export default Vue.extend({
  name: "tvpanel",
  data() {
    return{
      fusion: new Fusion()

    }as {
        fusion: Fusion;
    }
  },
  mounted(){
    this.fusion.start();

  },
  components: {
    TimesQueue,
    Map,
    TimeDisplay,
    eta: ETA,
  },
});
</script>

<style lang="scss" >
html{
  overflow-y:hidden;
}
html,body {
  height: 100%;
  width: 100%;
}
#title-RPI {
  text-align: center;
  position: relative;
  top: 100px;
  font-size: 75px;
  color: black;
}
#shuttle-logo{
  position:absolute;
  top:20px;
  left:180px;
  height:40px;
  width:40px;
}
.map{
  width: calc( 100% - 400px );
}
.tvPanel{
  display: flex;
  flex-flow: row;
  justify-content: space-between;
}
.sidebar {
  width: 400px;
  border-bottom: 0.5px solid #eee;
  box-shadow: -3px 0px 8px 0 #ddd;
  z-index: 3;
}
#queue {
  height:500px;
  width:100%;
  position:relative;
  top:180px;
}
#eta {
  height:200px;
  width:100%;
  position:relative;
  top:0;
}
#bottom{
  position:fixed;
  bottom:0;
  width:100%;
  height:70px;
  border-top:2px solid#D3D3D3; 
}
</style>
