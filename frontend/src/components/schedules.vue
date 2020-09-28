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
</template>

<script lang="ts">
import Vue from 'vue';
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
    };
  },
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
