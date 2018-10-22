<template>
<li @click="visible = !visible;" class="dropdown-main-item"><a class="dropdown-icon"><img src="~../assets/menu.svg" class="dropdown-icon"></a>
<ul v-if="visible" class="dropdown-menu" style="display: inline-block;">
    <li id="dropdown-submenu-item_styling" class="dropdown-submenu-item">
        <p class="dropdown-menu-item_p">Information</p>
        <router-link style="text-decoration: none; color: #333333;" to="/about">
        <div id="darkmode-icon" class="dropdown-submenu-item_div">
            <p class="dropdown-submenu-subitem">About</p>
        </div>
        </router-link>
    </li>
    <li id="dropdown-menu-item_shuttle-schedule" class="dropdown-menu-item">
        <p class="dropdown-menu-item_p">Shuttle Schedules</p>
        <ul id="dropdown-submenu_shuttle-schedule" class="dropdown-submenu" style="display: inline-block;">
            <li id="dropdown-submenu-item_shuttle-schedule" class="dropdown-submenu-item">
                <p class="dropdown-submenu-item_p"><a target="_blank" rel="noopener noreferrer" href="https://info.rpi.edu/sites/default/files/2018-19MonThursCampusShuttleScheduleEastRoute.pdf"
                        class="dropdown-submenu-item_link">East: Monday-Friday</a></p>
            </li>
            <li id="dropdown-submenu-item_shuttle-schedule" class="dropdown-submenu-item">
                <p class="dropdown-submenu-item_p"><a target="_blank" rel="noopener noreferrer" href="https://info.rpi.edu/sites/default/files/2018-19MonThursCampusShuttleScheduleWestRoute.pdf"
                        class="dropdown-submenu-item_link">West: Monday-Friday</a></p>
            </li>
            <li id="dropdown-submenu-item_shuttle-schedule" class="dropdown-submenu-item">
                <p class="dropdown-submenu-item_p"><a target="_blank" rel="noopener noreferrer" href="https://info.rpi.edu/sites/default/files/2018-19WeekendLateNightShuttleSchedule.pdf"
                        class="dropdown-submenu-item_link">Weekend Late Night</a></p>
            </li>
        </ul>
    </li>
</ul>
</li>
</template>
<script lang="ts">
// This component is the dropdown menu displayed on the front page
import Vue from 'vue';
export default Vue.extend({
    data() {
        return {
            visible: false,
        } as {
            visible: boolean,
        };
    },
});
</script>

<style lang="scss">
    /* CSS for dropdown menus */
    ul {
    list-style: none none;
    }

    ul.dropdown-main {
    position: relative;
    margin: 0;
    padding: 0;

    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    }

    li.dropdown-main-item {
    background-color: rgb(221, 221, 221);
    width: 34px;
    float: left;
    }

    li.dropdown-main-item:hover {
    background-color: rgba(226, 226, 226, 0.85);
    cursor: pointer;
    }

    a.dropdown-icon {
    color: rgb(51, 51, 51);
    display: block;
    font-weight: 700;
    font-size: 12px;
    line-height: 32px;
    text-decoration: none;
    height: 34px;
    }

    img.dropdown-icon {
    box-sizing: border-box;
    padding: 5px;
    }

    ul.dropdown-menu {
    background-color: rgba(238, 238, 238, 0.90);
    display: none;
    position: absolute;
    top: 100%;
    left: 0;
    padding: 0;
    width: auto;
    white-space: nowrap;
    }

    li.dropdown-menu-item {
    box-sizing: border-box;
    float: none;
    width: auto;
    }

    p.dropdown-menu-item_p {
    height: 24px;
    margin: 3px 0;
    padding: 0 7px;
    font-size: large;
    font-weight: 500;
    }

    p.dropdown-menu-item_p:hover {
    background-color: rgba(226, 226, 226, 0.85);
    }

    ul.dropdown-submenu {
    display: none;
    position: static;
    top: 100%;
    left: 0;
    padding: 0 0 0 24px;
    width: calc( 100% - 24px );
    }

    div.dropdown-submenu-item_div {
    box-sizing: border-box;
    width: 100%;
    position: relative;
    height: 26px;
    margin: 3px 0;
    padding: 0 7px;
    }

    .dropdown-submenu-item_p {
    box-sizing: border-box;
    height: 24px;
    margin: 3px 0;
    padding: 0 7px;
    }

    a.dropdown-submenu-item_link {
    box-sizing: border-box;
    color: rgb(51, 51, 51);
    padding: 3px 7px;
    width: 100%;
    text-decoration: none;
    }

    a.dropdown-submenu-item_link:visited {
    color: rgb(51, 51, 51);
    text-decoration: none;
    }

    p.dropdown-submenu-item_p:hover {
    background-color: rgba(226, 226, 226, 0.85);
    }

    div.dropdown-submenu-item_div:hover {
    background-color: rgba(226, 226, 226, 0.85);
    }

    .dropdown-submenu-subitem {
    box-sizing: border-box;
    display: inline-block;
    height: 24px;
    margin: 0;
    padding: 0;
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    cursor: pointer;
    }

    img.dropdown-submenu-subitem {
    height: 20px;
    width: 32px;
    padding: 0 6px;
    }

    p.dropdown-submenu-subitem {
    left: 40px;
    padding: 0 8px;
    }

</style>

