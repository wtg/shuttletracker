<template>
    <div style="margin-top: 50px;" class="container">
        <table class="table">
            <thead>
                <tr>
                    <th><abbr title="Name">Name</abbr></th>
                    <th><abbr title="ID">ID</abbr></th>
                    <th><abbr title="Desc">Description</abbr></th>
                    <th><abbr title="Lat.">Latitude</abbr></th>
                    <th><abbr title="Lng.">Longitude</abbr></th>
                    <!-- <th><abbr title="Routes">Routes</abbr></th> -->
                    <th><abbr title="Enabled">Enabled</abbr></th>
                    <th></th>
                    <th></th>
                </tr>
                <tr v-for="stop in stops" :key="stop.id">
                    <th>{{stop.name}}</th>
                    <th>{{stop.id}}</th>
                    <th>{{stop.description}}</th>
                    <th>{{stop.latitude.toFixed(3)}}</th>
                    <th>{{stop.longitude.toFixed(3)}}</th>
                    <!-- <th>{{stop.routesOn}}</th> -->
                    <th></th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th><button @click="$router.push('/admin/stops/-1/new')" class="button is-success">New</button></th>
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
                  <button v-on:click="addImportedStops()" class="button is-success">Submit</button>
                </div>
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import { Stop } from '@/structures/stop';
import AdminServiceProvider from '@/structures/serviceproviders/admin.service';
export default Vue.extend({
    name: 'stops',

    computed: {
        stops(): Stop[] {
            return this.$store.state.Stops;
        },
    },
    data() {
        return{
            file: null,
        } as {
            file: any,
        };
    },
    mounted() {
        this.$store.dispatch('grabStops');
    },
    methods: {
        addImportedStops() {
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
                    // Create a new Stop object using the data in obj
                    const newStop = new Stop(obj.id, obj.name, obj.description, obj.latitude, obj.longitude, obj.created, obj.updated);
                    // If creating the new Stop failed, the JSON file was not formatted correctly. Throw an Error
                    if (!newStop) {
                      throw new Error('Improper JSON formatting');
                    }
                    // Create the stop in the database
                    AdminServiceProvider.NewStop(newStop).then(() => {
                            this.$store.dispatch('grabStops');
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
