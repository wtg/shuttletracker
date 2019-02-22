<template>
    <ul>
        <li v-for="element in elements" :key="element">{{element}}</li>
    </ul>
</template>
<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
    data(){
        return {
            elements: ["asdf","asdf1","asdf2"],

        }
    }
})
</script>
