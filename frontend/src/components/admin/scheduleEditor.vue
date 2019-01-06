<template>
    <div>
        <div class="field">
            <div class="control">
                <button @click="addAndUpdate" class="button">Add</button>
            </div>
        </div>
        <span class="box" v-for="(interval, idx) in value" :key="idx">
        <p class="has-text-weight-semibold">Start:</p>
        <div class="field has-addons">
            <div class="control">
                <span class="select">
                    <select v-model.number="interval.start_day">
                    <option value=0>Sunday</option>
                    <option value=1>Monday</option>
                    <option value=2>Tuesday</option>
                    <option value=3>Wednesday</option>
                    <option value=4>Thursday</option>
                    <option value=5>Friday</option>
                    <option value=6>Satruday</option>
                    </select>
                </span>
            </div>
            <div class="control">
                <timepicker v-model="interval.start_time"/>
            </div>
        </div>
        <p class="has-text-weight-semibold">End:</p>
        <div class="field has-addons">
            <div class="control">
                <span class="select">
                    <select v-model.number="interval.end_day">
                    <option value=0>Sunday</option>
                    <option value=1>Monday</option>
                    <option value=2>Tuesday</option>
                    <option value=3>Wednesday</option>
                    <option value=4>Thursday</option>
                    <option value=5>Friday</option>
                    <option value=6>Satruday</option>
                    </select>
                </span>
            </div>
            <div class="control">
                <timepicker v-model="interval.end_time"/>
            </div>
        </div>
        <div class="field">
            <div class="control">
                <button @click="remove(idx);" class="button is-danger">Remove</button>
            </div>
        </div>
        </span>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import RotueScheduleInterval from '../../structures/routeScheduleInterval';
import timepicker from '@/components/admin/timepicker.vue';

// Takes a schedule interval as v-model and allows editing
export default Vue.extend({
    props: {
        value: {
            type: Array as () => RotueScheduleInterval[],
            default: () => [],
        },
    },
    methods: {
        addAndUpdate() {
            const tmp: RotueScheduleInterval[] = this.value;
            tmp.push(new RotueScheduleInterval(0, 0, 0, new Date(), 0, new Date()));
            this.$emit('input', tmp);
        },
        remove(idx: number) {
            const tmp: RotueScheduleInterval[] = this.value;
            tmp.splice(idx, 1);
            this.$emit('input', tmp);
        },
        dateAsTimeString(date: Date): string {
            return String(date.getHours() < 10 ? '0' + String(date.getHours()) : String(date.getHours())) + ':' + String(date.getMinutes() < 10 ? '0' + String(date.getMinutes()) : String(date.getMinutes()));
        },
    },
    components: {
        timepicker,
    },
});
</script>
