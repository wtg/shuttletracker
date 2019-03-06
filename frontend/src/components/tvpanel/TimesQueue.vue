<template>
    <div id="queue" v-if="$store.state.online">
        <div id="east">
            <ul>  
                <li id="type"> EAST  </li>
                <li id="time"> Shuttle Time 1 </li>
                <li id="time"> Shuttle Time 2 </li>
                <li id="time"> Shuttle Time 3 </li>
            </ul>
        </div>
        <div id="west">
            <ul>
                <li id="type"> WEST </li>
                <li id="time"> Shuttle Time 1 </li>
                <li id="time"> Shuttle Time 2 </li>
                <li id="time"> Shuttle Time 3 </li>
            </ul>
        </div>
    </div>
    <div v-else> Shuttles Services are currently down </div>
</template>

<script lang="ts">
// This component handles the Shuttle Time Queue on the TV Panel
import Vue from 'vue';
export default Vue.extend({
    name: 'TimesQueue',
});
</script>

// CSS Styling for the Shuttle Time Queue
<style scoped>
    #queue{
        margin-left:25px;
        float:left;
        height:700px;
        width:48%;
        top:159px;
        position:relative;
        text-align:center;
        color:black;
    }
    #time{
        font-size:40px;
    }
    ul {
        padding-left:0;
    }
    #east{
        float:left;
        padding-left:130px;
        padding-top:30px;
    }
    #west{
        float:right;
        padding-right:130px;
        padding-top:30px;
    }
    #type{
        font-size:80px;
    }
</style>
