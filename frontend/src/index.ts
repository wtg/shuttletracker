import Vue from 'vue';
import Router from 'vue-router';
import Public from './components/Public.vue';
import about from './components/about.vue';
import tvcomponent from './components/tv.vue';

Vue.use(Router);

export default new Router({
    mode: 'history',
    routes: [
      {
        path: '/',
        name: 'main',
        component: Public,
      },
      {
        path: '/about',
        name: 'about',
        component: about,
      },
      {
        path: '/tv',
        name: 'TV Component',
        component: tvcomponent,
      }
    ],
  });
