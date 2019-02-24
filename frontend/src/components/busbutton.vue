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
    border: 1px solid black;
    border-radius: 50%;
    
    width: 50px;
    margin: 0;
    padding: 0;
    display: flex;
    background-color: rgba(255,255,255,.8);
    justify-content: center;
    align-items: center;
    height: 50px;
    font-size: 20px;
    text-decoration: none;
    cursor: pointer;
}
</style>

