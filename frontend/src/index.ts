import Vue from 'vue';
import Router from 'vue-router';
import Public from './components/Public.vue';
import about from './components/about.vue';
import faq from './components/faq.vue';
import schedules from '@/components/schedules.vue';
import settings from '@/components/settings.vue';
import etas from '@/components/etas.vue';
import changes from '@/components/changes.vue';
import Resources from '@/resources';
import feedback from './components/feedback.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: Resources.BasePath,
  routes: [
    {
      path: '/',
      name: 'map',
      component: Public,
    },
    {
      path: '/about',
      name: 'about',
      component: about,
    },
    {
      path: '/faq',
      name: 'faq',
      component: faq,
    },
    {
      path: '/etas',
      name: 'etas',
      component: etas,
    },
    {
      path: '/schedules',
      name: 'schedules',
      component: schedules,
    },
    {
      path: '/settings',
      name: 'settings',
      component: settings,
    },
    {
      path: '/changes',
      name: 'changes',
      component: changes,
    },
    {
      path: '/feedback',
      name: 'feedback',
      component: feedback,
    },
  ],
});
