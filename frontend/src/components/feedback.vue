<template>
    <div class = "parent content">
        <h2>Feedback</h2>
        <div v-if="this.adminMessage !== undefined" >
            <div style="width: 100%;float:left;" v-html="this.adminMessage"></div>
        </div>
        <!-- <p>
            <br>We'd love to hear back from you!
            <br>If you have any comments, suggestions, or concerns, submit them down below and they'll be sent directly to us.
        </p> -->
        <div>
            <td style="padding-top: 5px;">
            <textarea 
                v-model="feedbackMessage"
                placeholder="beep beep I'm a sheep . . ." 
                rows = "8" cols = "127">
            </textarea>
            <br>
            <button @click="save" class="submit" type="submit" form="feedback" value="Submit"
            id="submit">Submit</button>
            </td>
        </div>
        <div v-if="fail" class="notification is-danger">
            <p>Failed to submit form.</p>
        </div>
        <div v-if="success" class="notification is-success">
            <p>Submission successful! Thanks for sending us your feedback!</p>
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
    data() {
        return {
            feedbackMessage: '',
            adminMessage: '',
            fail: false,
            success: false,
        }as {
            feedbackMessage: string;
            adminMessage: string;
            fail: boolean;
            success: boolean;
        };
    },
    mounted() {
        this.getAdminMessage();
    },
    methods: {
        save() {
            const myMessage = new Form(-1, this.feedbackMessage, new Date(), false);
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
            }).catch(() => {
                this.fail = true;
                setTimeout(() => {
                    this.fail = false;
                }, 1000);
            });
        },
        getAdminMessage() {
            sp.GrabFeedbackMessage().then((message: FeedbackMessageUpdate) => {
                this.adminMessage = message.message;
            });
        },
    },
});
</script>

<style lang="scss" scoped>
.parent{
    padding: 50px;
    p {
        margin: .10em;
    }
}
.submit{
    &:hover {
    color: white;
    background-color: lighten(red, 5%);
    cursor: pointer;
    cursor: hand;
    }
    height:40px;
    width:100px;
    border: none;
    border-radius: 8px;
    margin-top: .5em;
    transition: all 0.35s;
}

h2 {
    margin-bottom:.1em;
}
</style>