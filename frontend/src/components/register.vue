<template>
    <div class="container">
        <router-link to="/">Back</router-link>
        <h1>Notification Register</h1> 
        <h4>User Phone Number:</h4>
        <input v-model="phone_number" placeholder="ex. 1234567890">
        <!-- phone_number holds string -->
        <h4>Carrier Type</h4>
        <select v-model="carrier">
            <option disabled value="">Select a Carrier</option>
            <option>Verizon</option>
            <option>AT&T</option>
            <option>T-mobile</option>
            <option>Ooredoo</option>
        </select>
        <!-- carrier holds value -->
        <h4>Shuttle Stop</h4>
        <select v-model="stop">
            <option disabled value="">Select a Stop</option>
            <option>Student Union</option>
            <option>B</option>
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
        <tabulator />
        <br>
        <!-- TODO disable button until flag -->
        <button v-on:click="" v-bind:disabled="route == ''">Submit</button>
    </div>
</template>
<script lang="ts">

import Vue from 'vue';
import tabulator from './tabulator.vue'
export default Vue.extend({
    components: {
        tabulator
    },
    data: {
        phone_number: '',
        carrier: '',
        route: '',
        stop: '',
    },
    methods: {
        submit(){
            //TODO store info into db
        }
    },
});
</script>
<style lang="scss">
    .container{
        margin: 20px;
    }
</style>
