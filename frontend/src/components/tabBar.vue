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
      <button v-on:click="sendnotify()" :disabled="valid == false">Notify</button>
    </b-field>
  </ul>
</template>


<script lang="ts">
/*$ npm install -g web-push --save
$ web-push generate-vapid-keys */
// Public: BP_qB6Rfxb4PNwV89br7bq5WzXoEU5pJwvS_wji6iNgEOTYo2MiVNhmBq6zDMg2HPjNUr1MamHvEFttADuLni2g
import EventBus from '@/event_bus';
import Vue from 'vue';
import axios from 'axios';

export default Vue.extend({
  data() {
    return {
      valid: true,
      newMessage: '',
      eta: 0,
    };
  },
  computed: {
    etasEnabled(): boolean {
      return this.$store.state.settings.etasEnabled;
    },
  },
  created() {
    if (!('serviceWorker' in navigator)) {
      console.log('ServiceWorker isn\'t supported');
      this.valid = false;
    } else if ( Notification.permission === 'denied' ) {
      console.log('User has blocked notifications');
      this.valid = false;
    } else if (!('PushManager' in window)) {
      console.log('Push isn\'t supported');
      this.valid = false;
    } else {
      this.registerServiceWorker();
      this.readyServiceWorker();
      EventBus.$on('PUSH', (payload: any) => {
        this.eta = payload.eta;
        this.newMessage = payload.newMessage;
      });
    }
  },
  methods: {
    subscribe() {
      navigator.serviceWorker.ready
      .then((registration) => {
        const vapidPublicKey = 'BHu_01FAmOhIaQ1KXX4qqHiJ7ire9s5dYTK4TF2dFXbeWb0fFvfpjJl3zaQjonIjhx1bl7IlQ_MWFsQBzAYZV9I';
        return registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey: this.urlBase64ToUint8Array(vapidPublicKey),
        });
      })
      .then((subscription) => {
        console.log(
          JSON.stringify({
            subscript: subscription,
          }),
        );
      })
      .catch((err) => console.error(err));
    },

    readyServiceWorker() {
      navigator.serviceWorker.ready
      .then((registration) => {
        return registration.pushManager.getSubscription()
        .then((subscription) => {
          if ( subscription ) {
            console.log(JSON.stringify({
              subscript: subscription,
            }));
          } else {
          // subscribe
            this.subscribe();
          }
        });
      }).then((subscription) => {
        console.log(JSON.stringify({
          subscript: subscription,
        }));
      });
    },

    registerServiceWorker() {
      navigator.serviceWorker.register('/serviceworker.js')
      .then((reg) => {
        console.log('Service Worker Successfully Registered.');
      })
      .catch((err) => {
        console.error('Unable to Register Service Worker.', err);
      });
    },

    sendnotify() {
      if ( this.eta < 5.5 * 60 * 1000 ) {
        alert('No ETAs Found!');
      } else {
        alert('Notification Set!');
        console.log(this.eta);
        axios.post('./sendNotification', {delay: this.eta - 5 * 60 * 1000}, {headers: {'Content-Type': 'application/x-www-form-urlencoded'}})
        .then((res) => {
          console.log('sent' + res.data);
        })
        .catch((err) => {
          console.log(err);
        });
      }
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
