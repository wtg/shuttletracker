<template>
  <div class="parent">
    <h1 v-if="!isIntegrated" class="title">ETAs</h1>
    <hr v-if="!isIntegrated">

    <!--<div :class="isMobile ? '' : 'container'">
      <table class="table">
        <thead>
        <tr>
          <th>Route ID</th>
          <th v-show="!isMobile">Vehicle</th>
          <th v-show="!isMobile">Stop</th>
          <th>ETA</th>
          <th>Arriving</th>
        </tr>
        </thead>
        <tbody>
        <tr v-if="!etas.length">
          <td colspan="3"> No ETAS Currently Calculated</td>
        </tr>

        <template v-if="!isMobile">&lt;!&ndash; :key="`${i}-${eta.stopID}`"&ndash;&gt;
          <tr v-for="(eta, i) in etas" @click="selection=i" class="table-row"
              :class="[i % 2 === 1 ? 'stressed' : '']">
            <td>
              <div :style="{ color: eta.routeColor, display: 'inline'}"> &#9830;</div>
              {{ eta.routeName }}
            </td>
            <td>{{ eta.vehicleID }}</td>
            <td>{{ eta.stopName }}</td>
            <td>{{ eta.eta }}</td>
            <td>{{ eta.arriving }}</td>
          </tr>
        </template>
        <template v-else>&lt;!&ndash; :key="`${i}-${eta.stopID}`"&ndash;&gt;
          <div v-for="(eta, i) in etas" @click="selection=i">
            <div class="cell" style="width: 40%">
              <span :style="{ color: eta.routeColor}"> &#9830;</span>
              {{ eta.routeName }}
            </div>
            <div class="cell">{{ eta.eta }}</div>
            <div class="cell">{{ eta.arriving }}</div>
            <div></div>
          </div>
        </template>
        </tbody>
      </table>
    </div>-->

    <div class="container">
      <div>
        <div class="table-header" v-for="(info, j) in tableInfo" :key="j" :style="{width: info.width}">
          {{ info.header }}
        </div>
      </div>
      <div v-if="!etas.length">No ETAS Currently Calculated</div>
      <template v-else>
        <div v-for="(eta, i) in etas" @click="selection = selection === i ? -1 : i" :class="['table-row', i%2===1?'stressed':'']">
          <div v-for="(info, j) in tableInfo" :key="j" class="cell" :style="{width: info.width}">
            <span v-if="j===0" :style="{ color: eta.routeColor}"> &#9830;</span>
            {{ eta[info.content] }}
          </div>
          <div v-if="isMobile" style="width: 100%" class="modifiable-cell" :class="selection === i ? 'expanded':'collapsed'">
            vehicle-{{ eta.vehicleID }} at {{ eta.stopName }}
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import ETA from '@/structures/eta';
import axios from 'axios';

export default Vue.extend({
  props: {
    isIntegrated: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      routes: new Map(),
      colors: new Map(),
      stops: new Map(),
      selection: -1,
      isMobile: window.innerWidth < 768 || this.isIntegrated,
    };
  },
  mounted() {
    axios
        .get('http://localhost:8080/routes')
        .then((response) => {
          const dictionary = new Map();
          const dictionary2 = new Map();
          response.data.forEach((x: any) => {
            dictionary.set(x.id, x.name);
            dictionary2.set(x.id, x.color);
          });

          this.routes = dictionary;
          this.colors = dictionary2;
          // console.log(this.routes);
          // console.log(this.colors);
        })
        .catch((error) => console.log(error));
    axios
        .get('http://localhost:8080/stops')
        .then((response) => {
          const dictionary = new Map();
          response.data.forEach((x: any) => {
            dictionary.set(x.id, x.name);
          });
          this.stops = dictionary;
          // console.log(this.stops);
        })
        .catch((error) => console.log(error));
    // localStorage.clear();
    for (let i = 0; i < 10; i++) {
      const testString = [{stopID: 10, vehicleID: 3, routeID: 22, eta: Date.now() + 35000 * i, arriving: true}];
      localStorage.setItem(String(i), JSON.stringify(testString));
    }
    window.addEventListener('resize', this.resize);
  },
  computed: {
    etas(): any[] {

      const etaArray = [];
      for (let i = 0; i < 18; i++) {
        const etaString = localStorage.getItem(String(i + 1));
        // const testString = '[{"stopID": 10, "vehicleID": 3, "routeID": 22, "eta": "2020-10-28T01:09:09.826Z", "arriving": true}]';
        // console.log('abcd');
        if (etaString) {
          const localETA = JSON.parse(etaString);
          // const ret = [];
          if (localETA.length) {
            for (const eta of localETA) {

              /*Example of eta:
              {
                stopID: 10
                vehicleID: 3,
                routeID: 24,
                eta: "2020-10-27T21:09:09.826Z",
                arriving: true}

                (All times are in UTC)

              */

              const now = new Date();
              const from = new Date(eta.eta);
              const minuteMs = 60 * 1000;
              const elapsed = from.getTime() - now.getTime();
              // console.log(elapsed);
              let etaMinutes = `A while`;
              // cap display at 15 min
              if (elapsed < minuteMs * 15) {
                if (Math.round(elapsed / minuteMs) === 0) {
                  etaMinutes = `Less than a minute`;
                } else if (Math.round(elapsed / minuteMs) === 1) {
                  etaMinutes = `${Math.round(elapsed / minuteMs)} minute`;
                } else {
                  etaMinutes = `${Math.round(elapsed / minuteMs)} minutes`;
                }
              }

              const e = {
                eta: etaMinutes,
                arriving: eta.arriving,
                vehicleID: eta.vehicleID,
                stopID: eta.stopID,
                stopName: this.stops.get(eta.stopID),
                routeID: eta.routeID,
                routeName: this.routes.get(eta.routeID),
                routeColor: this.colors.get(eta.routeID),
              };
              if (elapsed >= 0) {
                // ret.push(e);
                etaArray.push(e);
              }
            }
            // etaArray.push(ret);
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
    tableInfo(): any[] {
      return this.isMobile ? [
        {header: 'Route ID', content: 'routeName', width: '50%'},
        // {header: 'Vehicle', content: 'vehicleID', width: '20%'},
        // {header: 'Stop', content: 'stopName', width: '20%'},
        {header: 'ETA', content: 'eta', width: '25%'},
        {header: 'Arriving', content: 'arriving', width: '25%'},
      ] : [
        {header: 'Route ID', content: 'routeName', width: '30%'},
        { header: 'Vehicle', content: 'vehicleID', width: '10%' },
        { header: 'Stop', content: 'stopName', width: '30%' },
        {header: 'ETA', content: 'eta', width: '20%'},
        {header: 'Arriving', content: 'arriving', width: '10%'},
      ];
    },
  },
  methods: {
    resize(e: any): any {
      this.isMobile = e.target.innerWidth < 768 || this.isIntegrated;
    },
  },
});

window.addEventListener('storage', () => {
  location.reload();
});
</script>

<style lang="scss" scoped>
* {
  white-space: nowrap;
}

.parent {
  padding: 20px;
}

.container {
  width: 100%;
  transition: 200ms ease-in-out;
}

.caption {
  margin-bottom: 1em;
}

.table {
  position: relative;
  width: 100%;
}

.table-header {
  display: inline-block;
  border-bottom: 1px solid #eee;
  padding-left: 15px;
}

.table-row.stressed {
  background-color: #f8f8f8;
}

.table-row {
  padding: 5px 0;
  &:hover {
    background: #eee;
    cursor: pointer;
  }
}

.cell {
  display: inline-block;
  padding-left: 10px;
}
.modifiable-cell {
  //transition: 0.2s ease-in-out;
}
@keyframes collapse {
  0% {margin-top: 0; opacity: 1}
  100% {margin-top: -24px; opacity: 0}
}
@keyframes expand {
  0% {margin-top: -24px; opacity: 0}
  100% {margin-top: 0; opacity: 1;}
}
.expanded {
  animation: expand 200ms ease-in forwards;
}
.collapsed {
  animation: collapse 200ms ease-in forwards;
}
</style>
