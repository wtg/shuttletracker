<template>
    <div id="map">
        <map-view id="map-view"></map-view>
    </div>
</template>

<script lang="ts">
// This component handles the Map for the TV Panel
import Vue from 'vue';
import mapView from '../admin/map.vue';
export default Vue.extend({
    name: 'Map',
    components:{
        mapView,
    },
});
</script>

<style scoped>
    h1 {
        font-size:200px;
    }
    #map-view{
        height:100%;
        width:100%;
    }
    #map{
        width:48%;
        margin-right:25px;
        float:right;
        height:600px;
        position:relative;
        top:200px;
        text-align:center;
    }
</style>
