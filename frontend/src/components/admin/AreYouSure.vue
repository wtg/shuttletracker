<template>
    <div class="modal" :class="{'is-active': active}">
        <div class="modal-background"></div>
        <div class="modal-content">
            <div class="box">
                <p class="title">Are you sure?</p>
                <p class="subtitle">This action is permanent</p>
                <div class="level">
                    <div class="level-left"></div>
                    <div class="level-right">
                        <div class="control">
                            <button @click="$emit('yes');" style="margin-right: 10px;" class="button is-success">Yes</button>
                        </div>
                        <div class="control">
                            <button @click="$emit('no');" class="button">No</button>
                        </div>
                    </div>
                </div>
                
            </div>
        </div>
        <button @click="$emit('no');" class="modal-close is-large" aria-label="close"></button>
    </div>
</template>

<script lang="ts">
// This is an 'are you sure' modal, it is active on the active property, and emits no when close or no is clicked, and emits yes when yes is clicked
import Vue from 'vue';
export default Vue.extend({
    props: {
        active: {
            type: Boolean,
            default: () => false,
        },
    },

});
</script>
