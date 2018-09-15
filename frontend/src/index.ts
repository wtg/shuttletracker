import Vue from 'vue';
import Router from 'vue-router';
import Public from './components/Public.vue';
import about from './components/about.vue';

Vue.use(Router);

export default new Router({
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
    ],
  });
