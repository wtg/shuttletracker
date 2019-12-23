<template>
  <div class="parent content">
    <h2>Spring 2020 Route Changes</h2>
    <h3>The Changes</h3>
    <p>Starting January 2020 the current East and West route system will be changed to include three routes: North, South, and West. An image of the new routes is shown below. The new routes can additionally be viewed on the Tracker by toggling them in the <a href="/settings">/settings</a> page.</p>

    <p>All Shuttle hours will remain the same. These route changes will be occurring only apply for Weekdays, the current East and West routes will remain on Weekends.</p>

    <p>The overall goal of these changes was to reduce shuttle loop time, remove underutilized stops, and add stops where there is current demand.</p>

    <p><strong>Note</strong>: These changes were updated on 12/23/2019. West route will now travel down Congress (instead of College Ave), and stops have been added at Winslow (5), West Hall (6), and Georgian Ct (13).</p>

    <img src="https://i.imgur.com/kQfvzxn.png"/>
    <br>
    <div id="text">
    <p><strong>West Route</strong>
        <br>
        <br>1. CBIS stop added by crosswalk
        <br>2. Polytech stop moved to northwest corner of intersection of 15th St and Congress St (Polytech Transfer)
        <br>3. The City Station stop moved across the street to City Station East by the parking garage entrance
        <br>4. Blitman Commons stop is moving to the bus shelter on 6th Ave
        <br>5. Winslow stop added at north side of building facing parking lot
        <br>6. West Hall stop added in front of West Hall lot
        <br>7. Troy stop moved to the east end of Ricketts building (Ricketts Transfer)
        <br>
    </p>
    <p><strong>North Route</strong>
        <br>
        <br>8. All stops along Sage Ave are being merged to in front of E-Complex
        <br>9. 9th St and Sage Ave will stop on 9th St
        <br>10. Alumni House stop is moving across the street to in front of the Alumni House
        <br>11. Beman Park stop added at 15th St and Jacob St
        <br>12. Colonie stop moved to directly in front of Colonie A on street
        <br>13. Stop on Georgian Ct added
        <br>14. Brinsmade Terrace, Sunset Terrace 1, and Sunset Terrace 2 are unchanged
        <br>15. E Lot stop is moving to the intersection of Cook Dr and Sunset Terrace
        <br>16. B Lot stop added by bus shelter
        <br>
    </p>
    <p><strong>South Route</strong>
        <br>
        <br>17. B Lot stop added by bus shelter
        <br>18. LXA stop is moving to the southwest corner of the intersection of Cook Dr and Sunset Terrace
        <br>19. Tibbits Ave stop is moving to intersection of Tibbits Ave and Orchard Ave
        <br>20. Polytech stop added in front of Muza
        <br>21. 15th St and College St stop added in front of Fire House</p>
    </div>

    <h3>Credit</h3>
    <p>These changes were designed and proposed by the Facilities and Services Committee of the Student Senate in collaboration with Parking and Transportation. Thank you to all members of FSC for their work over the past several months and to Parking and Transportation for being extremely receptive to feedback and very willing to make these changes.</p>
    <p class="u">Parking and Transportation</p>
      Michael Ramella, Director of Parking and Transportation and Auxiliary Services<br>
      Brent Shako, Operations Manager
    </p>
    <p class="u">Facilities and Services Committee</p>
    <p>
      Nick Dybas, Facilities and Services Chair <br>
      Izzy Dashiell, Transportation Subcommittee Chair<br>
      Evan Lazaro, Web Technologies Group Chair<br>
      Tejas Narayan<br>
      Colleen Corrigan<br>
      Iris Erlanger<br>
      Cait Bennett<br>
      Neha Keshan<br>
      Noah Cohen<br>
      Gretchen Forbush<br>
      Jakob Waryas<br>
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
