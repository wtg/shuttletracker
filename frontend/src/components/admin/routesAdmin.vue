<template>
    <div style="margin-top: 50px;" class="container">
        {{routes}}
    </div>
</template>
<script lang="ts">
import Vue from 'vue'
import Route from '../../structures/route';
export default Vue.extend({
    name: "routes",
    computed: {
        routes(): Route[]{
            return this.$store.state.Routes;
        }
    }
    
})
</script>
