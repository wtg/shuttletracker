<template>
    <div id="live" v-bind:style="liveStyle">
    <p>{{lv + " " + text}}</p>
    <div v-if="live" class="pulsate" style=""></div></div>
</template>
<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  data() {
    return {
      liveStyle: {},
      text: '',
      lv: 'Live',
      live: true,
    };
  },
  methods: {
},
});
</script>
<style lang="scss">

.pulsate {
  position:absolute;
  float:right;
  width:10px;
  height:10px;
  background-color:blue;
  border-radius:50%;
  top:9px;
  right:3px;
  animation: pulsate 2.5s ease-out;
  animation-iteration-count: infinite;
}
#live{
  color:black;
  width: auto;
  height:18px;
  padding:5px;
  border-radius:5px;
  font-size:15px;
  background-color:rgba(255, 255, 255, 0.88);
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.8);
  display:none;
  position: absolute;
  right:10px;
  top:42px;
}

#live p{
  float:left;
  padding: none;
  margin-top: 0px;
  white-space: nowrap;
}

@keyframes pulsate {
  0% {
    opacity: 0;

  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0;

  }
}


.pulsating-circle {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translateX(-50%) translateY(-50%);
  width: 30px;
  height: 30px;

  &:before {
    content: '';
    position: relative;
    display: block;
    width: 300%;
    height: 300%;
    box-sizing: border-box;
    margin-left: -100%;
    margin-top: -100%;
    border-radius: 45px;
    background-color: #01a4e9;
    animation: pulse-ring 1.25s cubic-bezier(0.215, 0.61, 0.355, 1) infinite;
  }

  &:after {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    display: block;
    width: 100%;
    height: 100%;
    background-color: white;
    border-radius: 15px;
    box-shadow: 0 0 8px rgba(0,0,0,.3);
    animation: pulse-dot 1.25s cubic-bezier(0.455, 0.03, 0.515, 0.955) -.4s infinite;
  }
}

@keyframes pulse-ring {
  0% {
    transform: scale(.33);
  }
  80%, 100% {
    opacity: 0;
  }
}

@keyframes pulse-dot {
  0% {
    transform: scale(.8);
  }
  50% {
    transform: scale(1);
  }
  100% {
    transform: scale(.8);
  }
}
</style>

