<template>
  <div class="parent content">
    <h2>Important Notes and Changes for Fall 2020</h2>
    <h3>What's New</h3>
    <p>
      Given the current events, there have been some substantial changes made to the Shuttle routes and the Shuttle system to help keep Covid-19 in check. 
      This page serves to detail the changes that have been made this semester along with general question both answered and unanswered by Parking & Transporation. 
      If you have any additional questions, feel free to contact FSC Chair Cait Bennett (<a href="mailto:bennec6@rpi.edu">bennec6@rpi.edu</a>) or 
      WTG Chair Sean Heffernan (<a href="mailto:heffes2@rpi.edu">heffes2@rpi.edu</a>)</p>

    <h4><strong>Important Route Changes</strong></h4>
    <p>There will only be two Shuttle routes this semester: North and West, which will be reflected on the Shuttle Tracker.</p>
    <div id="text">
      <p>The following stops have been removed from North Route:
        <br>1. 9th/Sage
        <br>2. 13th/Peoples
        <br>3. Georgian
        <br>4. LXA
      </p>
      <p>The following stops have been removed from West Route:
        <br>1. CBIS
        <br>2. 15th/College
        <br>3. Winslow
        <br>4. West Hall
      </p>
      <p>According to the Department of Parking & Transportation, these stops were removed for efficiency given the decreased capacity of the shuttles.</p>
    </div>

    <h4><strong>Hudson Valley College Suites Shuttles</strong></h4>
    <p>
      To ease transportation for students living at College Suites in Hudson Valley, two temporary shuttles will be running back and forth between 
      Hudson Valley College Suites and Campus.
    </p>
    <p>The shuttle will stop at the front door of College Suites and travel to the Student Union.</p>
    <p>
      These two shuttles will be tracked by Shuttle Tracker, but the route they run on will <strong>not</strong> be displayed on 
      the Shuttle Tracker per the request of Parking and Transportation.
    </p>

    <h4><strong>CDTA Express Route</strong></h4>
    <p>An express route will be offerred Monday to Friday 7:00am - 7:00pm aboard two RPI student only CDTA buses.</p>
    <p>Bus route will be City Station &rarr; Ricketts &rarr; East Campus &rarr; City Station.</p>
    <p>Occupancy will be limited to 26 riders, masks and RPI IDs are required.</p>
    <p>Please note that these buses will <strong>NOT</strong> be tracked by Shuttle Tracker!</p>

    <h4><strong>Schedule Changes</strong></h4>
    <p>There will be no set schedule, or "paper schedule" this semester.</p>

    <div id="text">
    <p>The hours of operation for Shuttles this semester is as follows:
      <br>
      <br><strong>On Campus:</strong>
      <br>7am - 12am (Weekdays)
      <br>9am - 12am (Saturday)
      <br>9am - 8pm (Sunday)
      <br><strong>College Suites:</strong>
      <br>7am - 7pm (All Week)
      <br><strong>CDTA Express</strong>
      <br>7am - 7pm (Weekdays)
    </p>
    <p>Weekend times for On Campus shuttles are subject to change due to Safe Ride scheduling. This page will be updated as necessary.</p>
    <p>Please note that this information can also be viewed on Shuttle Tracker's <a href="/schedules">/schedules</a> page.</p>
    </div>


    <h3>FAQ</h3>
    <p>
      We are working to clarify certain elements of the shuttle system with the Department of Parking and Transportation, and will update 
      this document as more information becomes available. Some of these questions include:
      <br>How many Shuttles will be running on each route?
      <br>How were the removed stops chosen?
      <br>What will the availability of on call Shuttles be?
      <br>More details about weekend Shuttle Schedule/Safe Ride
    </p>

    <p>
      As always, if you have any questions, feel free to contact FSC Chair Cait Bennett (<a href="mailto:bennec6@rpi.edu">bennec6@rpi.edu</a>) 
      or WTG Chair Sean Heffernan (<a href="mailto:heffes2@rpi.edu">heffes2@rpi.edu</a>)
    </p>

    
    <p class="u">Student Senate</p>
    <p>
      Cait Bennett, Facilities and Services Chair <br>
      Sean Heffernan, Web Technologies Group Chair<br>
    </p>
    <p class="u">Parking and Transportation</p>
    <p>
      Michael Ramella, Director of Parking and Transportation and Auxiliary Services<br>
      Brent Shako, Operations Manager
    </p>

  </div>
</template>
<script lang="ts">
// This component is the about window that is displayed when the `about` button on the menu is selected

import Vue from 'vue';
export default Vue.extend({});
</script>
<style lang="scss" scoped>
.parent {
  padding: 20px;
}
@media only screen and (min-width: 70em) {
  .parent {
    padding: 20px 20em;
  }
}
p {
  margin-bottom: 1em;
}
.u {
  text-decoration: underline;
}
img{
  width:90%;
  max-width: 60em;
  margin-left:auto;
  margin-right:auto;
}
</style>
