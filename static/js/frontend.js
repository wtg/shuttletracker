Vue.component('dropdown-menu',{
  template: `<div>
      <ul class="dropdown">
          <li>
              <a href="#" class="dropdown-schedule">
                  <img src="static/images/menu.svg">
              </a>
              <ul class="dropdown-menu">
                  <li>
                      <p id="schedule-menu">Shuttle Schedules</p>
                      <!-- http://www.rpi.edu/dept/parking/shuttle/ -->
                  </li>
                  <li v-for="item in list_data">
                      <p><a target="_blank" rel="noopener noreferrer" :href="item.link">{{item.name}}</a></p>
                  </li>
              </ul>
          </li>
      </ul>
  </div>`,
  data (){
    return{
        list_data: [
          {name: "East: Monday-Thursday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleEastRoute.pdf"},
          {name: "East: Friday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyEastShuttleSchedule.pdf"},
          {name: "West: Monday-Thursday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleWestRoute.pdf"},
          {name: "West: Friday", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyWestShuttleSchedule.pdf"},
          {name: "Weekend Late Night", link: "http://www.rpi.edu/dept/parking/shuttle/2017-2018Weekend-LateNightShuttleSchedule.pdf"}
        ],

    }
  }
})

Vue.component('title-bar', {
  template:
  `  <div class="titleBar">
        <div class="titleContent">
            <dropdown-menu></dropdown-menu>
            <p class="title">{{title}}</p>
            <a href="https://webtech.union.rpi.edu" class="logo">
                <img src="static/images/wtg.svg">
            </a>
        </div>
    </div>`,
    data (){
      return {
        title: "RPI Shuttle Tracker"
      }
    }

})

var ShuttleTracker = new Vue({
  el: '#app-vue',
  data: {
  }

});
