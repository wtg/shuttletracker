<template>
  <div>
    <!-- Map Component -->
    <Map class="map"/>
  </div>
</template>

<script lang='ts'>
import Vue from 'vue';
import Map from './tvpanel/Map.vue';
import Fusion from '@/fusion';
export default Vue.extend({
  name: 'tvpanel',
  data() {
  return{
    fusion: new Fusion(),
  }as {
      fusion: Fusion;
  };
},
mounted() {
  this.fusion.start();
},
  components: {
    Map,
  },
});
</script>

<style lang='scss'>

#view-wrapper {
  height: 100%;
}

.map {
  position: relative;
  height: 100%;
  width: 70%;
}
</style>
