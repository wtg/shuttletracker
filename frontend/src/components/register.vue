<template>
    <div class="container">
        <router-link to="/">Back</router-link>
        <h1>Notification Register</h1>
        <h4>Shuttle Stop : <span v-html="stop_id"></span></h4>
        <h4>User Phone Number :</h4>
        <input v-model.trim="phone_number" placeholder="ex. 1234567890">
        <!-- phone_number holds string -->
        <h4>Carrier Type</h4>
        <select v-model="carrier">
            <option selected disabled>Select a Carrier</option>
            <option value="verizon">Verizon</option>
            <option value="at&t">AT&T</option>
            <option value="t-mobile">T-mobile</option>
            <option value="ooredoo">Ooredoo</option>
        </select>
        <!-- carrier holds value -->
        <h4>Shuttle Route :</h4>
        <input type="radio" v-model="route" value="east">
        <label for="east"> East</label>
        <br>
        <input type="radio" v-model="route" value="west">
        <label for="west"> West</label>
        <br>
        <input type="radio" v-model="route" value="wln">
        <label for="wln"> Weekend Late Night</label>
        <br>
        <!-- route holds value -->
        <br><br>
        <Tabulator/>
        <br>
        <Selected/>
        <br>
        <!-- TODO make time reactive -->
        <!-- <span>Times are : {{ time }}</span> -->
        <br>
        <!-- TODO disable button until flag -->
        <button v-on:click="submit()" v-bind:disabled="(route == '')||(carrier == '')">Submit</button>
        <br><br><br>
    </div>
</template>
<script lang="ts">

import Vue from 'vue';
import Tabulator from './Tabulator.vue'
import Selected from './Selected.vue'
import { time } from './Tabulator.vue'
// console.log('why is this'+time);
export default Vue.extend({
    components: {
        Tabulator,
        Selected,
    },
    data() {
        return {
            phone_number: '',
            carrier: '',
            carriers: ['verizon', 'at&t'],
            route: '',
            stop_url:  decodeURI(window.location.href),
        }
    },
    computed: {
        stop_id: function() {
            return ((this.stop_url.split('?')[1]).split('=')[1]);
        },
    },
    methods: {
        submit() {
            //test phone_num
            if ( this.phone_number.length !== 10 || Number.isNaN(this.phone_number as any) ) {
                //error phone
                console.log("Error : Phone Number Invalid")
                return
            }
            //test carrier
            if (  this.carriers.indexOf(this.carrier) <= -1 ) {
                //error carrier
                console.log("Error : Carrier Invalid")
                return
            }
            //test route
            if ( this.route == '' ) {
                //error route
                console.log("Error : Route Not Selected")
                return
            }
            //test time
            if ( !time ) {
                //error time
                console.log("Error : Time Not Selected")
                return
            } 
            
            //TODO submit 


        },
    },
});
</script>
<style lang="scss">
    .container{
        margin: 20px;
    }
</style>
