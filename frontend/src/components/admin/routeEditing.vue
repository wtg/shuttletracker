<template>
<div class="container">
    <div class="columns">
        <div class="column is-8">
            <map-view></map-view>
        </div>
        <form class="form-horizontal column" >
            <!-- Text input-->
            <div class="field">
            <label class="label" for="Name">Name</label>
            <div class="control">
                <input id="Name" name="Name" type="text" placeholder="Name" class="input " required="">
                <p class="help is-danger">Please Enter a name</p>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="Description">Description</label>
            <div class="control">
                <input id="Description" name="Description" type="text" placeholder="Description" class="input " required="">
                <p class="help is-danger">Please enter a description</p>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="color">Color</label>
            <div class="control">
                <input id="color" name="color" type="text" placeholder="#ff00ff" class="input " required="">
                <p class="help is-danger">Please set a color</p>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="width">Width</label>
            <div class="control">
                <input id="width" name="width" type="text" placeholder="4" class="input " required="">
                <p class="help is-danger">Please set a width</p>
            </div>
            </div>

            <!-- Multiple Radios -->
            <div class="field">
            <label class="label" for="enabled">Enabled/Disabled</label>
            <div class="control">
                <label class="radio" for="enabled-0">
                <input type="radio" name="enabled" id="enabled-0" value="enabled" checked="checked">
                enabled
                </label>
                <label class="radio" for="enabled-1">
                <input type="radio" name="enabled" id="enabled-1" value="disabled">
                disabled
                </label>
            </div>
            </div>
        </form>
    </div>
</div>
</template>
<script lang="ts">
import Vue from 'vue';
import mapView from '@/components/admin/map.vue';

export default Vue.extend({
    components: {
        mapView,
    },
});
</script>
