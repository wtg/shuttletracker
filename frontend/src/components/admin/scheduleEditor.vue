<template>
    <div>
        <span v-for="interval in RotueScheduleInterval" :key="interval.id">
        <p class="has-text-weight-semibold">Start:</p>
        <div class="field has-addons">
            <div class="control">
                <span class="select">
                    <select>
                    <option>Select dropdown</option>
                    <option>With options</option>
                    </select>
                </span>
            </div>
            <div class="control">
                <input class="input" type="time"/>
            </div>
        </div>
        <p class="has-text-weight-semibold">End:</p>
        <div class="field has-addons">
            <div class="control">
                <span class="select">
                    <select>
                    <option>Select dropdown</option>
                    <option>With options</option>
                    </select>
                </span>
            </div>
            <div class="control">
                <input class="input" type="time"/>
            </div>
        </div>
        </span>
        <button class="button">Add</button>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import RotueScheduleInterval from '../../structures/routeScheduleInterval';

export default Vue.extend({
    props: {
        value: {
            type: Array as () => RotueScheduleInterval[],
            default: () => [],
        }
    }
})
</script>
