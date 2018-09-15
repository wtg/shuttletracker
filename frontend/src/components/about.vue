<template>
    <div class="container">
        <router-link to="/">Back</router-link>
        <h2>About the Shuttle Tracker</h2> 
        <p>
        The Shuttle Tracker is <a href="https://github.com/wtg/shuttletracker">open source</a> and maintained by the <a href="https://webtech.union.rpi.edu">Web Technologies Group</a> of the Rensselaer Union Student Senate for the benefit of the student body.
        It is also an active <a href="https://rcos.io">RCOS</a> project, and dozens of students have contributed to it through RCOS in previous semesters.
        </p> 
        <h3>Why is there no app?</h3> 
        <p>
            We're focused on making the Shuttle Tracker website as good as it can be. While most of our users are on mobile, we think the best way to deliver a great experience to the most people is to focus our efforts on this website.
        </p> 
        <h3>How can I contribute?</h3> 
        <p>
            We have an active <a href="https://github.com/wtg/shuttletracker">repository</a> on GitHub where we track issues and accept pull requests. If you want to get involved in another way, email <a href="mailto:webtech@union.lists.rpi.edu">webtech@union.lists.rpi.edu</a> for more information.
        </p>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
export default Vue.extend({

});
</script>
<style lang="scss">
    .container{
        margin: 20px;
    }
</style>
