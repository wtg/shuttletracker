
Vue.component('message-panel', {
  template:
    `<div class ="column">
      <div class="box">
        <div class="notification" v-bind:class="{'is-danger':fail, 'is-white':!fail}">
          <div class="field">
            <label class="label">Message</label>
            <div class="control">
              <input class="textbox input" v-model="message">
            </div>
          </div>
          <div class="field">
            <div class="control">
              <label class="checkbox">
                <input type="checkbox" v-model="enabled">
                Enabled
              </label>
            </div>
          </div>
          <div class="field">
            <div class="control">
              <button class="button" @click="send">Submit</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>`,
  data() {
    return {
      message: "",
      enabled: false,
      fail: false
    };
  },
  methods: {
    send: function () {
      let el = this;
      toSend = { message: this.message, enabled: this.enabled };
      $.post("/adminMessage", JSON.stringify(toSend)).then(resp => {
        if (resp != "Success") {
          el.fail = true;
        } else {
          el.fail = false;
        }

      }
      ).catch(function () {
        el.fail = true;
      });
    }
  },
  mounted() {
    let el = this;
    fetch("/adminMessage").then(
      ret => { return ret.json(); }).then(val => {
        el.message = val.message;
        el.enabled = val.enabled;
      });
  }
});
