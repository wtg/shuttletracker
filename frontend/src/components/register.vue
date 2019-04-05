<template>
    <div class="container">
        <router-link to="/">Back</router-link>
        <h1 class="title">Notification Register</h1>
        <p class="subtitle">Shuttle Stop : <span v-html="stop_name"></span></p>
        <p class="subtitle">User Phone Number :
        <input v-model.trim="phone_number" placeholder="ex. 1234567890">
        </p>
        <!-- phone_number holds string -->
        <p class="subtitle">Carrier Type :
        <select v-model="carrier">
            <option selected disabled>Select a Carrier</option>
            <option v-for="c in carriers" :value="c">{{c}}</option>
        </select></p>
        <!-- carrier holds value -->
        <p class="subtitle">Shuttle Route :
            <br>
            <input type="radio" v-model="route" value="2">
            <label for="2"> East</label>
            <br>
            <input type="radio" v-model="route" value="1">
            <label for="1"> West</label>
            <br>
        </p>
        <!-- route holds value -->
        <p class="subtitle">Select Times
            <Tabulator/>
        </p>
        <Selected/>
        <br>
        <!-- TODO make time reactive -->
        <!-- <span>Times are : {{ time }}</span> -->
        <!-- TODO disable button until flag -->
        <button v-on:click="submit()">Submit</button>
        <br><br><br>
    </div>
</template>
<script lang="ts">

import Selected from './Selected.vue'
import Tabulator from './Tabulator.vue'
import Vue from 'vue';
import EventBus from '../event_bus'
export default Vue.extend({
    components: {
        Selected,
        Tabulator,
    },

    data() {
        const url = decodeURI(window.location.href);
        return {
            phone_number: '',
            carrier: '',
            carriers: ['AT&T','T-Mobile','Verizon','Sprint','XFinity Mobile','Virgin Mobile','Tracfone','Metro PCS','Boost Mobile','Cricket','Republic Wireless','Google Fi','U.S. Cellular','Ting','Consumer Cellular','C-Spire','Page Plus'],
            route: url.split('?')[1].split('&')[2].split('=')[1],
            stop_url: url,
            times: [] as string[],
        }
    },

    computed: {
        stop_id: function() {
            return this.stop_url.split('?')[1].split('&')[0].split('=')[1];
        },
        stop_name: function() {
            return this.stop_url.split('?')[1].split('&')[1].split('=')[1];
        },
    },

    methods: {
        addData (payload:any) {
            let time = payload.day + ' ' + payload.time;
            this.times.push(time);
        },
        removeData (payload:any) {
            let time = payload.day + ' ' + payload.time;
            let index = this.times.indexOf(time);
            this.times[index] = this.times[0];
            this.times.shift();
        },
        submit() {
            //test phone_num
            if ( this.phone_number.length !== 10 || Number.isNaN(this.phone_number as any) ) {
                //error phone
                console.log("Error : Phone Number Invalid");
                alert("Error : Phone Number Invalid.");
                return
            }
            //test carrier
            if ( this.carriers.indexOf(this.carrier) <= -1 ) {
                //error carrier
                console.log("Error : Carrier Invalid");
                alert("Error : Carrier Invalid.")
                return
            }
            //test stop
            if ( Number.isNaN(this.stop_id as any) ) {
                console.log("Error : Invalid Stop Id");
                alert("Error : Invalid Stop Id.");
                return
            }
            //test route
            if ( this.route.length != 1 ) {
                //error route
                console.log("Error : Route Not Selected");
                alert("Error : Route Not Selected.");
                return
            }
            //test time
            if ( this.times.length < 1 ) {
                //error time
                console.log("Error : Time Not Selected");
                alert("Error : Time Not Selected.");
                return
            } 
            
            //TODO submit
        },
    },

    mounted() {
        EventBus.$on('TIME_ADDED', (payload:any) => {
            this.addData(payload)
        });
        EventBus.$on('TIME_REMOVED', (payload:any) => {
            this.removeData(payload)
        });
    },
});

</script>
<style lang="scss">
.parent {
  padding: 20px;
}
.container{
    margin: 20px;
}
.button {
    background-color: blue;
    border: none;
    color: white;
    padding: 15px 32px;
    text-align: center;
    display: inline-block;
    font-size: 16px; 
}
</style>
