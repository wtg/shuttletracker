<template>
    <div>
        <h1>Hello asdf</h1>
        <button @click="toggleList">Toggle</button>
        <list v-if="listVisible"/>
        <another v-if="!listVisible"/>
    </div>
</template>
<script lang="ts">
import Vue from 'vue'
import list from '@/components/list.vue';
import another from '@/components/another.vue';

export default Vue.extend({
    data(){
        return {
            listVisible: false,
        } as {
            listVisible: Boolean;
        }
    },
    methods: {
        toggleList(){
            this.listVisible = !this.listVisible;
        }
    },
    components: {
        list,
        another
    }
})
</script>
