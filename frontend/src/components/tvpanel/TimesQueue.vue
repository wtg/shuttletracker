<template>
<div id ="main">
    <!-- East Queue -->
    <div id="east" v-if="curr_east">
        <ul>  
            <li id="type-east"> East  </li>
            <div class="times" >
            <li id="east1" class="time" v-if="east1_time" key="east1_time">{{display(east1_time)}}</li>
            <li id="east2" class="time" v-if="east2_time" key="east2_time">{{display(east2_time)}}</li>
            <li id="east3" class="time" v-if="east3_time" key="east3_time">{{display(east3_time)}}</li>
            <li id="east4" class="time" v-if="east4_time" key="east4_time">{{display(east4_time)}}</li>
            
            </div> 
        </ul>
    </div>
    <!-- West Queue -->
    <div id="west" v-if="curr_west">
        <ul>
            <li id="type-west"> West </li>
            <div class="times">
            <li id="west1" class="time" v-if="west1_time" key="west1_time">{{display(west1_time)}}</li>
            <li id="west2" class="time" v-if="west2_time" key="west2_time">{{display(west2_time)}}</li>
            <li id="west3" class="time" v-if="west3_time" key="west3_time">{{display(west3_time)}}</li>
            <li id="west4" class="time" v-if="west4_time" key="west4_time">{{display(west4_time)}}</li>
            </div>
        </ul>
    </div>
    <!-- Late Night/Weekend Queue -->
    <div id="weekendlate" v-if="curr_weekend_late">
        <ul>
            <li id="type" > Late Night </li>
            <div class="times">
            <li id="late1" class="time" v-show="this.curr_weekend_late">{{display(late1_time)}}</li>
            <li id="late2" class="time" v-show="this.curr_weekend_late">{{display(late2_time)}}</li>
            <li id="late3" class="time" v-show="this.curr_weekend_late">{{display(late3_time)}}</li>
            <li id="late4" class="time" v-show="this.curr_weekend_late">{{display(late4_time)}}</li>
            </div>
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
            curr_time: new Date(),    

            curr_west: weekdayW.scheduleTime,          
            curr_east: weekdayE.scheduleTime,                            
            curr_weekend_late: null,                   

            east1_time: null,
            east2_time: null,
            east3_time: null,
            east4_time: null,

            west1_time: null,
            west2_time: null,
            west3_time: null,
            west4_time: null,

            late1_time: null,
            late2_time: null,
            late3_time: null,
            late4_time: null,
            
            no_shuttles: "No Avaliable Shuttles",
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
                console.log("Hourly Queue check/update done");
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
                    this.curr_weekend_late = weekendlate.scheduleTime;
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
            // this.formatWeekendLate();
        },
        // // Function to manipulate weekendlate formatting
        // formatWeekendLate(){
        //     //If East and West shuttles are also running, format the queue
        //     if (this.curr_east && this.curr_west && this.curr_weekend_late){
        //         document.getElementById("weekendlate").style.marginTop = "250px";
        //     }
        //     // If only the Late Night shuttles are running, center the queue    
        //     else if (this.curr_weekend_late){
        //         document.getElementById("weekendlate").style.margin = "0 auto";
        //     }
        // },  
        // --------------------------------------------------------------------------
        // Function to handle updating shuttle times 
        updateShuttleTimes(){
            if (this.curr_east !== null){
                this.updateQueueEast();
            }
            if (this.curr_west !== null){
                this.updateQueueWest();
            }
            if (this.curr_weekend_late !== null){
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
            this.east1_time = this.curr_east[0];
            this.east2_time = this.curr_east[1];
            this.east3_time = this.curr_east[2];
            this.east4_time = this.curr_east[3];
            
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
            this.west1_time = this.curr_west[0];
            this.west2_time = this.curr_west[1];
            this.west3_time = this.curr_west[2];    
            this.west4_time = this.curr_west[3];
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
            this.late1_time = this.curr_weekend_late[0];
            this.late2_time = this.curr_weekend_late[1];
            this.late3_time = this.curr_weekend_late[2];  
            this.late4_time = this.curr_weekend_late[3];
        },
        // --------------------------------------------------------------------------
        // Function to convert 24 to 12 hour format and display AM or PM
        display(time){
            let temp = time.split(":");
            let hour = temp[0];
            let minutes = temp[1];
            if (parseInt(hour) < 12) {
                if (parseInt(hour) == 0){
                    hour = 12;
                }
                return hour + ':' + minutes + ' AM';
            }
            else if (parseInt(hour) == 12) {
                return hour + ':' + minutes + ' PM';
            }
            else {
                if (parseInt(hour) == 24) {
                    hour = hour - 12;
                    return hour + ':' + minutes + ' AM';
                }
                hour = hour - 12;
                return hour + ':' + minutes + ' PM';
            }
        },
        // --------------------------------------------------------------------------
    },
    mounted() {
        // Interval every 30 seconds; 30,000 milliseconds
        setInterval(() => {
            this.updateCurTime();
            this.checkHour();
        }, 30000);

        // Interval every three minutes; 180,000 milliseconds ****
        setInterval(() => {
            this.updateShuttleTimes();
        }, 18);
    },
});
</script>
<style lang="scss" scoped>
    #main{
        text-align:center;
        position:absolute;
        display:flex;
        width:100%;
        height:600px;
        color:black;
    }
    .time{
        font-size:22px;
        width:160px;
        height:65px;
        padding-top:13px;
        margin-left:25px;
        margin-top:20px;
        border: 1.5px solid #eee;
        border-radius: 4px;
        z-index:300;
        font-weight:400;
    }
    #east{
        flex:50%;
    }
    #west{
        flex:50%;
    }
    #weekendlate{ 
        flex:50%;
    }
    #type-east{
        font-size:60px; 
        color:#96C03A;
    }
    #type-west{
        font-size:60px;
        color:#E1501B;
    }
    .fade-enter-active, .fade-leave-active {
     transition: opacity .5s;
    }
    .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
    opacity: 0;
    }
</style>