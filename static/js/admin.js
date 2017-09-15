var Admin = {
  RoutesMap: null,
  /*
  <div class = "route-description-box">
    <span class = "emphasis">name:</span><span class ="content"> West Campus</span><br>
    <span class = "emphasis">description:</span><span class ="content"> The west campus Red Hawk shuttle will make stops at the Student Union Horseshoe, Sage Ave & Troy Building crosswalks, Blitman Residence Commons, 6th Avenue and City Station, Polytechnic Residence Commons, and 15th Street and College Avenue. This shuttle route will have a yellow indicator labeled "WEST" on the front of the shuttle.</span><br>
    <span class = "emphasis">enabled:</span><span class="content">true</span><br>
    <span class = "emphasis">color:</span><span class="content">#fff</span><br>
    <span class = "emphasis">time:</span><span class="content">idek</span><br>
    <span class = "emphasis">id:</span><span class="content">13445213a</span><br>
    <div style="float: right;"><a class="button">edit</a><a class="button">delete</a></div><br>

  </div>
  */
  populateRoutesPanel: function(data){
    console.log(data);
    $(".routePanel").html("");
    if(data == null){

    }else{
      for(var i = 0; i < data.length; i ++){
        console.log(data[i]);
        Admin.buildRouteBox(data[i]);
      }
      $(".deleteroute").click(function(){
        $.ajax({
          url: '/routes/' + $(this).attr("routeId"),
          type: 'DELETE',
          success: function(result) {
            $.get( "/routes", Admin.populateRoutesPanel);
          }
        });
      });
    }
  },

  buildRouteBox: function(routeInfo){
    var box = "";
    box += "<div id = " + routeInfo.id +" class = 'route-description-box'>";
    box += "<span class = 'emphasis'>name:</span><span class ='content'> " + routeInfo.name + "</span><br>";
    box += "<span class = 'emphasis'>description:</span><span class ='content'>" + routeInfo.description + "</span><br>"
    box += "<span class = 'emphasis'>enabled:</span><span class='content'>"+routeInfo.enabled + "</span><br>";
    box += "<span class = 'emphasis'>color:</span><span class='content'>" + routeInfo.color + "</span><br>";
    box += "<span class = 'emphasis'>time:</span><span class='content'>lel</span><br>";
    box += "<span class = 'emphasis'>id:</span><span class='content'>"+ routeInfo.id + "</span><br>";
    box += "<div style='float: right;width:auto;'><button class='button cbutton'>edit</button><button id='delete' routeId="+routeInfo.id +" class='button cbutton deleteroute'>delete</button></div><br></div>"
    $(".routePanel").append(box);

  },

  routeDeleteHandler: function(info){
    console.log(info);
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
  $.get( "/routes", Admin.populateRoutesPanel);


});
