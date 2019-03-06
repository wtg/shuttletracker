<template>
    <div id="map">
        <h1>MAP</h1>
    </div>
</template>

<script lang="ts">
// This component handles the Map for the TV Panel
import Vue from 'vue';
export default Vue.extend({
    name: 'Map',
});
</script>

<style scoped>
    h1 {
        font-size:200px;
    }

    #map{
        width:48%;
        margin-right:25px;
        background-color:green;
        float:right;
        height:700px;
        position:relative;
        bottom:80px;
        text-align:center;
    }
</style>
