<template>
    <div class="container">
        <router-link to="/">Back</router-link>
        <h1>Notification Register</h1>
        <h4>User Phone Number:</h4>
        <input v-model="phone_number" placeholder="ex. 1234567890">
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
        <h4>Shuttle Stop</h4>
        <select v-model="stop">
            <option disabled selected>Select a Stop</option>
            <option value="union">Student Union</option>
            <option value="b">B</option>
        </select>
        <h4>Shuttle Route</h4>
        <input type="checkbox" id="east" value="east" v-model="route">
        <label for="east">East</label>
        <br>
        <input type="checkbox" id="west" value="west" v-model="route">
        <label for="west">West</label>
        <br>
        <input type="checkbox" id="wln" value="wln" v-model="route">
        <label for="wln">Weekend Late Night</label>
        <!-- route holds value -->
        <br><br>
        <Tabulator/>
        <br>
        <Selected/>
        <br>
        <!-- TODO make time reactive -->
        <span>Times are : {{ time }}</span>
        <br>
        <!-- TODO disable button until flag -->
        <button v-on:click="" v-bind:disabled="route == ''">Submit</button>
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
    data: {
        phone_number: '',
        carrier: '',
        route: '',
        stop: '',
    },
    methods: {
        submit() {
            // TODO store info into db
        },
    },
});
</script>
<style lang="scss">
    .container{
        margin: 20px;
    }
</style>
