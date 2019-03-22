<template>
  <div>
    <p id="title-RPI">Shuttle Tracker</p>

    <!-- Display Current Time and Day -->
    <div id="timeDisplay">
      <h1>Current Time: {{displayTime()}}</h1>
      <h1>Today is {{displayDay()}}</h1>
    </div>

    <TimesQueue/>
    <Map/>
    <div class="titleBar bot">
      <img src="~../assets/icon.svg" id="shuttlepic">
      <h1 id="WTG">Web Technologies Group</h1>
    </div>
  </div>
</template>

<script lang="ts">
// This component handles the main TV Panel Application
import Vue from "vue";
import Map from "./tvpanel/Map.vue";
import TimesQueue from "./tvpanel/TimesQueue.vue";
export default Vue.extend({
  name: "tvpanel",
  components: {
    TimesQueue,
    Map
  },
  methods: {
    // Function to display the time
    displayTime() {
      let today = new Date();
      let hour = today.getHours();
      let minutes = today.getMinutes();
      let am_pm;

      // Format Hours
      if (hour < 12) {
        if (hour == 0) hour = 12;
        am_pm = "AM";
      } else if (hour == 12) am_pm = "PM";
      else {
        if (hour == 24) {
          hour = hour - 12;
          am_pm = "AM";
        }
        hour = hour - 12;
        am_pm = "PM";
      }
      // Format Minutes and return
      if (minutes <= 9) {
        return hour.toString() + ":0" + minutes.toString() + " " + am_pm;
      }
      return hour.toString() + ":" + minutes.toString() + " " + am_pm;
    },

    // Function to display the day of the week
    displayDay() {
      let today = new Date();
      let day = today.getDay();

      switch (day) {
        case 0:
          return "Sunday";
        case 1:
          return "Monday";
        case 2:
          return "Tuesday";
        case 3:
          return "Wednesday";
        case 4:
          return "Thursday";
        case 5:
          return "Friday";
        case 6:
          return "Saturday";
      }
    }
  },

  mounted() {
    // Update time every minute
    setInterval(() => {
      this.updateCurTime();
    }, 60000);
  }
});
</script>

 
<style>
html,
body {
  height: 100%;
  width: 100%;
  overflow: hidden;
}
#title-RPI {
  text-align: center;
  position: relative;
  top: 100px;
  font-size: 75px;
  color: black;
}
#WTG {
  position: absolute;
  right: 20px;
}
.bot {
  position: absolute;
  bottom: 0px;
  height: 60px;
}
#timeDisplay {
  position: absolute;
  right: 200px;
  top: 70px;
}
</style>
