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
          <tr v-for="(eta, i) in etas" v-bind:key="i">
            <td>{{ eta.vehicle.name }}</td>
            <td>{{ eta.stop.name }}</td>
            <td>{{ eta.eta }}</td>
            <td>{{ eta.arriving }}</td>
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
    etas(): any[] {
      const ret = [];
      for (const eta of this.$store.state.etas) {
        const e = Object.create({eta: eta.eta, arriving: eta.arriving});

        for (const vehicle of this.$store.state.Vehicles) {
          if (eta.vehicleID === vehicle.id) {
            e.vehicle = vehicle;
            break;
          }
        }

        for (const stop of this.$store.state.Stops) {
          if (eta.stopID === stop.id) {
            e.stop = stop;
            break;
          }
        }

        ret.push(e);
      }
      return ret.sort((a, b) => {
        if (a.vehicle.name > b.vehicle.name) {
          return 1;
        } else if (a.vehicle.name < b.vehicle.name) {
          return -1;
        }
        return 0;
      });
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
