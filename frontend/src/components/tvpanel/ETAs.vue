<template>
  <div id='etas'>
    <div id='title'>ETAs</div>
      <div id='upcoming'>
        <ul id='queue'>
          <li v-if="unionEtas[0]">{{displayTime(unionEtas[0])}}</li>
          <li v-if="unionEtas[1]">{{displayTime(unionEtas[1])}}</li>
          <li v-if="unionEtas[2]">{{displayTime(unionEtas[2])}}</li>
          <li v-if="unionEtas.length === 0"> No Shuttles Avaliable</li>
        </ul>
      </div>
  </div>
</template>

<script lang='ts'>

// This component is the ETAs component for the tv panel
import Vue from 'vue';
import ETA from '@/structures/eta';
import axios from 'axios';

export default Vue.extend({
  props: ['etaInfo', 'show'],
  data() {
     return {
        unionEtas: [],
        routeNames: {19: 'Modified East Campus', 18: 'Modified West Campus', 1: 'West Campus', 15: 'East Campus'},
     } as {
        unionEtas: string[][],
        routeNames: {[id: number]: string; },
     };
  },
  methods: {
    changeTextColor() {
      /* Change the color of the text depending on if the
        display is for the East or West shuttle */
      const liElems = document.getElementsByTagName('li');
      for (let i = 0; i < liElems.length; i++) {
        const etaText = liElems[i].innerHTML;
        const subStr = etaText.substring(0, 4);
        liElems[i].style.fontWeight = 'bold';
        if (subStr === 'West') {
          liElems[i].style.color = '#0080FF';
        }
        if (subStr === 'East') {
          liElems[i].style.color = '#71922b';
        }
      }
    },
    // Function to retrieve union vehicle stops from all ETAS
    parseUnionEtas(allEtas: any): string[][] {
      let unionStops: string[][];
      unionStops = [];
      if (allEtas.length === 0) {
        return unionStops;
      }
      for (let i = 1; i < 12; i++) {
         const charRep = String(i);
         for (let j = 0; j < allEtas[charRep].stop_etas.length; j++) {
            if (allEtas[charRep].stop_etas[j].stop_id === 1) {
               let routeNameAndEta: string[];
               routeNameAndEta = [this.routeNames[allEtas[charRep].route_id], allEtas[charRep].stop_etas[j].eta];
               unionStops.push(routeNameAndEta);
            }
         }
      }
      return unionStops;
    },
    // Get request to retrieve eta data from the "eta" endpoint
    retrieveEtaData() {
      axios.get('/eta', {})
      .then((res) => {
         const key = 'data';
         res = res[key];
         return res;
      })
      .then((res) => {
         this.unionEtas = this.parseUnionEtas(res);
      });
    },
    // Function to find the estimated time of arrival for each ETA of the union
    displayTime(currETA: string[]): string {
      let message: string;
      const currTimeDate = Date.parse(currETA[1]);
      const now = new Date();
      message = currETA[0] + ' arriving ' + this.relativeTime(now, currTimeDate);
      return message;
    },
    // Function to calculate the estimated time in minutes from two Date objects
    relativeTime(from: Date, to: number): string {
      const minuteMs = 60 * 1000;
      const elapsed = to - from.getTime();
      const minutesTillArrival = Math.round(elapsed / minuteMs);
      if (minutesTillArrival === 0) {
         return 'now';
      }
      return `in ${minutesTillArrival} minutes`;
    },
  },
  // Allow this function to be called on page load
  mounted() {
    this.changeTextColor();
    this.retrieveEtaData();
    setInterval(() => {
      this.retrieveEtaData();
    }, 30000);
  },
});

</script>

<style lang='scss' scoped>
  #etas {
    height: 45%;
    width: 100%;
    text-align: center;
    border-top: 3px solid #F8F8F8;
    position: relative;
    top: 10%;
  }
  #title {
    margin-top: 20px;
    margin-bottom: 8px;
    font-size: 50px;
    display: inline-block;
    border-bottom: 2px solid #4c4c4c;
  }
  li {
    border: 1px solid rgb(228, 228, 228);
    border-radius: 3px;
    margin-left: 10px;
    margin-right: 10px;
    margin-top: 18px;
    margin-bottom: 18px;
    padding: 10px;
    font-weight: 700;
  }
// West color = #0080FF
// Weekend late night color = #9b59b6
// East inclement weather color = #ff9900
// West inclement weather color = #FF0
// East color = #96C03A
</style>
