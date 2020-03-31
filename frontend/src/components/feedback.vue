<template>
    <div class = "parent content">
        <h2>Feedback</h2>
        <p>
            <br>We'd love to hear back from you!
            <br>If you have any comments, suggestions, or concerns, submit them down below and they'll be sent directly to us.
        </p>
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
    </div>
</template>


<script lang="ts">
import Vue from 'vue';
export default Vue.extend({});
// maybe clear the page and show thanks message on button click
// function clearAndThanks(elementID)
// {
//     document.getElementById(elementID).innerHTML = "";
// }
</script>
<style lang="scss" scoped>
// <button class="submit" type="submit" form="feedback" value="Submit" id="submit" onclick="clearAndThanks('submit')" href="#popup1">Submit</button>
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