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
            <input type="radio" v-model="route" name="route" value="2">
            <label for="2"> East</label>
            <br>
            <input type="radio" v-model="route" name="route" value="1">
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
            carriers: ['AT&T','T-Mobile','Verizon','Sprint','XFinity Mobile','Virgin Mobile','Tracfone','Metro PCS','Boost Mobile','Cricket','Republic Wireless','Google Fi','U.S. Cellular','Ting','Consumer Cellular','C-Spire','Page Plus'],
            err: 0,
            route: '',
            stop_url:  decodeURI(window.location.href),
            times: [],
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
            this.route = this.stop_url.split('?')[1].split('&')[2].split('=')[1];
            return this.route;
        }
    },
    methods: {

        submit() {
            //test phone_num
            if ( this.phone_number.length !== 10 || Number.isNaN(this.phone_number as any) ) {
                //error phone
                console.log("Error : Phone Number Invalid");
                this.error(1);
                return
            }
            //test carrier
            if ( this.carriers.indexOf(this.carrier) <= -1 ) {
                //error carrier
                console.log("Error : Carrier Invalid");
                this.error(2);
                return
            }
            //test stop
            if ( Number.isNaN(this.stop_id as any) ) {
                console.log("Error : Invalid Stop Id");
                this.error(3);
                return
            }
            //test route
            if ( this.route.length != 1 ) {
                //error route
                console.log("Error : Route Not Selected");
                this.error(4);
                return
            }
            //test time
            if ( this.times.length < 1 ) {
                //error time
                console.log("Error : Time Not Selected");
                this.error(5);
                return
            } 
            
            //TODO submit 


        },

        error( err: any ) {
            if ( err == 1 ) {
                alert("Error : Phone Number Invalid.")
            } else if ( err == 2 ) {
                alert("Error : Carrier Invalid.")
            } else if ( err == 3 ) {
                alert("Error : Invalid Stop Id.")
            } else if ( err == 4 ) {
                alert("Error : Route Not Selected.")
            } else if ( err == 5 ) {
                alert("Error : Time Not Selected.")
            } else {
                alert("Error : Something isn't Working!")
            }
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
