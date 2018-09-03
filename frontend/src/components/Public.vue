<template>
  <div style="padding: 0px; margin: 0px;">
    <div class="titleBar">
      <ul class="titleContent"><li class="title">RPI Shuttle Tracker</li></ul>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import InfoService from '../structures/serviceproviders/info.service';
import Vehicle from '../structures/vehicle';
import Route from '../structures/route';
import Stop from '../structures/stop';

export default Vue.extend({
  name: 'Public',
  data() {
    return ({
      vehicles: [],
      routes: [],
      stops: [],
    } as {
      vehicles: Vehicle[],
      routes: Route[],
      stops: Stop[],
    });
  },
  mounted() {
    const a  = new InfoService();
    a.GrabVehicles().then((data: Vehicle[]) => this.vehicles = data);
    a.GrabRoutes().then((data: Route[]) => this.routes = data);
    a.GrabStops().then((data: Stop[]) => this.stops = data);
  },
});

</script>

<style lang="scss">
.titleBar {
    height: 34px;
    float: none;
    position: absolute;
    z-index: 1;
    width: 100%;
    background-color: rgba(255, 255, 255, 0.88);
    box-shadow: 0 -5px 10px rgba(0, 0, 0, 0.8);
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;

    & .titleContent {
      height: 100%;
      z-index: 1;
      width: auto;
      list-style: none;
      position: relative;
      top: 0px;
      margin: 0px;
      padding: 0px;

      & .title {
        font-size: 22px;
        padding: 0px;
        margin: 3px 6px 0px;
        float: left;
      }
    }
}
</style>
