
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
                <div class="container">
                  <div class="large-12 medium-12 small-12 cell">
                    <label>File
                      <input type="file" id="inputfile" ref="inputfile" @change="handleFileUpload"/>
                    </label>
                      <button v-on:click="addImportedRoutes()">Submit</button>
                  </div>
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
              try {
                json = JSON.parse(e.target.result);
                console.log('json global var has been set to parsed json of this file here it is unevaled = \n' + JSON.stringify(json));
                if (json.length === undefined || json.length === 0) {
                  throw new Error('Improper JSON formatting');
                }
                for (let i = 0; i < json.length; i++) {
                    const obj = json[i];
                    console.log(obj);
                    const newRoute = new Route(-1, obj.name, obj.description, obj.enabled, obj.color, obj.width, obj.points,
                      obj.schedule, obj.active, obj.stop_ids);
                    if (!newRoute) {
                      throw new Error('Improper JSON formatting');
                    }
                    AdminServiceProvider.CreateRoute(newRoute).then(() => {
                            this.$store.dispatch('grabRoutes');
                    });
                }
              } catch (e) {
                if (e instanceof SyntaxError) {
                  e = new Error('Non-JSON file submitted');
                }
                alert(e);
                console.log(e);
              }
            };
          })(this.file);
          reader.readAsText(this.file);

          // }

        },

        handleFileUpload(event: any) {
          // this.file = (document.getElementById('inputfile') as HTMLInputElement);
          this.file = event.target.files[0];
          // make sure the file is recieved
          console.log(this.file);
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
