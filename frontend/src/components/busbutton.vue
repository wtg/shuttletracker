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
    border: 1px solid #ccc;
    border-radius: 50%;
    box-shadow: 0 4px 8px 0 #ddd;
    -webkit-box-shadow: 0 4px 8px 0 #ddd;
    width: 50px;
    margin: 0;
    padding: 0;
    display: flex;
    background-color: rgba(250,250,250,.8);
    justify-content: center;
    align-items: center;
    height: 50px;
    font-size: 20px;
    text-decoration: none;
    cursor: pointer;

    &:hover{
        background-color: #fff;
    }
}


</style>

