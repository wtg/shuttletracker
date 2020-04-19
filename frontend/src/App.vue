<template>
  <div id="app">
    <div id="view-wrapper">
      <keep-alive>
        <router-view id="router-view"></router-view>
      </keep-alive>
    </div>
    <tab-bar></tab-bar>   
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import TabBar from '@/components/tabBar.vue';
import UserLocationService from '@/structures/userlocation.service';
import {DarkTheme} from '@/structures/theme';

UserLocationService.getInstance();

export default Vue.extend({
  name: 'app',
  components: {
    TabBar,
  },
  computed: {
    // The theme code is in this file rather than Public.vue because we need the theme to apply to the tab bar and
    // router tabs as well.
    currentCSSTheme(): string {
      return DarkTheme.getCurrentCSSThemeAttribute(this.$store.state);
    },
  },
  watch: {
    // A watcher is used rather than a bind because we need to put the `data-theme` attribute on <body>, and Vue
    // cannot bind to the <body> element.
    // Why? Because inherited CSS properties whose values are a var(...) expression use the context of the inherited
    // source, not of the inheritor. An example is the color property. <body> has `color: var(...)`. Text elements like
    // <p> will inherit their color property from <body>. If data-theme is placed below <body>, those <p> elements will
    // NOT see the updated values of the theme variables because their inherited property comes from <body>.
    currentCSSTheme: {
      handler(newValue: string) {
        document.body.setAttribute('data-theme', newValue);
      },
      immediate: true,
    },
  },
});
</script>
<style lang="scss">
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
