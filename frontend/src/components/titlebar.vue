<template>
    <div class="titleBar" style="filter: invert(0)">
        <!-- left side of tile bar -->
        <ul class="titleContent" id="titleContent-left">
            <li>
                <dropdown @displaymodal="showModal=true" />
            </li>
            <li>
                <p class="title">RPI Shuttle Tracker</p>
            </li>
        </ul>
        <!-- right side of title bar -->
        <ul class="titleContent" id="titleContent-right">
            <li>
                <a href="https://webtech.union.rpi.edu" class="logo">
                    <img src="../assets/images/wtg.svg">
                </a>
            </li>
            <li>
                <!-- <live-indicator></live-indicator> -->
            </li>
        </ul>
        <about-modal v-if="showModal" @hidemodal="showModal=false" />
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    import dropdown from './menu.vue';
    export default Vue.extend({
    components: {
        dropdown,
    },
   });
</script>
