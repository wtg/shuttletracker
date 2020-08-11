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
    transition: all .8s ease;
    touch-action: manipulation;
    border: 1px solid var(--color-button-bg);
    border-radius: 50%;
    box-shadow: 0 6px 10px 0 rgba(0, 0, 0, 0.14),
              0 1px 18px 0 rgba(0, 0, 0, 0.12),
              0 3px 5px -1px rgba(0, 0, 0, 0.20);
    -webkit-box-shadow: 0 6px 10px 0 rgba(0, 0, 0, 0.14),
              0 1px 18px 0 rgba(0, 0, 0, 0.12),
              0 3px 5px -1px rgba(0, 0, 0, 0.20);
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

    &:hover {
        transition: all .8s ease;
        background-color: var(--color-button-hover);
    }

    &:active {
        transition: all .8s ease;
        background-color: var(--color-button-hover);
        box-shadow: 0px 12px 17px 2px hsla(0,0%,0%,0.14), 
            0px 5px 22px 4px hsla(0,0%,0%,0.12), 
            0px 7px 8px -4px hsla(0,0%,0%,0.2);
        -webkit-box-shadow: 0px 12px 17px 2px hsla(0,0%,0%,0.14), 
            0px 5px 22px 4px hsla(0,0%,0%,0.12), 
            0px 7px 8px -4px hsla(0,0%,0%,0.2);
    
    }
}


</style>

