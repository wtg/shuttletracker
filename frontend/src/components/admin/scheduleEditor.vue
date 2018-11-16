<template>
    <div>
        <div class="field">
            <div class="control">
                <button @click="addAndUpdate" class="button">Add</button>
            </div>
        </div>
        <span class="box" v-for="(interval, idx) in value" :key="idx">
        <p class="has-text-weight-semibold">Start:{{dateAsTimeString(interval.startTime)}}</p>
        <div class="field has-addons">
            <div class="control">
                <span class="select">
                    <select v-model.number="interval.endDate">
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
                <input @input="interval.startTime =  $event.target.valueAsDate" :value="dateAsTimeString(interval.startTime)" class="input" type="time"/>
            </div>
        </div>
        <p class="has-text-weight-semibold">End:</p>
        <div class="field has-addons">
            <div class="control">
                <span class="select">
                    <select v-model.number="interval.endDate">
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
                <input @input="interval.endTime =  $event.target.valueAsDate" :value="dateAsTimeString(interval.endTime)" class="input" type="time"/>
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

// Takes a schedule interval as v-model and allows editing
export default Vue.extend({
    props: {
        value: {
            type: Array as () => RotueScheduleInterval[],
            default: () => [],
        }
    },
    methods: {
        addAndUpdate(){
            let tmp: RotueScheduleInterval[] = this.value;
            tmp.push(new RotueScheduleInterval(0,0,0,new Date(),0,new Date()));
            this.$emit('input', tmp);
        },
        remove(idx: number){
            let tmp: RotueScheduleInterval[] = this.value;
            tmp.splice(idx,1);
            this.$emit('input',tmp);
        },
        dateAsTimeString(date: Date): string{
            return String(date.getHours()) + ':' + String(date.getMinutes());
        },

    }
})
</script>
