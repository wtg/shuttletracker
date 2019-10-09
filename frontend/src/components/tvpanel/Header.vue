<template>
  <div class='header'>
    <div id="header-content">
      <div id='shuttle-tracker'>
        Shuttle Tracker
      </div>
      <img id='shuttletrackerlogo' src="../../assets/icon.svg">
      <div id="weblink">shuttles.rpi.edu</div>
    </div>
  </div>
</template>

<script lang="ts">

// This component is the Header component for the tv panel
import Vue from 'vue';

export default Vue.extend({
  name: 'Header',
});
</script>

<style lang="scss" scoped>

.header {
  height: 10%;
  width: 70%;
  position: absolute;
  top: 0px;
  left: 0px;
  background-color: #F8F8F8;
  display: flex;
  justify-content: center;
  align-items: center;
  color: black;
}

#shuttle-tracker {
  display: inline-block;
  margin-right: 10px;
  margin-left: 10px;
  font-size: 40px;
}

#weblink {
  display: inline-block;
  margin-right: 10px;
  margin-left: 10px;
  font-size: 40px;
}

#shuttletrackerlogo {
  display: inline-block;
  height: 40px;
  width: 40px;
  margin-right: 10px;
  margin-left: 10px;
}

#header-content {
  padding: 1.35% 0;
  display: inline-block;
}

</style>
