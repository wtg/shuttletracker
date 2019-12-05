<template>
  <!-- Display Current Time and Day -->
  <div id='time-display'>
    {{displayDay() + ', '+ displayTime()}}
  </div>
</template>

<script>
// This component handles the Time Display on the TV panel
import Vue from 'vue';
export default Vue.extend({
  name: 'TimeDisplay',
  data() {
    return {
      today: new Date(),
    };
  },
  methods: {
    // Function to display the Current Time
    displayTime() {
      this.today = new Date();
      let hour = this.today.getHours();
      const minutes = this.today.getMinutes();
      let am_pm;
      // Format Hours
      if (hour < 12) {
        if (hour === 0) {
          hour = 12;
        }
        am_pm = 'AM';
      } else if (hour === 12) {
        am_pm = 'PM';
      } else {
        if (hour === 24) {
          hour = hour - 12;
          am_pm = 'AM';
        }
        hour = hour - 12;
        am_pm = 'PM';
      }
      // Format minutes and return
      if (minutes <= 9) {
        return hour.toString() + ':0' + minutes.toString() + ' ' + am_pm;
      }
      return hour.toString() + ':' + minutes.toString() + ' ' + am_pm;
    },
    // Function to display the day of the week
    displayDay() {
      const today = new Date();
      const day = today.getDay();
      switch (day) {
        case 0:
          return 'Sunday';
        case 1:
          return 'Monday';
        case 2:
          return 'Tuesday';
        case 3:
          return 'Wednesday';
        case 4:
          return 'Thursday';
        case 5:
          return 'Friday';
        case 6:
          return 'Saturday';
      }
    },
  },
  mounted() {
    // Update Time and Day every minute
    setInterval(() => {
      this.displayTime();
      this.displayDay();
    }, 30000);
  },
});
</script>

<style scoped>
  #time-display {
    width: 100%;
    height: 15%;
    position:absolute;
    text-align: center;
    margin-top: 22px;
    font-size: 40px;
    font-weight: 600;
    top: 0;
  }
</style>
