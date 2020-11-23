<template>
    <div class="parent">
        <h1 class="title">Feedback</h1>
        <h2 class="subtitle">All responses are greatly appreciated!<hr></h2>
        <div id="feedback-area">
            <div v-if="this.adminMessage !== undefined" >
                <b><div id="admin-msg" style="width: 100%;float:left;" v-html="this.adminMessage"></div></b>
            </div>
            <div id="thank-you" style="display:none">
                <p>Your submission has been recorded anonymously. Thank you for helping us improve Shuttle Tracker!</p>
            </div>
            <textarea
                id="form"
                v-model="feedbackMessage"
                placeholder="Type your response here..." 
                maxlength="512">
            </textarea>
            <br>
            <button @click="save" class="submit" type="submit" form="feedback" value="Submit"
                id="submit"><b>Submit</b></button>
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
        const wrapper = document.getElementById('feedback-area');
        const text = document.getElementById('form');
        const c_wrap  = document.createElement('div');
        const count   = document.createElement('span');
        c_wrap.id = 'char-count';
        c_wrap!.style.position = 'relative';
        c_wrap!.style.right = '0.78em';
        c_wrap!.style.bottom = '2.1em';
        c_wrap!.style.color = '#8a8a8a';
        c_wrap!.style.cssFloat = 'right';

        c_wrap!.appendChild(count);
        wrapper!.appendChild(c_wrap);

        function _set() {
            count!.innerHTML = (512 - (text as HTMLTextAreaElement).value.length || 0).toString();
            if ((text as HTMLTextAreaElement).value.length === 512) {
                const oldColor = count!.style.color;
                count!.style.color = 'var(--color-primary)';
                setTimeout(() => {
                    count!.style.color = oldColor;
                }, 150);
            }
        }

        text!.addEventListener('input', _set);
        _set.call(text);
    },
    methods: {
        save() {
            const myMessage = new Form(-1, this.feedbackMessage, new Date(), false);
            AdminServiceProvider.CreateForm(myMessage).then((resp) => {
                if (resp.ok) {
                    document.getElementById('char-count')!.style.display = 'none';
                    document.getElementById('admin-msg')!.style.display = 'none';
                    document.getElementById('form')!.style.display = 'none';
                    document.getElementById('submit')!.style.display = 'none';
                    (document.getElementById('submit') as HTMLButtonElement)!.disabled = true;
                    document.getElementById('thank-you')!.style.display = 'block';
                } else {
                    
                }
            }).catch(() => {
                
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

.parent{
    padding: 20px;
}
.subtitle{
    margin-bottom: 0 !important;
}
#admin-msg{
    font-size: 1.955em;
}
#feedback-area{
    margin-left: auto;
    margin-right: auto;
    width: 50%;
}
#thank-you{
    font-size:1.05em;
    margin-top: 10%;
}
#form{ 
    border-radius: 8px;
    outline: none !important;
    border:2px solid rgb(202, 202, 202);
    margin-top: 1em;
    margin-bottom: 0.25em;
    width:  100%;
    height: 16em; 
    padding: 5px;
    resize: none;
}
.submit{
     &:not(.no-hover):hover {
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
    height: 40px;
    width: 100px;
    border: none;
    border-radius: 8px;
    margin-top: .5em;
    transition: all 0.35s;
    margin-left: calc(50% - 50px);
}

@media only screen and (max-width: 1024px) {
    #feedback-area{ 
        margin-bottom: 0.5em;
        width: 95%;
    }
}
@media only screen and (max-width: 420px) {
    #form{
        margin-top: 0.75em;
        padding-right: 2.28em;
    }
    #feedback-area{
        margin-bottom: 0.5em;
    }
    #admin-msg{
        font-size: 1.2em; 
    }
}
@media only screen and (max-height: 700px) {
    #feedback-area {
        padding-bottom: 4em;
    }
}

@media only screen and (max-height: 450px) {
    #admin-msg {
        font-size: 1.2em;    
    }
}

@media only screen and (min-height: 700px) and (max-height: 800px) {
    #form {
        height: 19em; 
    }
}

@media only screen and (min-height: 800px) and (max-height: 1000px) {
    #form {
        height: 25em; 
    }
}
</style>
