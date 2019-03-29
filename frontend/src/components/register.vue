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
            <option value="verizon">Verizon</option>
            <option value="at&t">AT&T</option>
            <option value="t-mobile">T-mobile</option>
            <option value="ooredoo">Ooredoo</option>
        </select></p>
        <!-- carrier holds value -->
        <p class="subtitle">Shuttle Route :
            <br>
            <input type="radio" v-model="stop_route" value="2">
            <label for="2"> East</label>
            <br>
            <input type="radio" v-model="stop_route" value="1">
            <label for="1"> West</label>
            <br>
        </p>
        <!-- route holds value -->
        <p class="subtitle">Select Times</p>
        <Tabulator/>
        <Selected/>
        <br>
        <!-- TODO make time reactive -->
        <!-- <span>Times are : {{ time }}</span> -->
        <!-- TODO disable button until flag -->
        <button v-on:click="submit()" v-bind:disabled="(route == '')||(carrier == '')">Submit</button>
        <br><br><br>
    </div>
</template>
<script lang="ts">

import Selected from './Selected.vue'
import Tabulator from './Tabulator.vue'
import Vue from 'vue';
// console.log('why is this'+time);
export default Vue.extend({
    components: {
        Selected,
        Tabulator,
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
            return this.stop_url.split('?')[1].split('&')[0].split('=')[1];
        },
        stop_name: function() {
            return this.stop_url.split('?')[1].split('&')[1].split('=')[1];
        },
        stop_route: function() {
            return this.stop_url.split('?')[1].split('&')[2].split('=')[1];
        }
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
            /*if ( !time ) {
                //error time
                console.log("Error : Time Not Selected")
                return
            } */
            
            //TODO submit 


        },
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
</style>
