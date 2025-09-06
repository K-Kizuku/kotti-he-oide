// Service Worker for Web Push Notifications

// Push event handler
self.addEventListener('push', (event) => {
  const options = {
    body: 'Default notification body',
    icon: '/icons/icon-192.png',
    badge: '/icons/badge-72.png',
    tag: 'default',
    data: {},
    actions: [
      {
        action: 'open',
        title: '開く'
      },
      {
        action: 'close',
        title: '閉じる'
      }
    ]
  };

  let payload = null;
  if (event.data) {
    try {
      payload = event.data.json();
      Object.assign(options, {
        body: payload.body || options.body,
        icon: payload.icon || options.icon,
        tag: payload.tag || options.tag,
        data: payload.data || options.data,
        actions: payload.actions || options.actions
      });
    } catch (error) {
      console.error('Failed to parse push data:', error);
    }
  }

  event.waitUntil(
    self.registration.showNotification(
      payload?.title || 'New Notification',
      options
    )
  );
});

// Notification click handler  
self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  const clickAction = event.action || 'open';
  const notificationData = event.notification.data || {};
  
  if (clickAction === 'close') {
    return;
  }

  const urlToOpen = notificationData.url || '/';
  
  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true })
      .then((clientList) => {
        // Try to focus existing tab
        for (const client of clientList) {
          if (client.url === urlToOpen && 'focus' in client) {
            return client.focus();
          }
        }
        
        // Open new tab if no existing tab found
        if (clients.openWindow) {
          return clients.openWindow(urlToOpen);
        }
      })
  );

  // Track notification click
  if (notificationData.trackingId) {
    fetch('/api/push/click', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        trackingId: notificationData.trackingId,
        action: clickAction,
        timestamp: new Date().toISOString()
      })
    }).catch(error => {
      console.error('Failed to track click:', error);
    });
  }
});

// Install event handler
self.addEventListener('install', () => {
  console.log('Service Worker installing');
  self.skipWaiting();
});

// Activate event handler
self.addEventListener('activate', (event) => {
  console.log('Service Worker activating');
  event.waitUntil(clients.claim());
});