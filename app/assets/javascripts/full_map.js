function initializeMap() {
  var options = {
    center: {
      lat: 42.730172,
      lng: -73.678803
    },
    zoom: 15
  };

  var map = new google.maps.Map($('#map-canvas')[0], options);
}

google.maps.event.addDomListener(window, 'load', initializeMap);