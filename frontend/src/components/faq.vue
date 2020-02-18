<template>
  <div class="parent content">
    <h2>Frequently Asked Questions</h2>
    <h3>What is Shuttle Tracker?</h3>
    <p>
      Shuttle Tracker is
      <a href="https://github.com/wtg/shuttletracker">open source</a> and maintained by the
      <a href="https://webtech.union.rpi.edu">Web Technologies Group</a> of the Rensselaer Union Student Senate for the benefit of the student body.
      It is also an active 
      <a href="https://rcos.io">RCOS</a> project, and dozens of students have contributed to it through RCOS in previous semesters.
    </p>
    <h3>How can I contribute?</h3>
    <p>
      We have an active
      <a href="https://github.com/wtg/shuttletracker">repository</a> on GitHub where we track issues and accept pull requests. 
      If you want to get involved in another way, email
      <a href="mailto:webtech@union.lists.rpi.edu">webtech@union.lists.rpi.edu</a> for more information.
    </p>
    <h3>Why is there no app?</h3>
    <p>
        There is a Shuttle Tracker app in development right now, although we currently have no timeline to offer. 
        In the meantime, we're focused on making the Shuttle Tracker website as good as it can be.
    </p>
    <h3>Why are there shuttles in downtown Troy?</h3>
    <p>
        The shuttle station is located in downtown Troy. As shuttles are driving from the station to the Union to begin their route, they are still being tracked. 
        These are the shuttles you see in downtown Troy.
    </p>
    <h3>Why are the shuttles the wrong color?</h3>
    <p>
        There’s a few reasons that the shuttles you’re seeing could be wrong. Check the routes you have toggled on first. Different routes have different colors, 
        so the late night shuttle has a different color than the daytime shuttles. When a shuttle is off a route, it will be colored white. Secondly, some of the routes overlap, 
        so take a closer look. If you’re still experiencing problems, send us a message!
    </p>
    <h3>Why are some shuttles colored white?</h3>
    <p>
        If shuttles are off-route for multiple minutes, their color will change to a white color to indicate they are no longer on a route. 
    </p>
    <h3>Why do the shuttles jump on the map?</h3>
    <p>
        When the shuttle is traveling at a certain speed for 8 seconds, the shuttle location will update.
    </p>
    <h3>Why did the routes change?</h3>
    <p>
        What a complicated question! The short answer: Parking and Transportation, along with some Senate committees, changed the routes to make them shorter and more streamlined 
        (most notably, cutting the East route into two different routes). The long answer can be found <router-link to="/changes">here</router-link>.
    </p>
    <h3>Why are the shuttles spewing buses and tractor emojis?</h3>
    <p>
        A: The beloved 𝐵𝓊𝓈 𝐵𝓊𝓉𝓉𝑜𝓃 is used for a variety of reasons. You can press the 𝐵𝓊𝓈 𝐵𝓊𝓉𝓉𝑜𝓃 to let people know where a shuttle is, how fast you’re going, 
        or to assert your dominance.
    </p>
    <h3>How do you pronounce "track?"</h3>
    <p>
        “Tra-cc”, it's pretty straightforward.
    </p>
  </div>
</template>
<script lang="ts">
// This component is the faq window that is displayed when the `faq` button on the menu is selected

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
