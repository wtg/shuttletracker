<template>
  <ul>
    <router-link tag="li" to="/">
      <span>
        <span class="fas fa-location-arrow"></span>
        Map
      </span>
    </router-link>
    <router-link tag="li" to="/schedules">
      <span>
        <span class="fas fa-list-ul"></span>
        Schedules
      </span>
    </router-link>
    <!-- <router-link v-if="etasEnabled" tag="li" to="/etas">
      ETAs
    </router-link> -->
    <router-link tag="li" to="/settings">
      <span>
        <span class="fas fa-cog"></span>
        Settings
      </span>
    </router-link>
    <b-field>
      <b-switch v-model="pushEnabled">Push Notification</b-switch>
    </b-field>
  </ul>
</template>


<script lang="ts">
/*$ npm install -g web-push --save
$ web-push generate-vapid-keys */
import EventBus from '@/event_bus';
import Vue from 'vue';

export default Vue.extend({
  computed: {
    etasEnabled(): boolean {
      return this.$store.state.settings.etasEnabled;
    },
    pushEnabled: {
      get(): boolean {
        return this.$store.state.settings.pushEnabled;
      },
      set(value: boolean) {
        if ( value === true ) {
          this.askPermission();
          const payload = {
            register: this.registration,
          };
          EventBus.$emit('REGISTER', payload);
          /*const pushSubscription = this.subscribeUserToPush();
          const subscriptionObject = JSON.stringify(pushSubscription);*/
        }
        this.$store.commit('setSettingsPushEnabled', value);
      },
    },
    registration() {
      if (!('serviceWorker' in navigator)) {
        console.log('ServiceWorker isn\'t supported');
      }
      if (!('PushManager' in window)) {
        console.log('Push isn\'t supported');
      }
      return this.registerServiceWorker();
    },
  },
  methods: {
    registerServiceWorker() {
      navigator.serviceWorker.register('/serviceworker.js')
      .then((reg) => {
        console.log('Service Worker Successfully Registered.');
        return reg;
      })
      .catch((err) => {
        console.error('Unable to Register Service Worker.', err);
      });
    },
    askPermission() {
      new Promise((resolve, reject) => {
        const permissionResult = Notification.requestPermission((result) => {
          resolve(result);
        });
        if ( permissionResult ) {
          permissionResult.then(resolve, reject);
        }
      })
      .then((permissionResult) => {
        if ( permissionResult !== 'granted' ) {
          this.pushEnabled = false;
          console.log('No permission granted');
        }
      });
    },
    subscribeUserToPush() {
      navigator.serviceWorker.register('serviceworker.js')
      .then((reg) => {
        const options = {
          userVisibleOnly: true,
          applicationServerKey: this.urlBase64ToUint8Array('BFgDYavD_PqQryHBqpHa6w7Fh8dRok9tZ3E5nfY0_P3cFCNPSHN06REG-kgtEjuKkQE_UZp3bjREFKPtQR5wVqk'),
        };
        reg.pushManager.subscribe(options);
      })
      .then((pushSubscription) => {
        console.log('Subscribed Push:', JSON.stringify(pushSubscription));
      });
    },
    urlBase64ToUint8Array(base64String: string) {
      const padding = '='.repeat((4 - base64String.length % 4) % 4);
      const base64 = (base64String + padding)
      .replace(/\-/g, '+')
      .replace(/_/g, '/');

      const rawData = window.atob(base64);
      const outputArray = new Uint8Array(rawData.length);

      for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i);
      }
      return outputArray;
    },
  },
});
</script>

<style lang="scss" scoped>
@import "@/assets/vars.scss";
ul {
  display: flex;
  height: 40px;
  position: fixed;
  bottom: 0;
  display: flex;
  padding: 0;
  margin: 0;
  justify-content: center;
  align-items: center;
  width: 100%;
  border-top: 0.5px solid #eee;
  box-shadow: 0 3px 8px 0 #ddd;
  font-size: 13px;
  user-select: none;
  background: white;
}
li {
  cursor: pointer;
  width: auto;
  height: 100%;
  padding: 5px 15px;
  margin: 0 5px;
  border-top: 1px solid rgba(0, 0, 0, 0);
  position: relative;
  top: -0.5px;
  display: flex;
  align-items: center;
}
.router-link-exact-active {
  border-top: 1px solid $primary;
  color: $primary;
}
</style>
