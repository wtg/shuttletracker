<template>
    <div style="margin-top: 50px;" class="container">
        <table class="table">
            <thead>
                <tr>
                    <th><abbr title="ID">ID</abbr></th>
                    <th><abbr title="Message">Message</abbr></th>
                    <th><abbr title="Prompt">Prompt</abbr></th>
                    <th><abbr title="Created">Created</abbr></th>
                    <th><abbr title="Admin">Admin</abbr></th>
                    <th></th>
                    <th></th>
                </tr>
                <tr v-for="form in forms" :key="form.id">
                    <th>{{form.id}}</th>
                    <th>{{form.message}}</th>
                    <th>{{form.prompt}}</th>
                    <th>{{form.created}}</th>
                    <th>{{form.admin}}</th>
                    <th></th>
                    <th></th>
                </tr>
            </thead>
        </table>

        <p class="has-text-weight-bold">Current Message: "{{currentMessage}}"</p> 
        <div style="margin-top: 15px; margin-bottom: 15px;">
            <div class="field has-addons">
                <div class="control">
                    <input v-model="newMessage" class="input" type="text" placeholder="Message Contents"/>
                    <p v-if="newMessage.length > 251" class="help is-danger">Message must be shorter than 251 characters.</p>
                </div>
                <!-- <div class="control">
                    <button class="button is-info">Set Message</button>
                </div> -->
            </div>
            <div class="field">
                <div class="control">
                    <button @click="save" class="button is-success">Save</button>
                </div>
            </div>
        </div>
        <div v-if="fail" class="notification is-danger">
            <p>Failed to set message</p>
        </div>
        <div v-if="success" class="notification is-success">
            <p>Message set</p>
        </div>
    </div>  
</template>
<script lang="ts">
import Vue from 'vue';
import InfoServiceProvider from '@/structures/serviceproviders/info.service';
import AdminServiceProvider from '@/structures/serviceproviders/admin.service';
import FeedbackMessageUpdate from '@/structures/feedbackMessageUpdate';
import Form from '@/structures/form';
import { setTimeout } from 'timers';

const sp = new InfoServiceProvider();

export default Vue.extend({
    name: 'forms',
    computed: {
        forms(): Form[] {
            return this.$store.state.Forms;
        },
    },
    data() {
        return {
            currentMessage: '',
            newMessage: '',
            fail: false,
            success: false,
        }as {
            currentMessage: string;
            newMessage: string;
            fail: boolean;
            success: boolean;
        };
    },
    mounted() {
        this.getCurrentMessage();
    },
    methods: {
        save() {
            const myMessage = new Form(-1, this.newMessage, '', new Date(), true);
            AdminServiceProvider.CreateForm(myMessage).then((resp) => {
                if (resp.ok) {
                    this.success = true;
                    setTimeout(() => {
                        this.success = false;
                    }, 1000);
                } else {
                    this.fail = true;
                    setTimeout(() => {
                        this.fail = false;
                    }, 1000);
                }
                this.getCurrentMessage();
            }).catch(() => {
                this.fail = true;
                setTimeout(() => {
                    this.fail = false;
                }, 1000);
            });
        },
        getCurrentMessage() {
            sp.GrabFeedbackMessage().then((message: FeedbackMessageUpdate) => {
                this.currentMessage = message.message;
                this.newMessage = message.message;
            });
        },
    },
});
</script>
