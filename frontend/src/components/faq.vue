<template>
  <div class="parent content">
    <h2>Frequently Asked Questions</h2>
    <h4>What is Shuttle Tracker?</h4>
    <p>
      Shuttle Tracker is
      <a href="https://github.com/wtg/shuttletracker">open source</a> and maintained by the
      <a href="https://webtech.union.rpi.edu">Web Technologies Group</a> of the Rensselaer Union Student Senate for the benefit of the student body.
      It is also an active
      <a href="https://rcos.io">RCOS</a> project, and over 50 students have contributed to it through RCOS in previous semesters.
    </p>
    <h4>Why is there no app?</h4>
    <p>
        There is a Shuttle Tracker app in development right now, although we currently have no timeline to offer.
        In the meantime, we're focused on making the Shuttle Tracker website the best it can be.
    </p>
    <h4>Why are there shuttles in downtown Troy?</h4>
    <p>
        The shuttle station is located in downtown Troy. As shuttles are driving to and from the station each day, they are still being tracked.
        These are the shuttles you see in downtown Troy.
    </p>
    <h4>Why are the shuttles the wrong color?</h4>
    <p>
        Since both the shuttles and routes change on a daily basis, a shuttle's route is determined dynamically by looking at which route it has been closest to
        recently. However, due to routes overlapping and drivers going on break or off-course, a shuttle's route may be incorrectly determined. This
        should resolve itself within a few minutes, however if there are repeated issues feel free to send us a message.
    </p>
    <h4>Why are some shuttles colored white?</h4>
    <p>
        If shuttles are off-route for multiple minutes, their color will change to a white color to indicate they are no longer on a route.
    </p>
    <h4>Why do the shuttles jump on the map?</h4>
    <p>
        When the shuttle is traveling at a certain speed for 8 seconds, the shuttle location will update.
    </p>
    <h4>Why did the routes change?</h4>
    <p>
        Parking and Transportation, along with some Senate committees, updated the routes (beginning Spring 2020) to make them faster and serve more students
        (most notably by cutting the East route into two different routes). More information can be found <router-link to="/changes">here</router-link>.
    </p>
    <h4>How can I contribute?</h4>
    <p>
      We have an active
      <a href="https://github.com/wtg/shuttletracker">repository</a> on GitHub where we track issues and accept pull requests.
      Learn more about what we do on our <a href="https://webtech.union.rpi.edu">website</a> or by sending us an <a href="mailto:webtech@union.lists.rpi.edu">email</a>.
    </p>
    <h4>Why are the shuttles spewing buses and other emojis?</h4>
    <p>
        A: The beloved 𝐵𝓊𝓈 𝐵𝓊𝓉𝓉𝑜𝓃 can be used for a variety of reasons. You can press the 𝐵𝓊𝓈 𝐵𝓊𝓉𝓉𝑜𝓃 to let people know where a shuttle is, how fast you’re going,
        or to assert your dominance.
    </p>
    <h4>How do you pronounce "track?"</h4>
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
  p {
    margin-bottom: .8em;
  }
}
h4 {
  margin-bottom:.1em;
}
</style>
