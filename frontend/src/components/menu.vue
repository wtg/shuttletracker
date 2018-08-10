<template>
    <div class='dropdown'>
        <ul class='dropdown-main'>
            <li class='dropdown-main-item'>
                <a v-on:click='toggleDropdownMenuVisibility()' class='dropdown-icon'>
                    <img class='dropdown-icon' src='../assets/images/menu.svg'>
                </a>
                <ul class='dropdown-menu'>
                    <li class='dropdown-submenu-item' id='dropdown-submenu-item_styling'>
                        <p class='dropdown-menu-item_p'>Information</p>
                        <div class='dropdown-submenu-item_div' id='darkmode-icon'>
                            <p @click='$emit(' displaymodal ')' class='dropdown-submenu-subitem'>About</p>
                        </div>
                    </li>
                    <li class='dropdown-menu-item' id='dropdown-menu-item_shuttle-schedule'>
                        <p class='dropdown-menu-item_p'>Shuttle Schedules</p>
                        <!-- http://www.rpi.edu/dept/parking/shuttle/ -->
                        <ul class='dropdown-submenu' id='dropdown-submenu_shuttle-schedule'>
                            <li class='dropdown-submenu-item' :key='item.name' id='dropdown-submenu-item_shuttle-schedule' v-for='item in list_data'>
                                <p class='dropdown-submenu-item_p'><a class='dropdown-submenu-item_link' target='_blank' rel='noopener noreferrer' :href='item.link'>{{item.name}}</a></p>
                            </li>
                        </ul>
                    </li>
                    <li class='dropdown-menu-item' id='dropdown-menu-item_styling'>
                        <p class='dropdown-menu-item_p'>Styling</p>
                        <!-- for changing the view of the page -->
                        <ul class='dropdown-submenu' id='dropdown-submenu_styling'>
                            <li class='dropdown-submenu-item' id='dropdown-submenu-item_styling'>
                                <div class='dropdown-submenu-item_div' id='darkmode-icon'>
                                    <img v-on:click='toggleDarkmode' class='dropdown-submenu-subitem' :src='moonicon'>
                                    <p v-on:click='toggleDarkmode' class='dropdown-submenu-subitem'>Darkmode</p>
                                </div>
                            </li>
                        </ul>
                    </li>
                </ul>
            </li>
        </ul>
    </div>
</template>

<script lang='ts'>
    import Vue from 'vue';
    export default Vue.extend({
        data() {
            return {
                list_data: [{
                        name: 'East: Monday-Thursday',
                        link: 'http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleEastRoute.pdf',
                    },
                    {
                        name: 'East: Friday',
                        link: 'http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyEastShuttleSchedule.pdf',
                    },
                    {
                        name: 'West: Monday-Thursday',
                        link: 'http://www.rpi.edu/dept/parking/shuttle/2017-2018CampusShuttleScheduleWestRoute.pdf',
                    },
                    {
                        name: 'West: Friday',
                        link: 'http://www.rpi.edu/dept/parking/shuttle/2017-2018FridayOnlyWestShuttleSchedule.pdf',
                    },
                    {
                        name: 'Weekend Late Night',
                        link: 'http://www.rpi.edu/dept/parking/shuttle/2017-2018Weekend-LateNightShuttleSchedule.pdf',
                    },
                ],
                title: 'RPI Shuttle Tracker',
                moonicon: 'static/images/moon.svg',
                sunicon: 'static/images/sun.svg',
                darkmodeOn: 0,
                menuVisibility: 0,
                dropdownMenuList: document.getElementsByClassName('dropdown-menu'),
                dropdownSubmenuList: document.getElementsByClassName('dropdown-submenu'),
    
            };
        },
        methods: {
            toggleDropdownMenuVisibility() {
                // 0 is closed, 1 is open
                if (this.menuVisibility === 0) {
                    this.menuVisibility = 1;
                    for (var i = 0; i < this.dropdownMenuList.length; i++) {
                        this.dropdownMenuList[i].style.display = 'inline-block';
                    }
                    for (var i = 0; i < this.dropdownSubmenuList.length; i++) {
                        this.dropdownSubmenuList[i].style.display = 'inline-block';
                    }
                } else {
                    this.menuVisibility = 0;
                    for (var i = 0; i < this.dropdownMenuList.length; i++) {
                        this.dropdownMenuList[i].style.display = 'none';
                    }
                    for (var i = 0; i < this.dropdownSubmenuList.length; i++) {
                        this.dropdownSubmenuList[i].style.display = 'none';
                    }
                }
            },
        },
    });
</script>

