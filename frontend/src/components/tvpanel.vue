<template>
  <div>
    <Header />
    <Map />
    <div id="right-panel">
      <TimeDisplay />
      <ETAs />
      <News />
    </div>
  </div>
</template>

<script lang='ts'>
import Vue from 'vue';
import Map from './tvpanel/Map.vue';
import Header from './tvpanel/Header.vue';
import News from './tvpanel/News.vue';
import ETAs from './tvpanel/ETAs.vue';
import TimeDisplay from './tvpanel/TimeDisplay.vue';
import Fusion from '@/fusion';
import messagebox from './adminmessage.vue';
import AdminMessageUpdate from '../structures/adminMessageUpdate';
import UserLocationService from '../structures/userlocation.service';
import { Stop, StopSVG } from '../structures/stop';
import ETA from '../structures/eta';
import Route from '../structures/route';
export default Vue.extend({
  name: 'tvpanel',
  data() {
    return{
      fusion: new Fusion(),
    } as {
      fusion: Fusion;
    };
  },
  mounted() {
    this.fusion.start();
  },
  components: {
    Map,
    Header,
    News,
    ETAs,
    TimeDisplay,
  },
});
</script>

<style lang='scss'>

#view-wrapper {
  height: 100%;
}

#right-panel {
  position: absolute;
  right: 0px;
  height: 100%;
  width: 25%;
}

</style>
