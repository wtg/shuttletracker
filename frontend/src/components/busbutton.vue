<template>
    <button @click="$emit('bus-click');" class="bus-button">
        <span> 🚐</span>
    </button>
</template>
<script lang="ts">
import Vue from 'vue';
export default Vue.extend({

});
</script>

<style lang="scss" scoped>

.bus-button{
    border-radius: 50%;
    width: 50px;
    margin: 0;
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 50px;
    font-size: 20px;
    cursor: pointer;
}
</style>

