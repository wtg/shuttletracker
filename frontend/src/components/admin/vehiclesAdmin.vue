<template>
   <div style="margin-top: 50px;" class="container">
       <are-you-sure @no="shouldDelete = false;" @yes="vehicleDelete" :active="shouldDelete"/>
        <table class="table">
            <thead>
                <tr>
                <th><abbr title="Name">Name</abbr></th>
                <th><abbr title="ID">ID</abbr></th>
                <th><abbr title="Enabled">Enabled</abbr></th>
                <th></th>
                <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="vehicle in vehicles" :key="vehicle.id">
                    <th><router-link :to="'/admin/vehicles/' + String(vehicle.id)">{{vehicle.name}}</router-link></th>
                    <th><p>{{vehicle.id}}</p></th>
                    <th><p>{{vehicle.enabled}}</p></th>
                    <th><button @click="$router.push('/admin/vehicles/' + String(vehicle.id) + '/edit')" class="button">Edit</button></th>
                    <th><button @click="shouldDelete = true; vehicleToDelete = vehicle;" class="button is-danger">Delete</button></th>
                </tr>
                <tr>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th><button @click="$router.push('/admin/vehicles/-1/new')" class="button is-success">New</button></th>
                </tr>
                <div class="file">
                  <label class="file-label">
                    <input class="file-input" type="file" name="import" @change="handleFileUpload">
                    <span class="file-cta">
                      <span class="file-icon">
                        <img src="./../../assets/upload-icon.svg"></i>
                      </span>
                      <span class="file-label">
                        Choose a file…
                      </span>
                    </span>
                  </label>
                  <button v-on:click="addImportedVehicles()" class="button is-success">Submit</button>
                </div>
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import Vehicle from '../../structures/vehicle';
import AdminServiceProvider from '@/structures/serviceproviders/admin.service';
import AreYouSure from '@/components/admin/AreYouSure.vue';

const sp = new AdminServiceProvider();
export default Vue.extend({
    computed: {
        vehicles(): Vehicle[] {
            return this.$store.state.Vehicles;
        },
    },
    data() {
        return{
            shouldDelete: false,
            vehicleToDelete: undefined,
            file: null,
        } as {
            shouldDelete: boolean,
            vehicleToDelete: Vehicle | undefined,
            file: any,
        };
    },
    components: {
        AreYouSure,
    },
    mounted() {
        this.$store.dispatch('grabVehicles');
    },
    methods: {
        vehicleDelete() {
            this.shouldDelete = false;
            if (this.vehicleToDelete !== undefined) {
                AdminServiceProvider.DeleteVehicle((this.vehicleToDelete as Vehicle)).then(() => {
                    this.$store.dispatch('grabVehicles');
                });
            }
        },

        addImportedVehicles() {
          const reader = new FileReader();
          let json = null;
          // Capture file information
          reader.onload = ((theFile) => {
            return (e: any) => {
              try {
                json = JSON.parse(e.target.result);
                // If JSON.parse does not return a usable value, throw an Error
                if (json.length === undefined || json.length === 0) {
                  throw new Error('Improper JSON formatting');
                }
                for (let i = 0; i < json.length; i++) {
                    const obj = json[i];
                    // Create a new Vehicle object using the data in obj
                    const newVehicle = new Vehicle(obj.id, obj.name, obj.created, obj.updated, obj.enabled, obj.tracker_id);
                    // If creating the new Vehicle failed, the JSON file was not formatted correctly. Throw an Error
                    if (!newVehicle) {
                      throw new Error('Improper JSON formatting');
                    }
                    // Create the vehicle in the database
                    AdminServiceProvider.NewVehicle(newVehicle).then(() => {
                            this.$store.dispatch('grabVehicles');
                    });
                }
              } catch (e) {
                // If we get a SyntaxError, we have received a Non-JSON file. Replace the error to reflect this.
                if (e instanceof SyntaxError) {
                  e = new Error('Non-JSON file submitted');
                }
                // Display the Error to the user, and log it in the console
                alert(e);
                console.log(e);
              }
            };
          })(this.file);
          reader.readAsText(this.file);
        },

        // Gets the user-submitted file and stores it
        handleFileUpload(event: any) {
          this.file = event.target.files[0];
        },
    },
});
</script>
