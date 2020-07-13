<template>
  <div class="parent">
    <h1 class="title">Schedules</h1>
    <p class="subtitle">Official shuttle schedules from the Parking and Transportation office.</p>
    <hr>
    <div class="columns">
      <div class="column" v-for="link in links" v-bind:key="link.url">
        <div class="box dimmed">
          <strong class="link-header">{{ link.name }}</strong>
          <hr />
          <div class="box-section">
            <div class="caption">
              <p v-for="line in link.caption" v-bind:key="line">{{ line }}</p>
            </div>
            <a target="_blank" rel="noopener noreferrer" v-bind:href="link.url">View PDF</a>
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
          url: 'https://shuttles.rpi.edu/static/Weekday.pdf',
          name: 'Weekday Routes',
          caption: [
            'North, South, and New West Routes',
            'Monday–Friday 7am – 11pm',
          ],
          color: 'green',
        },
        {
          url: 'https://shuttles.rpi.edu/static/Weekend.pdf',
          name: 'Weekend Routes',
          caption: [
            'West and East Routes',
            'Saturday–Sunday 9:30am – 5pm',
            '⠀',
            'Weekend Express Route',
            'Saturday–Sunday 4:30pm – 8pm',
            '⠀',
            'Late Night Route',
            'Friday–Saturday 8pm – 4am',
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

.dimmed {
  position: relative;
  border-radius: 10px;
}

.dimmed:after {
  content: " ";
  z-index: 10;
  display: block;
  position: absolute;
  height: 100%;
  top: 0;
  left: 0;
  right: 0;
  border-radius: 10px;
  background: var(--color-box-color);
}
</style>
