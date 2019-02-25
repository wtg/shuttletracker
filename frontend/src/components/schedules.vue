<template>
  <div class="parent">
    <h2 class="title">Schedules</h2>
    <p class="subtitle">
      Official shuttle schedules from the Parking and Transportation office.
    </p>
    <hr>
    <div class="columns">
      <div class="column" v-for="link in links" v-bind:key="link.url">
        <div class="card">
          <header class="card-header">
            <p class="card-header-title">
              {{ link.name }}
            </p>
          </header>
          <div class="card-content">
            <p>{{ link.caption }}</p>
            <a class="" target="_blank" rel="noopener noreferrer" v-bind:href="link.url">View PDF</a>
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
          url: 'https://info.rpi.edu/sites/default/files/2018-19MonThursCampusShuttleScheduleEastRoute.pdf',
          name: 'East Campus',
          caption: 'Monday–Friday, 7 am–11 pm',
          color: 'green',
        },
        {
          url: 'https://info.rpi.edu/sites/default/files/2018-19MonThursCampusShuttleScheduleWestRoute.pdf',
          name: 'West Campus',
          caption: 'Monday–Friday, 7 am–11 pm',
          color: 'red',
        },
        {
          url: 'https://info.rpi.edu/sites/default/files/2018-19WeekendLateNightShuttleSchedule.pdf',
          name: 'Weekend/Late Night',
          caption: 'Friday–Sunday, 8 pm–4 am',
          color: 'purple',
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
</style>
