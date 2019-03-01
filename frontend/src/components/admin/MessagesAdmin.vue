<template>
    <div style="margin-top: 50px;" class="container">
        <p class="has-text-weight-bold">Current Message: "{{currentMessage}}" is {{ currentMessageVisible ? 'visible' : 'not visible'}}</p> 
        <div style="margin-top: 15px; margin-bottom: 15px;">
            <div class="field has-addons">
                <div class="control">
                    <input v-model="newMessage" class="input" type="text"/>
                    <p v-if="newMessage.length > 251" class="help is-danger">Message must be shorter than 251 characters.</p>
                </div>
                <!-- <div class="control">
                    <button class="button is-info">Set Message</button>
                </div> -->
            </div>
            <div class="field">
                <div class="control">
                    <label class="checkbox">
                        <input v-model="newMessageVisible" type="checkbox" />
                        Visible
                    </label>
                </div>
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
import AdminMessageUpdate from '@/structures/adminMessageUpdate';
import AdminServiceProvider from '@/structures/serviceproviders/admin.service';
import { setTimeout } from 'timers';

const sp = new InfoServiceProvider();

export default Vue.extend({
    data() {
        return {
            currentMessage: '',
            currentMessageVisible: false,
            newMessage: '',
            newMessageVisible: false,
            fail: false,
            success: false,
        }as {
            currentMessage: string;
            currentMessageVisible: boolean;
            newMessage: string;
            newMessageVisible: boolean;
            fail: boolean;
            success: boolean;
        };
    },
    mounted() {
        this.getCurrentMessage();
    },
    methods: {
        save() {
            const myMessage = new AdminMessageUpdate(this.newMessage, this.newMessageVisible, new Date(), new Date());
            AdminServiceProvider.SetMessage(myMessage).then((resp) => {
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
            sp.GrabAdminMessage().then((message: AdminMessageUpdate) => {
                this.currentMessage = message.message;
                this.currentMessageVisible = message.enabled;
                this.newMessage = message.message;
                this.newMessageVisible = message.enabled;
            });
        },
    },
});
</script>
