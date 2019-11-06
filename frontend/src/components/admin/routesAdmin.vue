<template>
    <div style="margin-top: 50px;" class="container">
        <are-you-sure @no="shouldDelete = false;" @yes="deleteRoute" :active="shouldDelete"/>
        <table class="table">
            <thead>
                <tr>
                <th><abbr title="Name">Name</abbr></th>
                <th><abbr title="ID">ID</abbr></th>
                <th><abbr title="Description">Desc.</abbr></th>
                <th><abbr title="Schedule">Sched.</abbr></th>
                <th><abbr title="Enabled">Enabled</abbr></th>
                <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="route in routes" :key="route.id">
                <th><router-link :to='"/admin/routes/" + String(route.id) + "/"' >{{route.name}}</router-link></th>
                <td>{{route.id}}</td>
                <td>{{route.description}}</td>
                <td><ul>
                        <li v-for="interval in route.schedule" :key="interval.id">{{String(interval)}}</li>
                    </ul></td>
                <td>{{route.enabled}}</td>
                <td><button class="button" @click="$router.push('/admin/routes/' + String(route.id) + '/edit');">Edit</button></td>
                <td><button class="button is-danger" @click="shouldDelete = true; routeToDelete = route;">Delete</button></td>

                </tr>
                <tr>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th><button @click="$router.push('/admin/routes/-1/new')" class="button is-success">New</button></th>
                </tr>
                <tr>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th></th>
                  <th><button @click="downloadRoutes();" class="button is-success">Export</button></th>
                </tr>
                <br>
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
                  <button v-on:click="addImportedRoutes()" class="button is-success">Submit</button>
                </div>
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import Route from '../../structures/route';
import AdminServiceProvider from '../../structures/serviceproviders/admin.service';
import AreYouSure from '@/components/admin/AreYouSure.vue';
import routeEditing from '../../components/admin/routeEditing.vue';

export default Vue.extend({
    name: 'routes',
    data() {
        return{
            shouldDelete: false,
            routeToDelete: undefined,
            file: null,
        } as {
            shouldDelete: boolean,
            routeToDelete: Route | undefined,
            file: any,
        };
    },
    computed: {
        routes(): Route[] {
            return this.$store.state.Routes;
        },
    },
    methods: {
        deleteRoute() {
            this.shouldDelete = false;
            if (this.routeToDelete !== undefined) {
                AdminServiceProvider.DeleteRoute(this.routeToDelete).then(() => {
                    this.$store.dispatch('grabRoutes');
                });
            }
        },


        addImportedRoutes() {
          // assume routes can be accessed from json
          const reader = new FileReader();
          let json = null;
          // Closure to capture the file information.
          reader.onload = ((theFile) => {
            return (e: any) => {
              // Wrap route creation in a try-catch block
              try {
                json = JSON.parse(e.target.result);
                // If JSON.parse does not return a usable value, throw an Error
                if (json.length === undefined || json.length === 0) {
                  throw new Error('Improper JSON formatting');
                }
                for (let i = 0; i < json.length; i++) {
                    const obj = json[i];
                    // Create a new Route object using the data stored in 'obj'
                    const newRoute = new Route(-1, obj.name, obj.description, obj.enabled, obj.color, obj.width, obj.points,
                      obj.schedule, obj.active, obj.stop_ids);
                    // If creating the new Route failed, the JSON file was not formatted correctly. Throw an Error
                    if (!newRoute) {
                      throw new Error('Improper JSON formatting');
                    }
                    // Create the route in the database
                    AdminServiceProvider.CreateRoute(newRoute).then(() => {
                            this.$store.dispatch('grabRoutes');
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

        // Gets the user-submitted file and stores it in 'file'
        handleFileUpload(event: any) {
          this.file = event.target.files[0];
        },

        downloadRoutes() {
          const link = document.createElement('a');
          link.href = '/routes';
          link.setAttribute('download', '/routes') ; // or any other extension
          document.body.appendChild(link);
          link.click();
        },

    },
    components: {
        AreYouSure,
    },

});
</script>
