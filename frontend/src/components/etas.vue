<template>
<div class="parent">
    <h1 class="title">ETAs</h1>
    <hr>

    <div class="container">
        <table class="table">
            <thead>
                <tr>
                    <th>Vehicles</th>
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
                         <td v-if="this.colors.get(info.routeID) == '#0080FF' " > 
                            <font style="color:#0080FF"> &#9830 </font> {{ info.vehicleID }}
                            </td>
                        <td v-else-if="this.colors.get(info.routeID) == '#9b59b6' " > 
                            <font style="color:#9b59b6"> &#9830 </font> {{ info.vehicleID }}
                            </td>
                        <td v-else-if="this.colors.get(info.routeID) == '#96C03A' " > 
                            <font style="color:#96C03A"> &#9830 </font> {{ info.vehicleID }}
                            </td>
                        <td v-else-if="this.colors.get(info.routeID) == '#DCC308' " > 
                            <font style="color:#DCC308"> &#9830 </font> {{ info.vehicleID }}
                            </td>   
                         <td v-else-if="this.colors.get(info.routeID) == '#FF0000' " > 
                            <font style="color:#FF0000"> &#9830 </font> {{ info.vehicleID }}
                            </td> 
                         <td v-else-if="this.colors.get(info.routeID) == '#0080FF' " > 
                            <font style="color:#0080FF"> &#9830 </font> {{ info.vehicleID }}
                            </td> 
                        <td v-else-if="this.colors.get(info.routeID) == '#aa08dc' " > 
                            <font style="color:#aa08dc"> &#9830 </font> {{ info.vehicleID }}
                            </td> 
                        <td v-else-if="this.colors.get(info.routeID) == '#00ff00' " > 
                            <font style="color:#00ff00"> &#9830 </font> {{ info.vehicleID }}
                            </td> 
                        <td v-else-if="this.colors.get(info.routeID) == '#ff9d00' " > 
                            <font style="color:#ff9d00"> &#9830 </font> {{ info.vehicleID }}
                            </td> 
                        <td v-else-if="this.colors.get(info.routeID) == '#ffa500' " > 
                            <font style="color:##ffa500"> &#9830 </font> {{ info.vehicleID }}
                            </td> 
                        <td v-else>{{info.vehicleID}}</td>
                        <td>{{ this.stops.get(info.stopID) }}</td>
                        <td>{{ info.eta }}</td>
                        <td>{{ info.arriving }}</td>
                        <td>{{ this.posts.get(info.routeID) }}</td>
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
            posts: new Map(),
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
                console.log(dictionary);
                this.posts = dictionary;
                console.log(this.posts.get(21));
                console.log(dictionary2);
                this.colors = dictionary2;
                console.log(this.colors.get(21));

                // this.posts = response.data;
                // console.log(this.posts);
                // console.log(this.posts[0].name);
            })
            .catch((error) => console.log(error));
        axios
            .get('http://localhost:8080/stops')
            .then((response) => {
                const dictionary = new Map();
                response.data.forEach((x: any) => {
                    dictionary.set(x.id, x.name);
                });
                console.log(dictionary);
                this.stops = dictionary;
                console.log(this.stops.get(21));
                // this.posts = response.data;
                // console.log(this.posts);
                // console.log(this.posts[0].name);
            })
            .catch((error) => console.log(error));
    },
    computed: {
        etas(): any[] {

            const etaArray = [];
            for (let i = 0; i < 18; i++) {
                const etaString = localStorage.getItem(String(i + 1));
                if (etaString) {
                    const localETA = JSON.parse(etaString);
                    const ret = [];
                    if (localETA.length) {
                        for (const eta of localETA) {
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
                                routeID: eta.routeID,
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
