self.addEventListener('push', function(event) {
  const newMessage = `${event.data.text()}` + ' Shuttle will arrive in approximately 3 minutes';
  const now = new Date();
  const title = 'Shuttle Tracker';
  console.log('[Service Worker] Push Received.');
  console.log(`[Service Worker] Push had this data: "${event.data.text()}"`);
  const options = {
    body: newMessage,
  	icon: '~../assets/icon.svg',
   	requireInteraction: true,
   	timestamp: now.getHours() + ':' + now.getMinutes(),
  };
  event.waitUntil(self.registration.showNotification(title, options));
});