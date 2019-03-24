<template>
<div id ="main">
    <!-- East Queue -->
    <div id="east" v-if="checkEast()">
        <ul>  
            <li id="type"> EAST  </li>
            <li id="east1" class="time" v-show="this.curr_east">1</li>
            <li id="east2" class="time" v-show="this.curr_east">2</li>
            <li id="east3" class="time" v-show="this.curr_east">3</li>
        </ul>
    </div>
    <!-- West Queue -->
    <div id="west" v-if="checkWest()">
        <ul>
            <li id="type"> WEST </li>
            <li id="west1" class="time" v-show="this.curr_west">1</li>
            <li id="west2" class="time" v-show="this.curr_west">2</li>
            <li id="west3" class="time" v-show="this.curr_west">3</li>
        </ul>
    </div>
    <!-- Late Night/Weekend Queue -->
    <div id="weekendlate" v-if="checkLate()">
        <ul>
            <li id="type" > LATE NIGHT </li>
            <li id="late1" class="time" v-show="this.curr_weekend_late">1</li>
            <li id="late2" class="time" v-show="this.curr_weekend_late">2</li>
            <li id="late3" class="time" v-show="this.curr_weekend_late">3</li>
        </ul>
    </div>
</div>
</template>

<script lang="ts">
// This component handles the Shuttle Times Queue on the TV Panel
import Vue from 'vue';

// Importing East Campus shuttle times (JSON)
import weekdayE from '@/assets/shuttle_times/weekdayE.json';
import weekendE from '@/assets/shuttle_times/weekendE.json';

// Importing Weekend/Late shuttle times (JSON)
import weekendlate from '@/assets/shuttle_times/weekendlate.json';

// Importing West Campus shuttle times (JSON)
import weekdayW from '@/assets/shuttle_times/weekdayW.json';
import satW from '@/assets/shuttle_times/satW.json';
import sunW from '@/assets/shuttle_times/sunW.json';

export default Vue.extend({
    name: 'TimesQueue',
    data(){
        return {
            curr_time: new Date(),        // Current Time (Date Object)
            curr_west: undefined,         // Current West Time Queue (Array)
            curr_east: undefined,         // Current East Time Queue (Array)
            curr_weekend_late: undefined, // Current Late/Weekend Time Queue (Array)
        }
    },
    methods: {
        // --------------------------------------------------------------------------
        // Display/Update the current time 
        updateCurTime(){
            this.curr_time = new Date();
        },
        // --------------------------------------------------------------------------
        // Function to handle checking/updating the queues every hour 
        checkHour(){
            if (this.curr_time.getMinutes() === 0){
                this.updateQueues();
                this.checkEast();
                this.checkWest();
                this.checkLate();
                console.log("Hourly check/update done");
            }
        },
        // Function to update/set the component's queue values
        updateQueues(){
            // Monday - Thursday 
            if (this.curr_time.getDay() >= 1 && this.curr_time.getDay() <= 4){
                // Morning (Update at 6AM)
                if (this.curr_time.getHours() === 6){
                    this.curr_weekend_late = null || undefined;
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
                    this.curr_weekend_late = null || undefined;
                }
                // Late night (Update at 7PM)
                else if (this.curr_time.getHours() === 19 ){
                    this.curr_east = null || undefined;
                    this.curr_west = null || undefined;
                    this.curr_weekend_late = weekendlate.scheduleTime;
                }
            }
            // Sunday
            else {
                // Morning (Update at 9AM)
                if (this.curr_time.getHours() === 9){
                    this.curr_east = weekendE.scheduleTime;
                    this.curr_west = sunW.scheduleTime;
                    this.curr_weekend_late = null || undefined;
                }
            }
            // Format weekendlate queue
            this.formatWeekendLate();
        },
        // Toggle East Shuttles Queue
        checkEast(){
            if (this.curr_east != null){
                return true;
            }
            return false;
        },
        // Toggle West Shuttles Queue
        checkWest(){
            if (this.curr_west != null){
                return true;
            }
            return false;
        },
        // Toggle Late/Weekend Shuttles Queue
        checkLate(){
            if (this.curr_weekend_late != null){
                return true;
            }
            return false;    
        },
        // Function to manipulate weekendlate formatting
        formatWeekendLate(){
            //If East and West shuttles are also running, format the queue
            if (this.curr_east && this.curr_west && this.curr_weekend_late){
                document.getElementById("weekendlate").style.marginTop = "250px";
            }
            // If only the Late Night shuttles are running, center the queue    
            else if (this.curr_weekend_late){
                document.getElementById("weekendlate").style.margin = "0 auto";
            }
        },
        // --------------------------------------------------------------------------
        // Function to handle updating shuttle times 
        updateShuttleTimes(){
            if (this.curr_east != null){
                this.updateQueueEast();
            }
            if (this.curr_west != null){
                this.updateQueueWest();
            }
            if (this.curr_weekend_late != null){
                this.updateQueueLate();
            }
            console.log("Shuttle times updated");
        },
        // Function to update shuttle times for the East Queue
        updateQueueEast(){  
            let now = this.curr_time;
            let first_shuttle_time = this.curr_east[0].split(":");
            if (now.getHours() > parseInt(first_shuttle_time[0])){
                this.curr_east.shift();
            }            
            else if (now.getHours() === parseInt(first_shuttle_time[0])){
                if (now.getMinutes() > parseInt(first_shuttle_time[1]) + 2){
                    this.curr_east.shift();
                }
            }
            // Display the first three shuttle times of the queue
            document.getElementById('east1').innerHTML = (this.curr_east[0]) ? this.display(this.curr_east[0]) : ""; 
            document.getElementById('east2').innerHTML = (this.curr_east[1]) ? this.display(this.curr_east[1]) : "";
            document.getElementById('east3').innerHTML = (this.curr_east[2]) ? this.display(this.curr_east[2]) : "";
        },
        // Function to update shuttles times for the West Queue
        updateQueueWest(){
            let now = this.curr_time;
            let first_shuttle_time = this.curr_west[0].split(":");
            if (now.getHours() > parseInt(first_shuttle_time[0])){
                this.curr_west.shift();
            }            
            else if (now.getHours() === parseInt(first_shuttle_time[0])){
                if (now.getMinutes() > parseInt(first_shuttle_time[1]) + 2){
                    this.curr_west.shift();
                }
            }
            // Display the first three shuttle times of the queue
            document.getElementById('west1').innerHTML = (this.curr_west[0]) ? this.display(this.curr_west[0]) : ""; 
            document.getElementById('west2').innerHTML = (this.curr_west[1]) ? this.display(this.curr_west[1]) : "";
            document.getElementById('west3').innerHTML = (this.curr_west[2]) ? this.display(this.curr_west[2]) : "";
        },
        // Function to update shuttle times for the Late/Weekend Queue
        updateQueueLate(){
            let now = this.curr_time;
            let first_shuttle_time = this.curr_weekend_late[0].split(":");
            if (now.getHours() > parseInt(first_shuttle_time[0])){
                this.curr_weekend_late.shift();
            }            
            else if (now.getHours() === parseInt(first_shuttle_time[0])){
                if (now.getMinutes() > parseInt(first_shuttle_time[1]) + 2){
                    this.curr_weekend_late.shift();
                }
            }
            // Display the first three shuttle times of the queue
            document.getElementById('late1').innerHTML = (this.curr_weekend_late[0]) ? this.display(this.curr_weekend_late[0]) : ""; 
            document.getElementById('late2').innerHTML = (this.curr_weekend_late[1]) ? this.display(this.curr_weekend_late[1]) : "";
            document.getElementById('late3').innerHTML = (this.curr_weekend_late[2]) ? this.display(this.curr_weekend_late[2]) : "";
        },
        // --------------------------------------------------------------------------
        // Function to convert 24 to 12 hour format and display AM or PM
        display(time){
            let temp = time.split(":");
            let hour = temp[0];
            let minutes = temp[1];
            if (parseInt(hour) < 12){
                if (parseInt(hour) == 0){
                    hour = 12;
                }
                return hour + ':' + minutes + ' AM';
            }
            else if (parseInt(hour) == 12){
                return hour + ':' + minutes + ' PM';
            }
            else {
                if (parseInt(hour) == 24){
                    hour = hour - 12;
                    return hour + ':' + minutes + ' AM';
                }
                hour = hour - 12;
                return hour + ':' + minutes + ' PM';
            }
        }
        // --------------------------------------------------------------------------
    },

    mounted() {
        // Interval every minute; 60,000 milliseconds
        setInterval(() => {
            this.updateCurTime();
            this.checkHour();
        }, 60000);

        // Interval every three minutes; 180,000 milliseconds
        setInterval(() => {
            this.updateShuttleTimes();
        }, 180000);
    },
});
</script>

<style scoped>
    #main{
        margin-left:25px;
        float:left;
        height:700px;
        width:49%;
        top:159px;
        position:relative; 
        text-align:center;
        color:black;
    }
    .time{
        font-size:35px;
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
        font-size:80px;
    }
    
    .fade-enter-active, .fade-leave-active {
        transition: opacity 0.25s ease-out;
    }

    .fade-enter, .fade-leave-to {
        opacity: 0;
    }
</style>