<template>
<div style="padding: 0px; margin: 0px;width: 100%; height: 100%;">
    <div class="titleBar">
        <ul class="titleContent">
            <dropdown />
            <li class="title">RPI Shuttle Tracker</li>
        </ul>
        <div v-if="$store.state.online" class="livebox">
          <p>Live </p>
          <div class="pulsate" style=""></div>
        </div>
        <div v-if="!$store.state.online" class="livebox">
          <p>Offline </p>
          <div class="caution-circle" style=""></div>
        </div>
        <div class="logo">
          <a href="https://webtech.union.rpi.edu/">
            <img src="~../assets/wtg.svg" />
          </a
>        </div>
    </div>
    <span style="width: 100%; height: 100%; position: fixed;">
      <div id="mymap"></div>
    <messagebox ref="msgbox" />
    </span>
</div>
</template>

<script lang="ts">
// This component handles everything on the shuttle tracker that is publicly facing.

import Vue from 'vue';
import InfoService from '../structures/serviceproviders/info.service';
import Vehicle from '../structures/vehicle';
import Route from '../structures/route';
import Stop from '../structures/stop';
import dropdown from './dropdown.vue';


import messagebox from './adminmessage.vue';
import * as L from 'leaflet';
import { setTimeout, setInterval } from 'timers';
import getMarkerString from '../structures/leaflet/rotatedMarker';
import { Position } from 'geojson';

const StopSVG = require('@/assets/circle.svg') as string;
const UserSVG = require('@/assets/user.svg') as string;

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
      legend: new L.Control({ position: 'bottomleft' }),
    } as {
      vehicles: Vehicle[],
      routes: Route[],
      stops: Stop[],
      ready: boolean,
      Map: L.Map | undefined, // Leaflets types are not always useful
      existingRouteLayers: L.Polyline[],
      initialized: boolean,
      legend: L.Control,
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
        maxZoom: 17,
        minZoom: 14,
      }).addTo(this.Map);

      this.Map.invalidateSize();
      this.showUserLocation();

    });
    this.renderRoutes();
    this.$store.subscribe((mutation: any, state: any) => {
      if (mutation.type === 'setRoutes') {
        this.renderRoutes();
        this.updateLegend();
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
    choose() {
      alert('hi');
    },
    updateLegend() {
      this.legend.onAdd = (map: L.Map) => {
        const div = L.DomUtil.create('div', 'info legend');
        let legendstring = '';
        this.$store.state.Routes.forEach((route: Route) => {
          if (route.enabled) {
            legendstring += `<li onclick="choose();"><img class="legend-icon" src=` + getMarkerString(route.color) + ` width="12" height="12"> ` +
            route.name + `</li>`;
            }
        });

        div.innerHTML = `<ul style="list-style:none">
					<li><img class="legend-icon" src='` + UserSVG + `' width="12" height="12"> You </li>` +
          legendstring +
          `<li><img class="legend-icon" src="` + StopSVG + `" width="12" height="12"> Shuttle Stop</li>
				</ul>`;
        return div;
      };
      if (this.Map !== undefined) {
        this.legend.addTo(this.Map);
      }
    },
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
        const marker = L.marker([stop.latitude, stop.longitude], {icon: StopIcon});
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
    showUserLocation() {
      const userIcon = L.icon({
        iconUrl: 'static/images/user.svg',

        iconSize:     [12, 12], // size of the icon
        iconAnchor:   [6, 6], // point of the icon which will correspond to marker's location
        shadowAnchor: [6, 6],  // the same for the shadow
        popupAnchor:  [0, 0], // point from which the popup should open relative to the iconAnchor
      });

      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(((position) => {
          const locationMarker = {
                name: 'You are here',
                marker: L.marker([position.coords.latitude, position.coords.longitude], {
                    icon: userIcon,
                    zIndexOffset: 1000,
                }),
          };
          locationMarker.marker.bindPopup(locationMarker.name);
          if (this.Map !== undefined) {
            locationMarker.marker.addTo(this.Map);
          }
        }) as PositionCallback);
      }
    },
  },
  components: {
    dropdown,
    messagebox,
  },
});
</script>

<style lang="scss">
.caution-circle {
  float:right;
  width:10px;
  height:10px;
  background-color:orange;
  border-radius:50%;
}
.pulsate {
  float:right;
  width:10px;
  height:10px;
  background-color:blue;
  border-radius:50%;
  animation: pulsate 2.5s ease-out;
  animation-iteration-count: infinite;
}
@keyframes pulsate {
  0% {
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0;

  }
}
.livebox{
    p {
      margin-right: 5px;
    }
    position: absolute;
    height: 26px;
    right: 10px;
    top: 40px;
    box-shadow: rgba(0, 0, 0, 0.8) 0px 1px 1px;
    border-radius: 5px;
    padding-left: 4px;   
    padding-right: 4px;   
    justify-self: flex-end;
    display: flex;
    flex-flow: row wrap;
    align-content: center;
    align-items: center;
    justify-content: space-around;
    background-color: rgba(255, 255, 255, 0.9);
    
}

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
    background-color: rgba(255, 255, 255, 0.9);
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
        line-height: 34px;
        vertical-align: center;
        padding: 0px;
        margin: 0px 6px 0px;
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

.info.legend{
  box-shadow: rgba(0, 0, 0, 0.8) 0px 1px 1px;
  border-radius: 5px;
  background-color: rgba(255, 255, 255, .9);
  padding: 5px;
  bottom: 25px;
  & ul{
    margin-top: 2px;
    margin-bottom: 2px;
    padding-left: 0px;
  }

  li.legendButton{
      padding: 1px;
      background-color: #e7e7e7;
      color: black;
  }
}
</style>
