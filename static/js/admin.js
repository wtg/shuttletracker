var Admin = {

  grabStops: function(){
    $.get( "http://127.0.0.1:8080/stops", Admin.populateStops);

  },

  populateStops: function(data){
    console.log(data);
    for(var i = 0; i < data.length; i ++){
      $("#Stops").append("<div class = 'stop' id='123'><b> Stop 1</b><p> Stop Information here</p><button class = 'rbtn'><i class='fa fa-trash' aria-hidden='true'></i></button><button class = 'rbtn'><i class='fa fa-pencil' aria-hidden='true'></i></button></div>");
      if(i % 2 == 0 && i != 0){
        $("#Stops").append("<br style='clear:both;'/>")
      }
    }
  }


}

$(document).ready(function(){
  Admin.grabStops();
});
