var Admin = {

  grabStops: function(){
    $.get( "http://127.0.0.1:8080/stops", Admin.populateStops);

  },

  populateStops: function(data){
    console.log(data);

  }


}

$(document).ready(function(){
  Admin.grabStops();
});
