<template>
  <div id='etas'>
    <div id='title'>ETAs</div>
      <div id='upcoming'>
        <ul id='queue'>
           {{unionEtas}}
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
        // Dictionary that stores the route ID with its corresponding route name
        routeNames: {19: 'Modified East Campus', 18: 'Modified West Campus', 1: 'West Campus', 15: 'East Campus'},
     };
  },
  computed: {
    message(): string | null {
      if (this.etaInfo === null) {
          return null;
      }
      const now = new Date();

      let newMessage = `${this.etaInfo.route.name} shuttle arriving at ${this.etaInfo.stop.name}`;
      // more than 1 min 30 sec?
      if (this.etaInfo.eta.eta.getTime() - now.getTime() > 1.5 * 60 * 1000 && !this.etaInfo.eta.arriving) {
        newMessage += ` in ${relativeTime(now, this.etaInfo.eta.eta)}`;
      }
      if (newMessage.substring(newMessage.length - 1) !== '.') {
        newMessage += '.';
      }
      return newMessage;
    },
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
    parseUnionEtas(allEtas) {
      const unionStops = [];
      // Loop through all 11 vehicles
      for (let i = 1; i < 12; i++) {
         const charRep = String(i);
         // Access the stop_etas for the vehicle using char i ('1' instead of 1)
         for (let j = 0; j < allEtas[charRep].stop_etas.length; j++) {
            // The union's stop_id is 1
            if (allEtas[charRep].stop_etas[j].stop_id === 1) {
               /* Create a tuple to hold the route name (looked up in the dictionary)
                  and the time this shuttle is arriving at the union */
               let routeNameAndEta: [string, string];
               routeNameAndEta = [this.routeNames[allEtas[charRep].route_id], allEtas[charRep].stop_etas[j].eta];
               unionStops.push(routeNameAndEta);
            }
         }
      }
      return unionStops;
    }
    // Get request to retrieve eta data from localhost:9000/eta
    retrieveEtaData() {
      axios.get('http://localhost:9000/eta', {})
      .then((res) => {
         const key = 'data';
         res = res[key];
         return res;
      })
      .then((res) => {
         this.unionEtas = this.parseUnionEtas(res);
      });
   },
  },
  // Allow this function to be called on page load
  mounted() {
    this.changeTextColor();
    this.retrieveEtaData();
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
    font-size: 50px;
    display: inline-block;
    border-bottom: 2px solid #4c4c4c;
  }
  li {
  border: 1px solid rgb(228, 228, 228);
  border-radius: 3px;
  font-size: 20px;
  margin: 10px;
  padding: 10px;
}
// West color = #0080FF
// Weekend late night color = #9b59b6
// East inclement weather color = #ff9900
// West inclement weather color = #FF0
// East color = #96C03A


</style>
