<template>
<div id ="queue">
    <h3> The Current Time is {{displayCurTime()}} </h3>
    <h3> The Day is {{curr_time.getDay()}} </h3>

    <!-- East Queue -->
    <div id="east" v-if="checkEast()">
        <ul>  
            <li id="type"> EAST  </li>
            <li id="time"> 1 </li>
            <li id="time"> 2 </li>
            <li id="time"> 3 </li>
        </ul>
    </div>

    <!-- West Queue -->
    <div id="west" v-if="checkWest()">
        <ul>
            <li id="type"> WEST </li>
            <li id="time"> 1 </li>
            <li id="time"> 2 </li>
            <li id="time"> 3 </li>
        </ul>
    </div>

    <!-- Late Night/Weekend Queue -->
    <div id="weekendlate" v-if="checkLate()">
        <ul>
            <li id="type" > LATE NIGHT </li>
            <li id="time"> 1 </li>
            <li id="time"> 2 </li>
            <li id="time"> 3 </li>
        </ul>
    </div>
</div>
</template>

<script lang="ts">
// This component handles the Shuttle Time Queue on the TV Panel
import Vue from 'vue';

// Importing East Campus shuttle times (JSON)
import weekdayE from '@/assets/shuttle_times/weekdayE.json';
import weekendE from '@/assets/shuttle_times/weekendE.json';

// Importing Weekend/Late Night shuttle times (JSON)
import weekendlate from '@/assets/shuttle_times/weekendlate.json';

// Importing West Campus shuttle times (JSON)
import weekdayW from '@/assets/shuttle_times/weekdayW.json';
import satW from '@/assets/shuttle_times/satW.json';
import sunW from '@/assets/shuttle_times/sunW.json';

export default Vue.extend({
    name: 'TimesQueue',
    data(){
        return {
            curr_time: new Date(),
            curr_west: undefined,
            curr_east: undefined, 
            curr_weekend_late: undefined,
        }
    },
    methods: {

        // Display/update the current time 
        displayCurTime(){
            this.curr_time = new Date();
            return this.curr_time.getHours() + ': ' + this.curr_time.getMinutes();
        },

        // --------------------------------------------------------------------------
        // Check if EAST shuttles are running
        checkEast(){
            if (this.curr_east){
                return true;
            }
            return false;
        },
        // Check if WEST shuttles are running 
        checkWest(){
            if (this.curr_west){
                return true;
            }
            return false;
        },
        // Check if LATENIGHT/WEEKEND shuttles are running 
        checkLate(){
            if (this.curr_weekend_late){
                // If east and west shuttles are also running, format 
                if (this.curr_east && this.curr_west){
                    document.getElementById("weekendlate").style.marginTop = "300px";
                }   
                else {
                    document.getElementById("weekendlate").style.margin = "0 auto";
                }
                return true; 
            }
            return false; 
        },
        // --------------------------------------------------------------------------

        // Update queues
        updateQueues(){
            // Monday - Thursday 
            if (this.curr_time.getDay() >= 1 && this.curr_time.getDay() <= 4){
                // Morning (Update at 6AM)
                if (this.curr_time.getHours() === 6){
                    delete this.curr_weekend_late;
                    this.curr_west = weekdayW.scheduleTime;
                    this.curr_east = weekdayE.scheduleTime;
                } 
            }
            // Friday 
            else if (this.curr_time.getDay() === 5){
                // Morning (Update at 6AM)
                if (this.curr_time.getHours() === 6){
                    this.curr_west = weekdayW.scheduleTime;
                    this.curr_east = weekdayE.scheduleTime;
                }
                // Late night (Update at 7PM)
                else if (this.curr_time.getHours() === 19){
                    this.weekendlate = weekendlate.scheduleTime;
                }
            }
            // Saturday
            else if (this.curr_time.getDay() === 6) {
                // Morning (Update at 9AM)
                if (this.curr_time.getHours() === 9 ){
                    this.curr_east = weekendE.scheduleTime;
                    this.curr_west = satW.scheduleTime;
                    delete this.curr_weekend_late;
                }
                // Late night (Update at 7PM)
                else if (this.curr_time.getHours() === 19 ){
                    delete this.curr_east;
                    delete this.curr_west;
                    this.curr_weekend_late = weekendlate.scheduleTime;
                }
            }
            // Sunday
            else {
                // Morning (Update at 9AM)
                if (this.curr_time.getHours() === 9){
                    this.curr_east = weekendE.scheduleTime;
                    this.curr_west = sunW.scheduleTime;
                    delete this.curr_weekend_late;
                }
            }
        },


    },

    mounted() {
        // Interval to update current time, check to display EAST,WEST,LATE queues
        setInterval(() => {
            this.displayCurTime();
        }, 60000);


        // Interval to update queues every hour 
        setInterval(() => {
            this.updateQueues();
            this.checkEast();
            this.checkWest();
            this.checkLate();
        }, 100);

    },
});
</script>

// CSS Styling for the Shuttle Time Queue
<style scoped>
    #queue{
        margin-left:25px;
        float:left;
        height:700px;
        width:49%;
        top:159px;
        position:relative;
        text-align:center;
        color:black;
    }
    #time{
        font-size:20px;
    }
    ul {
        padding-left:0;
    }
    #east{
        float:left;
        margin-left:200px;
        padding-top:30px;
    }
    #west{
        float:right;
        margin-right:200px;
        padding-top:30px;
    }
    #weekendlate{
        margin:auto 0;
        padding-top:30px;
    }
    #type{
        font-size:60px;
    }
</style>