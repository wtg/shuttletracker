<template>
  <div>
    <div id="eta-title">
      ETA
    </div>
    <div class="eta-message">
      East Campus shuttle arriving in 1 minute.
    </div>
    <div class="eta-message">
      West Campus shuttle arriving in 3 minutes.
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
export default Vue.extend({
  data() {
      return {
          etaInfo: null,
      }
  },
  computed: {
    message(): string | null {
      return null;
      // find nearest stop
      let unionStop: Stop | null;
      for (const stop of this.$store.state.Stops) {
        if (stop.name === "Student Union") {
          unionStop = stop;
          break;
        }
      }
      if (unionStop === null) {
        return "No ETA";
      }

      // do we have an ETA for this stop? find the next soonest
      let eta: ETA | null = null;
      for (const e of this.$store.state.etas) {
        if (e.stopID === unionStop.id) {
          // is this the soonest?
          if (eta === null || e.eta < eta.eta || e.arriving) {
            eta = e;
          }
        }
      }
      if (eta === null) {
        this.currentETAInfo = null;
        return "No ETA";
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

      // this.currentETAInfo = {eta, route, stop: closestStop};
      const etas = this.$store.state.etas;
      
      if (this.etaInfo === null) {
          return null;
      }
      const now = new Date();
      let newMessage = `${this.etaInfo.route.name} shuttle arriving at ${this.etaInfo.stop.name}`;
      // more than 1 min 30 sec?
      if (this.etaInfo.eta.eta.getTime() - now.getTime() > 1.5 * 60 * 1000 && !this.etaInfo.eta.arriving) {
        newMessage += ` in ${relativeTime(now, this.etaInfo.eta.eta)}`;
      }
      newMessage += '.';
      return newMessage;
    },
  },
});
function relativeTime(from: Date, to: Date): string {
  const minuteMs = 60 * 1000;
  const elapsed = to.getTime() - from.getTime();
  // cap display at thirty min
  if (elapsed < minuteMs * 30) {
    return `${Math.round(elapsed / minuteMs)} minutes`;
  }
  return 'a while';
}
</script>

<style lang="scss" scoped>
.eta-message {
    // width: 320px;
    // position: fixed;
    // top: 0;
    // left: 0;
    // right: 0;
    margin: 20px 20px;
    // z-index: 1000;
    background: white;
    padding: 20px 28px;
    border: 0.5px solid #eee;
    border-radius: 4px;
    // box-shadow: 0 1px 16px -4px #bbb;
    font-size: 18px;
}
#eta-title{
  font-size:60px;
  font-weight:400;
}
@media screen and (max-width: 500px) {
    .eta-message {
        width: auto;
        margin: 50px 10px 0 10px;
        padding: 16px 22px;
        font-size: 16px;
    }
}
</style>