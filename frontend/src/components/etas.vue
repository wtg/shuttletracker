<template>
    <div class="parent">
        <h1 v-if="!isIntegrated" class="title">ETAs</h1>
        <hr v-if="!isIntegrated">

        <div :class="isMobile ? '' : 'container'">
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

                <template v-for="(eta, i) in etas"><!-- :key="`${i}-${eta.stopID}`"-->
                    <tr @click="selection=i" class="table-row"
                        :class="[i % 2 === 1 ? 'stressed' : '']">
                        <td>
                            <div :style="{ color: eta.routeColor, display: 'inline'}"> &#9830;</div>
                            {{ eta.routeName }}
                        </td>
                        <td v-show="!isMobile">{{ eta.vehicleID }}</td>
                        <td v-show="!isMobile">{{ eta.stopName }}</td>
                        <td>{{ eta.eta }}</td>
                        <td>{{ eta.arriving }}</td>
                    </tr>
                    <tr :class="selection === i && isMobile ? 'extension-row-active' : ''" class="extension-row">
                        <td colspan="2">Stop: {{ eta.stopName }}</td>
                        <td>
                            V.ID: {{ eta.vehicleID }}
                        </td>
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

.table-row.stressed {
    background-color: #f8f8f8;
}

.table-row {
    &:hover {
        background: #eee;
        cursor: pointer;
    }
    &-selected {
        row-span: 2;
        //height: 100px;
    }
}
.extension-row {
    display: none;
    &-active {
        display: table-row;
    }
}
</style>
