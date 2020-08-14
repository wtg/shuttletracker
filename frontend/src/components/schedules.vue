<template>
  <div class="parent">
    <h1 class="title">Schedules</h1>
    <h2 class="subtitle">Official shuttle schedules from the Parking and Transportation office.<hr></h2>
    <div class="columns">
      <div class="column" v-for="link in links" v-bind:key="link.url">
        <div class="dimmed">
          <div class="box">
            <strong class="link-header">{{ link.name }}</strong>
            <hr>
            <div class="box-section">
              <div class="caption">
                <p v-for="line in link.caption" v-bind:key="line">{{ line }}</p>
              </div>
              <div style="color: var(--color-primary)">No paper schedules this semester</div>
              <!-- <a target="_blank" rel="noopener noreferrer" v-bind:href="link.url">View PDF</a> -->
            </div>
          </div>
          
        </div>
      </div>
      </div>
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
      links: [
        {
          // url: 'https://shuttles.rpi.edu/static/Weekday.pdf',
          name: 'Weekday Routes',
          caption: [
            'North and West Routes',
            'Monday–Friday 7am – 11:45pm',
            '⠀ ',
            'Hudson Valley College Suites Shuttles',
            'Monday – Friday 7am – 7pm',
            '⠀',
            'CDTA Express Route',
            'Monday–Friday 7am – 7pm',
          ],
          color: 'green',
        },
        {
          // url: 'https://shuttles.rpi.edu/static/Weekend.pdf',
          name: 'Weekend Routes',
          caption: [
            'North and West Routes',
            'Saturday 9am – 11:45pm',
            'Sunday 9am – 8pm',
            '⠀ ',
            'Hudson Valley College Suites Shuttles',
            'Saturday–Sunday 7am – 7pm',
          ],
          color: 'red',
        },
      ],
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
.link-header {
  margin-left: 1em;
}
.box-section {
  margin-left: 1em;
}

hr {
  margin: 1rem 0;
}

</style>
