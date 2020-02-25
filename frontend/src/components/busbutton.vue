<template>
    <button @click="$emit('bus-click');" class="bus-button">
<span> {{this.$store.state.settings.busButtonChoice}} </span>
    </button>
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
    border: 1px solid var(--color-fg-light);
    border-radius: 50%;
    box-shadow: 0 4px 8px 0 var(--color-bg-least);
    -webkit-box-shadow: 0 4px 8px 0 var(--color-bg-least);
    width: 50px;
    margin: 0;
    padding: 0;
    display: flex;
    background-color: rgba(var(--color-bg-less-rgb), .8);
    justify-content: center;
    align-items: center;
    height: 50px;
    font-size: 20px;
    text-decoration: none;
    cursor: pointer;
    outline: 0;

    &:hover{
        background-color: var(--color-bg-normal);
    }
}


</style>

