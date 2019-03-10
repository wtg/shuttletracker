<template>
  <div class="parent content">
    <h2>About Shuttle Tracker</h2>
    <p>
      Shuttle Tracker is
      <a href="https://github.com/wtg/shuttletracker">open source</a> and maintained by the
      <a href="https://webtech.union.rpi.edu">Web Technologies Group</a> of the Rensselaer Union Student Senate for the benefit of the student body.
      It is also an active
      <a
        href="https://rcos.io"
      >RCOS</a> project, and dozens of students have contributed to it through RCOS in previous semesters.
    </p>
    <h3>Why is there no app?</h3>
    <p>We're focused on making the Shuttle Tracker website as good as it can be. While most of our users are on mobile, we think the best way to deliver a great experience to the most people is to focus our efforts on this website.</p>
    <h3>How can I contribute?</h3>
    <p>
      We have an active
      <a href="https://github.com/wtg/shuttletracker">repository</a> on GitHub where we track issues and accept pull requests. If you want to get involved in another way, email
      <a
        href="mailto:webtech@union.lists.rpi.edu"
      >webtech@union.lists.rpi.edu</a> for more information.
    </p>

    <h3>Privacy policy</h3>
    <p>
      Shuttle Tracker is operated by the
      <a
        href="https://webtech.union.rpi.edu"
      >Web Technolgies Group</a> (WebTech), a committee of the Rensselaer Union Student Senate. WebTech uses Google Analytics to gather anonymous metrics about the users of its services. This information cannot and will not be used to identify you or any specific user of the service.
    </p>
    <p>If you grant permission for Shuttle Tracker to access your location, this information is used for two purposes. The first is to indicate your location on the map. The second is to enhance the accuracy of vehicle tracking. The information Shuttle Tracker gathers is limited to your device's latitude, longitude, speed, and heading. These data are associated with a random identifier that is generated whenever you open Shuttle Tracker. In order to protect your privacy, no two visits to Shuttle Tracker are associated. This identifier is not used to identify any specific user of the service. All data gathered are only analyzed in aggregate in order to improve the quality of vehicle tracking for all users.</p>
    <p>
      Any questions about this privacy policy should be directed to
      <a
        href="mailto:webtech@union.lists.rpi.edu"
      >webtech@union.lists.rpi.edu</a>.
    </p>
  </div>
</template>
<script lang="ts">
// This component is the about window that is displayed when the `about` button on the menu is selected

import Vue from 'vue';
export default Vue.extend({});
</script>
<style lang="scss" scoped>
.parent {
  padding: 20px;
}
p {
  margin-bottom: 1em;
}
</style>
