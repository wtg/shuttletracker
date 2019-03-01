<template>
<div class="container" style="margin-top: 50px;">
    <div v-if="error" class="notification is-danger">
        <p>There was an error saving the vehicle.</p>
    </div>
    <div v-if="success" class="notification is-success">
        <p>Vehicle saved successfully</p>
    </div>
    <p v-if="!creation" class="is-size-2">Vehicle Editing</p>
    <p v-if="creation" class="is-size-2">Vehicle Creation</p>
    
    <!-- Text input-->
    <div class="field">
    <label class="label" for="vehicleID">Tracker ID</label>
    <div class="control">
        <input v-model.number="vehicle.tracker_id" id="vehicleID" name="vehicleID" type="text" placeholder="iTrak Vehicle ID" class="input ">        
    </div>
    </div>

    <!-- Text input-->
    <div class="field">
    <label class="label" for="vehicleName">Vehicle Name</label>
    <div class="control">
        <input v-model="vehicle.name" id="vehicleName" name="vehicleName" type="text" placeholder="Name" class="input ">
        <p v-if="vehicle.name.length === '0'" class="help is-danger">Required</p>
    </div>
    </div>

    <!-- Multiple Radios (inline) -->
    <div class="field">
    <label class="label" for="">Enabled/Disabled</label>
    <div class="control">
        <label class="radio inline" for="enabled-0">
        <input :checked="vehicle.enabled" @click="vehicle.enabled=true;" type="radio" name="enabled" id="enabled-0">
        Enabled
        </label>
        <label class="radio inline" for="enabled-1">
        <input :checked="!vehicle.enabled" @click="vehicle.enabled=false;" type="radio" name="enabled" id="enabled-1">
        Disabled
        </label>
    </div>
    </div>

    <!-- Button -->
    <div class="field">
    <label class="label" for=""></label>
    <div class="control">
        <button @click="send" id="" name="" class="button is-info" :class="{'is-loading': sending}">Submit</button>
    </div>
    </div>

</div>

</template>
<script lang="ts">
import Vue from 'vue';
import Vehicle from '../../structures/vehicle';
import AdminServiceProvider from '../../structures/serviceproviders/admin.service';

// This is the vehicle editing and creation interface.
export default Vue.extend({
    props: {
        creation: {
            type: Boolean,
        },
    },
    data() {
        return {
            vehicle: new Vehicle(0, '', new Date(), new Date(), false, -1),
            sending: false,
            error: false,
            success: false,
        } as {
            vehicle: Vehicle;
            sending: boolean;
            error: boolean;
            success: boolean;
        };
    },
    methods: {
        grabMyVehicle() {
            if (this.creation) {
                return;
            }
            for (let i = 0; i < this.$store.getters.getVehicles.length; i ++) {
                const tempRotue = this.$store.getters.getVehicles[i];
                const id = this.$route.params.id;
                if (Number(tempRotue.id) === Number(id)) {
                    this.vehicle.id = tempRotue.id;
                    this.vehicle.name = tempRotue.name;
                    this.vehicle.enabled = tempRotue.enabled;
                    this.vehicle.tracker_id = tempRotue.tracker_id;
                }
            }
        },
        send() {
            this.sending = true;
            if (this.creation) {
                AdminServiceProvider.NewVehicle(this.vehicle).then(() => {
                    this.success = true;
                    this.sending = false;
                    this.$store.dispatch('grabVehicles');
                    setTimeout(() => {
                        this.success = false;
                        this.$router.push('/admin/vehicles');
                    }, 2000);
                }).catch((err) => {
                    this.error = true;
                    this.sending = false;
                    setTimeout(() => {
                        this.error = false;
                    }, 2000);
                });
            } else {
                AdminServiceProvider.EditVehicle(this.vehicle).then(() => {
                    this.success = true;
                    this.sending = false;
                    this.$store.dispatch('grabVehicles');
                    setTimeout(() => {
                        this.success = false;
                    }, 2000);
                }).catch((err) => {
                    this.error = true;
                    this.sending = false;
                    setTimeout(() => {
                        this.error = false;
                    }, 2000);
                });
            }

        },
    },
    mounted() {
        this.$store.subscribe((mutation) => {
            if (mutation.type === 'setVehicles') {
                this.grabMyVehicle();
            }
        });
        if (this.$store.getters.getVehicles.length !== 0) {
            this.grabMyVehicle();

        }
    },

});
</script>
