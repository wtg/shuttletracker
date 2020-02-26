<template>
  <div class="parent">
    <h1 class="title">ETAs</h1>
    <hr>

    <div class="container">
      <table class="table">
        <thead>
          <tr>
            <th>Vehicle</th>
            <th>Stop</th>
            <th>ETA</th>
            <th>Route ID</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!etas.length">No ETAS Currently Calculated</tr>
            <template v-for="(eta) in etas">
              <tr v-for="(info, j) in eta" v-bind:key="j">
                <td>{{ info.vehicleID }}</td>
                <td>{{ info.stopID }}</td>
                <td>{{ info.eta }}</td>
                <td>{{ info.routeID }}</td>
              </tr>
            </template>
        </tbody>
      </table>
    </div>

  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import ETA from '@/structures/eta';

export default Vue.extend({
  computed: {
    etas(): any[] {
      const etaArray = [];
      for (let i = 0; i < 11; i++) {
        const etaString = localStorage.getItem(String(i + 1));
        if (etaString) {
          const localETA = JSON.parse(etaString);
          const ret = [];
          if (localETA.length) {
            for (const eta of localETA) {
              const now = new Date();
              const from = new Date(eta.eta);
              const minuteMs = 60 * 1000;
              const elapsed = from.getTime() - now.getTime();

              let etaMinutes = `A while`;
              // cap display at 15 min
              if (elapsed < minuteMs * 15) {
                etaMinutes =  `${Math.round(elapsed / minuteMs)} minutes`;
              }

              const e = {eta: etaMinutes,
                        arriving: eta.arriving,
                        vehicleID: eta.vehicleID,
                        stopID: eta.stopID,
                        routeID: eta.routeID};
              ret.push(e);
            }
            etaArray.push(ret);
          }
        }

      }
      console.log(etaArray);
      return etaArray;
      // return ret.sort((a, b) => {
      //   if (a.vehicle.name > b.vehicle.name) {
      //     return 1;
      //   } else if (a.vehicle.name < b.vehicle.name) {
      //     return -1;
      //   }
      //   return 0;
      // });
    },
  },
  watch: {
    etas(): function() {
      console.log(this.$store.state);
    }
  }
});

window.addEventListener('storage', () => {
  location.reload();
});
</script>

<style lang="scss" scoped>
.parent {
  padding: 20px;
}
.caption {
  margin-bottom: 1em;
}
</style>
