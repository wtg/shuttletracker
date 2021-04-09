<template>
   <div style="margin-top: 50px;" class="container">
       <are-you-sure @no="shouldDelete = false;" @yes="FeedbackDelete" :active="shouldDelete"/>
        <table class="table">
            <thead>
                <tr>
                <th><abbr title="Number">Number</abbr></th>
                <th><abbr title="Feedback">Feedback</abbr></th>
                <th><abbr title="Comment">Comment</abbr></th>
                <th></th>
                <th></th>
                </tr>
                <tr v-for="form in forms" :key="form.id">
                    <th>{{form.id}}</th>
                    <th>{{form.message}}</th>
                    <th>{{form.created}}</th>
                    <th>{{form.admin}}</th>
                    <th></th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <tr>
                  <td></td>
                  <td></td>
                  <td></td>
                </tr>
                <tr>
                  <td></td>
                  <td></td>
                  <td></td>
                </tr>
                <tr>
                  <td></td>
                  <td></td>
                  <td></td>
                 </tr>
                <tr>
                  <td></td>
                  <td></td>
                  <td></td>
                </tr>
                <tr>
                  <td></td>
                  <td></td>
                  <td></td>
                </tr>
                <tr>
                  <td></td>
                  <td></td>
                  <td></td>
                    <button @click="save" class="submit" type="submit" form="feedback" value="Submit"
                id="submit"><b>Submit</b></button>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import { Stop } from '@/structures/form';
export default Vue.extend({
    name: 'form',
    computed: {
        form(): form[] {
            return this.$form.state;
        },
    },
});
</script>