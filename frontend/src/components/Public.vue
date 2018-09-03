<template>
<div style="padding: 0px; margin: 0px;">
    <div class="titleBar">
        <ul class="titleContent">
            <dropdown />
            <li class="title">RPI Shuttle Tracker</li>
        </ul>
        <div class="logo">
          <img src="~../assets/wtg.svg" />
        </div>
    </div>
    {{vehicles}}
    {{routes}}
    {{stops}}
</div>
</template>

<script lang="ts">
import Vue from 'vue';
import InfoService from '../structures/serviceproviders/info.service';
import Vehicle from '../structures/vehicle';
import Route from '../structures/route';
import Stop from '../structures/stop';
import dropdown from './dropdown.vue';

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
  components: {
    dropdown,
  },
});
</script>

<style lang="scss">
.titleBar {
    height: 34px;
    float: none;
    position: absolute;
    z-index: 1;
    display: flex;
    align-content: space-around;
    justify-content: space-between;
    flex-flow: row;
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

.logo {
  height: 24px;
  float: right;
  padding-right: 10px;
  align-self: center;
  & img{
    height: 100%;
  }
}
</style>
