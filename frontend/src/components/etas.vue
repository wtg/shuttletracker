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
                    <th>Route ID</th>
                </tr>
            </thead>
            <tbody>
                <tr v-if="!etas.length">
                    <td colspan="3"> No ETAS Currently Calculated </td>
                </tr>
                    
                <template v-for="(eta) in etas">
                    <tr v-for="(info, i) in eta" v-bind:key="`${i}-${info.stopID}`">
                        <td> 
                            {{ info.vehicleID }}
                        </td>
                        <td>{{ info.stopName }}</td>
                        <td>{{ info.eta }}</td>
                        <td>{{ info.arriving }}</td>
                        <td><div :style="{ color: info.routeColor, display: inline}"> &#9830; </div> {{ info.routeName }}</td>
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
import axios from 'axios';

export default Vue.extend({
    data() {
        return {
            routes: new Map(),
            colors: new Map(),
            stops: new Map(),
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
    },
    computed: {
        etas(): any[] {

            const etaArray = [];
            for (let i = 0; i < 18; i++) {
                const etaString = localStorage.getItem(String(i + 1));
                // const testString = '[{"stopID": 10, "vehicleID": 3, "routeID": 22, "eta": "2020-10-28T01:09:09.826Z", "arriving": true}]';
                // console.log(JSON.parse(testString));
                if (etaString) {
                    const localETA = JSON.parse(etaString);
                    const ret = [];
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
                                ret.push(e);
                            }
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
