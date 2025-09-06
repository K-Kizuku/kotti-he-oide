import { useState, useEffect } from 'react';
import { NotificationPermissionState } from '@/types/push';

export function useNotificationPermission() {
  const [permission, setPermission] = useState<NotificationPermissionState>('default');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!('Notification' in window)) {
      setPermission('unsupported');
      return;
    }

    setPermission(Notification.permission);

    // Listen for permission changes (if supported)
    if ('permissions' in navigator) {
      navigator.permissions.query({ name: 'notifications' as PermissionName })
        .then((permissionStatus) => {
          setPermission(permissionStatus.state as NotificationPermissionState);
          
          permissionStatus.addEventListener('change', () => {
            setPermission(permissionStatus.state as NotificationPermissionState);
          });
        })
        .catch(() => {
          // Fallback to Notification.permission
        });
    }
  }, []);

  const requestPermission = async (): Promise<NotificationPermissionState> => {
    if (!('Notification' in window)) {
      return 'unsupported';
    }

    if (Notification.permission === 'granted') {
      return 'granted';
    }

    setIsLoading(true);
    try {
      const result = await Notification.requestPermission();
      setPermission(result);
      return result;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    permission,
    isLoading,
    requestPermission,
    isSupported: permission !== 'unsupported',
    isGranted: permission === 'granted',
    isDenied: permission === 'denied',
    isDefault: permission === 'default'
  };
}