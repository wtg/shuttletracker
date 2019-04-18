self.addEventListener('push', event => {
  const newMessage = 'Shuttle will arrive in approximately 5 minutes';
  const now = new Date();
  const title = 'Shuttle Tracker';
  const options = {
    body: newMessage,
  	icon: '~../assets/icon.svg',
   	requireInteraction: true,
   	timestamp: now.getHours() + ':' + now.getMinutes(),
  };
  event.waitUntil(self.registration.showNotification(title, options));
});