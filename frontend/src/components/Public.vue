<template>
<div style="padding: 0px; margin: 0px;width: 100%; height: 100%;">
    <div class="titleBar">
        <ul class="titleContent">
            <dropdown />
            <li class="title">RPI Shuttle Tracker</li>
        </ul>
        <div class="logo">
          <img src="~../assets/wtg.svg" />
        </div>
    </div>
    <span style="width: 100%; height: 100%; position: fixed;">
      <div id="mymap"></div>
    <messagebox />
    </span>
</div>
</template>

<script lang="ts">
import Vue from 'vue';
import InfoService from '../structures/serviceproviders/info.service';
import Vehicle from '../structures/vehicle';
import Route from '../structures/route';
import Stop from '../structures/stop';
import dropdown from './dropdown.vue';
import messagebox from './adminmessage.vue';
import * as L from 'leaflet';
import { setTimeout, setInterval } from 'timers';

const StopSVG = require('../assets/circle.svg') as string;

const StopIcon = L.icon({
  iconUrl: StopSVG,
  iconSize:     [12, 12], // size of the icon
  iconAnchor:   [6, 6], // point of the icon which will correspond to marker's location
  shadowAnchor: [6, 6],  // the same for the shadow
  popupAnchor:  [0, 0], // point from which the popup should open relative to the iconAnchor
});

export default Vue.extend({
  name: 'Public',
  data() {
    return ({
      vehicles: [],
      routes: [],
      stops: [],
      ready: false,
      Map: undefined,
      existingRouteLayers: [],
      initialized: false,
    } as {
      vehicles: Vehicle[],
      routes: Route[],
      stops: Stop[],
      ready: boolean,
      Map: L.Map | undefined, // Leaflets types are not always useful
      existingRouteLayers: L.Polyline[],
      initialized: boolean;
    });
  },
  mounted() {
    const a  = new InfoService();
    this.$store.dispatch('grabRoutes');
    this.$store.dispatch('grabStops');
    this.$store.dispatch('grabVehicles');
    this.$store.dispatch('grabUpdates');
    this.$store.dispatch('grabAdminMesssage');
    setInterval(() => {
      this.$store.dispatch('grabUpdates');
    }, 5000);

    this.$nextTick(() => {
      this.ready = true;
      this.Map = L.map('mymap', {
        zoomControl: false,
        attributionControl: false,
      });

      this.Map.setView([42.728172, -73.678803], 15.3);

      this.Map.addControl(L.control.attribution({
          position: 'bottomright',
          prefix: '',
      }));
      L.tileLayer('https://stamen-tiles.a.ssl.fastly.net/toner-lite/{z}/{x}/{y}{r}.png', {
        attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, ' +
                      'under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. ' +
                      'Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under ' +
                      '<a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
        // maxZoom: 17,
        // minZoom: 14,
      }).addTo(this.Map);

      this.Map.invalidateSize();

    });
    this.renderRoutes();
    this.$store.subscribe((mutation: any, state: any) => {
      if (mutation.type === 'setRoutes') {
        this.renderRoutes();
      }
      if (mutation.type === 'setStops') {
        this.renderStops();
      }
      if (mutation.type === 'setVehicles') {
        this.addVehicles();
      }
    });
  },
  methods: {
    routePolyLines(): L.Polyline[] {
      return this.$store.getters.getRoutePolyLines;
    },
    renderRoutes() {
      if (this.routePolyLines().length > 0 && !this.initialized) {
        if (this.Map !== undefined && !this.$store.getters.getBoundsPolyLine.isEmpty()) {
          this.initialized = true;
          this.Map.fitBounds(this.$store.getters.getBoundsPolyLine.getBounds());
        }
      }
      this.existingRouteLayers.forEach((line) => {
        if (this.Map !== undefined) {
          this.Map.removeLayer(line);
        }
      });
      this.existingRouteLayers = new Array<L.Polyline>();
      this.routePolyLines().forEach((line: L.Polyline) => {
        if (this.Map !== undefined) {
          this.Map.addLayer(line);
          this.existingRouteLayers.push(line);
        }
      });
    },
    renderStops() {
      this.$store.state.Stops.forEach((stop: Stop) => {
        const marker = L.marker([stop.lat, stop.lng], {icon: StopIcon});
        if (this.Map !== undefined) {
          marker.bindPopup(stop.name);
          marker.addTo(this.Map);
        }
      });
    },
    addVehicles() {
      this.$store.state.Vehicles.forEach((veh: Vehicle) => {
        if (this.Map !== undefined) {
          veh.addToMap(this.Map);
        }
      });
    },
  },
  components: {
    dropdown,
    messagebox,
  },
});
</script>

<style lang="scss">

#mymap{
    height: 100%;
    width: 100%;
    position: relative;
    filter: invert(0);
    
}

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
