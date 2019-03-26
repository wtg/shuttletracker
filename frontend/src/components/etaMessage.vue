<template>
    <transition name="pop">
        <div class="eta-message" v-if="message && show">
            {{ message }}
        </div>
    </transition>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  props: ['etaInfo', 'show'],
  computed: {
    message(): string | null {
      if (this.etaInfo === null) {
          return null;
      }
      const now = new Date();

      let newMessage = `${this.etaInfo.route.name} shuttle arriving at ${this.etaInfo.stop.name}`;
      // more than 1 min 30 sec?
      if (this.etaInfo.eta.eta.getTime() - now.getTime() > 1.5 * 60 * 1000 && !this.etaInfo.eta.arriving) {
        newMessage += ` in ${relativeTime(now, this.etaInfo.eta.eta)}`;
      }
      newMessage += '.';

      return newMessage;
    },
  },
});

function relativeTime(from: Date, to: Date): string {
  const minuteMs = 60 * 1000;
  const elapsed = to.getTime() - from.getTime();

  // cap display at thirty min
  if (elapsed < minuteMs * 30) {
    return `${Math.round(elapsed / minuteMs)} minutes`;
  }

  return 'a while';
}
</script>

<style lang="scss" scoped>
.eta-message {
    width: 320px;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    margin: 60px 20px auto auto;
    z-index: 1000;
    background: white;
    padding: 20px 28px;
    border: 0.5px solid #eee;
    border-radius: 4px;
    box-shadow: 0 1px 16px -4px #bbb;
    font-size: 17px;
}
@media screen and (max-width: 500px) {
    .eta-message {
        width: auto;
        margin: 50px 10px 0 10px;
        padding: 16px 22px;
        font-size: 16px;
    }
}
.pop-enter-active {
    animation: pop-in 0.4s;
}
.pop-leave-active {
    animation: pop-out 0.5s;
}
@keyframes pop-in {
    0% {
        transform: scale(0);
    }
    60% {
        transform: scale(1.05);
    }
    100% {
        transform: scale(1);
    }
}
@keyframes pop-out {
    0% {
        transform: scale(1);
    }
    100% {
        transform: scale(0);
    }
}
</style>
 