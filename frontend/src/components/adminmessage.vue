<template>
<span>
    <div v-if="(this.message !== undefined && this.message.enabled === true) && !hide" id="messagebox">
        <div style="width: 100%;float:left;" v-html="this.message.message"></div>
        <div @click="hide = !hide" style="cursor: pointer;position:absolute;right:10px;top:6px;color:#333;font-size:20px;">&times;</div>
    </div>
</span>
</template>
<script lang="ts">
// This component handles rendering of a single administrator message

import Vue from 'vue';
import AdminMessageUpdate from '@/structures/adminMessageUpdate';
export default Vue.extend({
    data() {
        return {
            hide: false,
        };
    },
    computed: {
        message(): AdminMessageUpdate {
            return this.$store.state.adminMessage;
        },
    },
});
</script>

<style lang="scss" scoped>
#messagebox {
  max-height: 150px;
  overflow-y: auto;
  position: relative;
  z-index: 0;
  bottom: 0px;
  left: 0;
  right: 0;
  padding: 10px 30px 10px 15px; /* leave some space for the close button */
  background: #ee2222;
  color: #ffffff;
  text-align: center;
  font-size: 15px;
}
</style>
