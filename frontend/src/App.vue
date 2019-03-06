<template>
  <div id="app">
    <div id="view-wrapper">
      <keep-alive>
        <router-view id="router-view"></router-view>
      </keep-alive>
    </div>  
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import UserLocationService from '@/structures/userlocation.service';

UserLocationService.getInstance();

export default Vue.extend({
  name: 'app',
});
</script>
<style lang="scss">
@import "~bulma/sass/utilities/_all";

$primary: #ED1C24;
$link: red;
$primary-invert: findColorInvert($primary);
$twitter: #4099FF;
$twitter-invert: findColorInvert($twitter);

$colors: (
    "white": ($white, $black),
    "black": ($black, $white),
    "light": ($light, $light-invert),
    "dark": ($dark, $dark-invert),
    "primary": ($primary, $primary-invert),
    "info": ($info, $info-invert),
    "success": ($success, $success-invert),
    "warning": ($warning, $warning-invert),
    "danger": ($danger, $danger-invert),
    "twitter": ($twitter, $twitter-invert)
);


@import "~bulma";
@import "~buefy/src/scss/buefy";
#app {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
#view-wrapper {
  height: calc(100% - 40px);
  box-sizing: border-box;
}
#router-view {
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  height: 100%;
}
</style>
