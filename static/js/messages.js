
Vue.component('message-panel', {
  template:
  `<div class ="column">
    <div class="tile is-parent">
      <div class="tile box is-parent">
        <div class="tile notification" v-bind:class="{'is-danger':fail, 'is-white':!fail}">
          <div class="container">
          <p>Message: </p> <input class="textbox" v-model="message"></input> Characters Left: {{45 - message.length}}
          <br>
          <input type="checkbox" v-model="display">Enabled</input><br>
          <button class="button" @click="send">Submit</button>
          </div>
        </div>
      </div>
    </div>
  </div>`,
  data (){
    return {
      message: "",
      display: false,
      fail: false
    };
  },
  methods: {
    send: function(){
      let el = this;
      toSend = {Message: this.message, Display: this.display};
      $.post("/adminMessage",JSON.stringify(toSend)).then(resp =>{
        if (resp != "Success"){
          el.fail = true;
        }else{
          el.fail = false;
        }

      }
    ).catch(function(){
      el.fail = true;
    });
    }
  },
  mounted (){
    let el = this;
    fetch("/adminMessage").then(
      ret => {return ret.json();}).then(val =>{
        el.message = val.Message;
        el.display = val.Display;
      });
  }
});
