<template>
    <div class = "parent content">
        <h2>Feedback</h2>
        <div v-if="(this.message !== undefined && this.message.enabled === true)" >
            <div style="width: 100%;float:left;" v-html="this.message.message"></div>
        </div>
        <!-- <p>
            <br>We'd love to hear back from you!
            <br>If you have any comments, suggestions, or concerns, submit them down below and they'll be sent directly to us.
        </p> -->
        <p>
            <td style="padding-top: 5px;">
            <textarea 
                placeholder="beep beep I'm a sheep . . ." 
                rows = "8" cols = "127">
            </textarea>
            <br>
            <button class="submit" type="submit" form="feedback" value="Submit"
            id="submit">Submit</button>
            </td>
        </p>
    </div>
</template>


<script lang="ts">
import Vue from 'vue';
import FeedbackMessageUpdate from '@/structures/feedbackMessageUpdate';
export default Vue.extend({
    // what does this do?
    data() {
        return {
            hide: false,
        };
    },
    computed: {
        message(): FeedbackMessageUpdate {
            return this.$store.state.feedbackMessage;
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