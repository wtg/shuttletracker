<template>
  <div>
    <h2>Schedules</h2>
    <p>
      View the official shuttle schedules from the Parking and Transportation office.
    </p>
    <ul>
      <li>
        <a target="_blank" rel="noopener noreferrer" href="https://info.rpi.edu/sites/default/files/2018-19MonThursCampusShuttleScheduleEastRoute.pdf">East: Monday-Friday</a>
      </li>
      <li>
        <a target="_blank" rel="noopener noreferrer" href="https://info.rpi.edu/sites/default/files/2018-19MonThursCampusShuttleScheduleWestRoute.pdf">West: Monday-Friday</a>
      </li>
      <li>
        <a target="_blank" rel="noopener noreferrer" href="https://info.rpi.edu/sites/default/files/2018-19WeekendLateNightShuttleSchedule.pdf">Weekend/Late Night</a>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
export default Vue.extend({});
</script>

<style lang="scss" scoped>
div {
  padding: 20px;
}
ul {
  padding: 0;
}
li {
  list-style: none;
  padding: 0;
}
</style>
