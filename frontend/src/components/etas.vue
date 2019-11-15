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
            <th>Arriving</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!etas.length">No ETAS Currently Calculated</tr>
          <tr v-for="(eta, i) in etas" v-bind:key="i">
            <td>{{ String(eta[i]) }}</td>
            <td>{{ String(eta[i]) }}</td>
            <td>{{ String(eta[i]) }}</td>
            <td>{{ String(eta[i]) }}</td>
          </tr>
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
    etas(): ETA[] {
      const ret = [];
      const etaString = localStorage.getItem('etas');
      if (etaString) {
        const localETA = JSON.parse(etaString);
        console.log(localETA);
        for (const eta of localETA) {
          const e = Object.create({eta: localETA.eta, arriving: localETA.arriving,
          vehicleID: localETA.vehicleID, stopID: localETA.stopID, routeID: localETA.routeID});
          ret.push(e);
        }
      }
      console.log(ret);
      return ret;
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
