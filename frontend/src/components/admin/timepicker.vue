<template>
    <input @input="input($event.target.value);" :value="timestampString" class="input" type="time"/>
    
</template>
<script lang="ts">
import Vue from 'vue';

// timepicker is a component which takes a javascript Date in v-model, and allows the user to change the time via an html time picker.
export default Vue.extend({
    props: {
        value: {
            type: Date as () => any, // TODO: leaving this as any right now because typechecking doesn't want to play nice with the date type for some reason
        },
    },
    data() {
        return {
            hours: (this.value as Date).getHours(),
            minutes: (this.value as Date).getMinutes(),
        } as {
            hours: number,
            minutes: number,
        };
    },
    watch: {
        value() {
            this.hours = (this.value as Date).getHours();
            this.minutes = (this.value as Date).getMinutes();
        },
    },
    methods: {
        input(value: string) {
            const arr = value.split(':');
            const ret = new Date();
            ret.setHours(Number(arr[0]));
            ret.setMinutes(Number(arr[1]));
            this.$emit('input', ret);
        },
    },
    computed: {
        timestampString(): string {
            return String(this.hours < 10 ? '0' + String(this.hours) : String(this.hours)) + ':' + String(this.minutes < 10 ? '0' + String(this.minutes) : String(this.minutes));
        },
    },

});
</script>
