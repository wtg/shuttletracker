<template>
  <div>
    <Header />
    <Map />
    <div id="right-panel">
      <ETAs v-bind:eta-info="currentETAInfo" v-bind:show="shouldShowETAMessage"></ETAs>
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
import Fusion from '../fusion';
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
    currentETAInfo: null,
  }as {
      fusion: Fusion;
     currentETAInfo: {} | null,
  };
},
mounted() {
  this.fusion.start();
},
  computed: {
    message(): AdminMessageUpdate {
        return this.$store.state.adminMessage;
    },
    busButtonActive(): boolean {
      return this.$store.getters.getBusButtonVisible;
    },
    shouldShowETAMessage(): boolean {
      return this.$store.state.settings.etasEnabled;
    },
    reconnecting(): boolean {
      return this.$store.state.fusionConnected === false;
    },
  },
  methods: {
        updateETA() {
      // find nearest stop
      const pos = UserLocationService.getInstance().getCurrentLocation();
      if (pos === undefined) {
        this.currentETAInfo = null;
        return;
      }
      const c = pos.coords as Coordinates;

      let minDistance = Infinity;
      let closestStop: Stop | null = null;
      for (const stop of this.$store.state.Stops) {
        const d = Math.hypot(c.longitude - stop.longitude, c.latitude - stop.latitude);
        if (d < minDistance) {
          minDistance = d;
          closestStop = stop;
        }
      }
      if (closestStop === null) {
        this.currentETAInfo = null;
        return;
      }

      // do we have an ETA for this stop? find the next soonest
      let eta: ETA | null = null;
      for (const e of this.$store.state.etas) {
        if (e.stopID === closestStop.id) {
          // is this the soonest?
          if (eta === null || e.eta < eta.eta || e.arriving) {
            eta = e;
          }
        }
      }
      if (eta === null) {
        this.currentETAInfo = null;
        return;
      }

      // get associated route
      let route: Route | null = null;
      for (const r of this.$store.state.Routes) {
        if (r.id === eta.routeID) {
          route = r;
          break;
        }
      }
      if (route === null) {
        this.currentETAInfo = null;
        return;
      }

      this.currentETAInfo = {eta, route, stop: closestStop};
    },
  },
components: {
    Map,
    Header,
    News,
    ETAs,
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
