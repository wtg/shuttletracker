<template>
  <div id="app">
    <p id="title-RPI">RPI Shuttle Tracker</p> 
    <div id="queue">
      <TimesQueue />
      <Map />
    </div>
  </div>
</template>

<script lang="ts">
// This component handles the main TV Panel Application
import Vue from 'vue';
import Map from './tvpanel/Map.vue';
import TimesQueue from './tvpanel/TimesQueue.vue';
export default Vue.extend({
  name: 'tvpanel',
  components: {
    TimesQueue,
    Map,
  },
});
</script>

 
<style>
html, body {
    height:100%;
    width:100%;
    overflow:hidden;
}
#title-RPI {
  text-align:center;
  position:relative;
  bottom:40px;
  font-size:75px;
}


</style>
