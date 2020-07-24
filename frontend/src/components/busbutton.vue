<template>
    <b-button @click="$emit('bus-click');" class="bus-button">
<span> {{this.$store.state.settings.busButtonChoice}} </span>
    </b-button>
</template>
<script lang="ts">
import Vue from 'vue';
export default Vue.extend({

});
</script>

<style lang="scss" scoped>
@import "@/assets/vars.scss";

.bus-button{
    touch-action: manipulation;
    // border: 2px solid var(--color-primary);
    border-radius: 50%;
    box-shadow: 0.5px 0.5px 3px 1px rgba(var(--color-box-shadow), 0.3);
    -webkit-box-shadow: 0.5px 0.5px 3px 1px rgba(var(--color-box-shadow), 0.3);
    width: 50px;
    margin: 0;
    padding: 0;
    display: flex;
    background-color: var(--color-button-bg);
    justify-content: center;
    align-items: center;
    height: 50px;
    font-size: 20px;
    text-decoration: none;
    cursor: pointer;
    outline: 0;

    &:hover{
        background-color: var(--color-button-hover);
    }
}


</style>

