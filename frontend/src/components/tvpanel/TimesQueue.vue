<template>
    <div id="app">
        <div id ="title"> QUEUE </div>
        <ul>
            <li id="time"> Shuttle Time 1 </li>
            <li id="time"> Shuttle Time 2 </li>
            <li id="time"> Shuttle Time 3 </li>
            <li id="time"> Shuttle Time 4 </li>
            <li id="time"> Shuttle Time 5 </li>
        </ul>
    </div>
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
    #title{
        font-size:80px;
    }
    #app{
        margin-left:25px;
        float:left;
        height:700px;
        width:48%;
        bottom:80px;
        position:relative;
        text-align:center;
        background-color:red;
    }
    #time{
        font-size:40px;
    }
    ul {
        padding-left:0;
    }
</style>
