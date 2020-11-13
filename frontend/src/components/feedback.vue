<template>
    <div class="parent">
        <h1 class="title">Feedback</h1>
        <h2 class="subtitle">All responses are greatly appreciated!<hr></h2>
        <div id="fail-n" class="hidden notification is-danger">
            <p>Failed to submit form.</p>
        </div>
        <div id="success-n" class="hidden notification is-success">
            <p>Feedback sent!</p>
        </div>
        <div id="body">
        <div id="feedback-area">
            <div v-if="this.adminMessage !== undefined" >
                <b><div id="admin-msg" style="width: 100%;float:left;" v-html="this.adminMessage"></div></b>
            </div>
            <textarea
                id="form"
                v-model="feedbackMessage"
                placeholder="Type your response here..." 
                rows="12"
                cols="80">
            </textarea>
            <br>
            <button @click="save" class="submit" type="submit" form="feedback" value="Submit"
            id="submit"><b>Submit</b></button>
        </div>
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
import {DarkTheme} from '@/structures/theme';

// const tinycolor = require('tinycolor2');

// const darkColor = tinycolor(line.options.color);
// darkColor.darken(15);

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
            let success = document.getElementById("success-n");
            let fail = document.getElementById("fail-n");
            AdminServiceProvider.CreateForm(myMessage).then((resp) => {
                if (resp.ok) {
                    success.classList.remove("hidden")
                    success.classList.add("visible")
                    setTimeout(() => {
                        success.classList.remove("visible")
                        success.classList.add("hidden")
                    },2000);
                } else {
                    fail.classList.add("visible")
                    setTimeout(() => {
                        fail.classList.remove("visible")
                        fail.classList.add("hidden")
                    },2000);
                }
            }).catch(() => {
                fail.classList.add("visible")
                    setTimeout(() => {
                        fail.classList.remove("visible")
                        fail.classList.add("hidden")
                    },2000);
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
@import "@/assets/vars.scss";

/*
1.955em - desktop question font
1.2em - mobile phone question font
0.5em - between the q and the tbox - mobile
6 cols textarea on mobile
padding between submit and bottom row - mobile
Replace submit button with feedback sent
Add dark mode support
*/

.parent{
    padding: 20px;
}
#admin-msg{
    font-size: 1.955em;
}
#feedback-area{
    display: block;
    margin-left: auto;
    margin-right: auto;
    width: 50%;
}
#form{ 
    border-radius: 8px;
    outline: none !important;
    border:2px solid rgb(202, 202, 202);
    margin-top: 1em;
    margin-bottom: 0.25em;
    width:  100%;
    height: 100%;
    padding: 5px;
    resize: none;
}
.submit{
     &:hover {
        cursor: pointer;
        transition: all .8s ease;
        background-color: var(--color-primary);
        box-shadow: 0px 12px 17px 2px hsla(0,0%,0%,0.14), 
            0px 5px 22px 4px hsla(0,0%,0%,0.12), 
            0px 7px 8px -4px hsla(0,0%,0%,0.2);
        -webkit-box-shadow: 0px 12px 17px 2px hsla(0,0%,0%,0.14), 
            0px 5px 22px 4px hsla(0,0%,0%,0.12), 
            0px 7px 8px -4px hsla(0,0%,0%,0.2);
    }

    &:active {
        transition: all .8s ease;
        background-color: var(--color-primary);
        box-shadow: 0px 12px 17px 2px hsla(0,0%,0%,0.14), 
            0px 5px 22px 4px hsla(0,0%,0%,0.12), 
            0px 7px 8px -4px hsla(0,0%,0%,0.2);
        -webkit-box-shadow: 0px 12px 17px 2px hsla(0,0%,0%,0.14), 
            0px 5px 22px 4px hsla(0,0%,0%,0.12), 
            0px 7px 8px -4px hsla(0,0%,0%,0.2);
    }
    box-shadow: 0 6px 10px 0 rgba(0, 0, 0, 0.14),
              0 1px 18px 0 rgba(0, 0, 0, 0.12),
              0 3px 5px -1px rgba(0, 0, 0, 0.20);
    -webkit-box-shadow: 0 6px 10px 0 rgba(0, 0, 0, 0.14),
              0 1px 18px 0 rgba(0, 0, 0, 0.12),
              0 3px 5px -1px rgba(0, 0, 0, 0.20);
    color: white;
    background-color: var(--color-primary);
    height:40px;
    width:100px;
    border: none;
    border-radius: 8px;
    margin-top: .5em;
    transition: all 0.35s;
    margin-left: calc(50% - 50px);
}
.notification{
    top: 2px;
    right: 2px;
    position: fixed;
}
.visible{
    visibility: visible;
    opacity: 1;
    transition: opacity 0s;
}
.hidden{
    visibility: hidden;
    opacity: 0;
    transition: visibility 0s 1s, opacity 1s ease-out;
}
@media only screen and (max-width: 1024px) {
    #feedback-area{ 
        margin-bottom: 0.5em;
        width: 95%;
    }
    #form{
        padding-bottom: 6em;
    }
}
@media only screen and (max-width: 400px) {
    #feedback-area{
        margin-bottom: 0.5em;
        width: 95%;
    }   
}
@media only screen and (max-height: 700px) {
    #feedback-area {
        padding-bottom: 4em;
    }
}
</style>
