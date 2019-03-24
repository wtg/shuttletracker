<template>
  <!-- Display Current Time and Day -->
  <div id="timeDisplay">
    <h1>Current Time: {{displayTime()}}</h1>
    <h1>Today is {{displayDay()}}</h1>
  </div>
</template>

<script>
// This component handles the Time Display on the TV panel
import Vue from "vue";

export default Vue.extend({
  name: "TimeDisplay",
  methods: {
    // Function to display the Current Time
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

<style scoped>
#timeDisplay {
  position: absolute;
  right: 200px;
  top: 70px;
}
</style>
