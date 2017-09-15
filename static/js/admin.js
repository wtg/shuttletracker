var Admin = {
  RoutesMap: null,

  buildRouteBox: function(){

  },

  initMap: function(){
    Admin.RoutesMap = L.map('mapid', {
        zoomControl: false,
        attributionControl: false
    });
    Admin.RoutesMap.setView([42.728172, -73.678803], 15.3);
    Admin.RoutesMap.addControl(L.control.attribution({
        position: 'bottomright',
        prefix: ''
    }));
    L.tileLayer('http://tile.stamen.com/toner-lite/{z}/{x}/{y}{r}.png', {
      attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under <a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
      maxZoom: 17,
      minZoom: 14
    }).addTo(Admin.RoutesMap);

  }

};

$(document).ready(function(){
  Admin.initMap();
});
