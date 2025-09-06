const STORAGE_KEYS = {
  PUSH_SUBSCRIPTION_ID: 'push_subscription_id',
  NOTIFICATION_SETTINGS: 'notification_settings',
  VAPID_PUBLIC_KEY: 'vapid_public_key',
} as const;

export const storage = {
  // Push subscription ID management
  setPushSubscriptionId(id: string): void {
    try {
      localStorage.setItem(STORAGE_KEYS.PUSH_SUBSCRIPTION_ID, id);
    } catch (error) {
      console.error('Failed to store subscription ID:', error);
    }
  },

  getPushSubscriptionId(): string | null {
    try {
      return localStorage.getItem(STORAGE_KEYS.PUSH_SUBSCRIPTION_ID);
    } catch (error) {
      console.error('Failed to get subscription ID:', error);
      return null;
    }
  },

  removePushSubscriptionId(): void {
    try {
      localStorage.removeItem(STORAGE_KEYS.PUSH_SUBSCRIPTION_ID);
    } catch (error) {
      console.error('Failed to remove subscription ID:', error);
    }
  },

  // VAPID public key management
  setVapidPublicKey(key: string): void {
    try {
      localStorage.setItem(STORAGE_KEYS.VAPID_PUBLIC_KEY, key);
    } catch (error) {
      console.error('Failed to store VAPID key:', error);
    }
  },

  getVapidPublicKey(): string | null {
    try {
      return localStorage.getItem(STORAGE_KEYS.VAPID_PUBLIC_KEY);
    } catch (error) {
      console.error('Failed to get VAPID key:', error);
      return null;
    }
  },

  removeVapidPublicKey(): void {
    try {
      localStorage.removeItem(STORAGE_KEYS.VAPID_PUBLIC_KEY);
    } catch (error) {
      console.error('Failed to remove VAPID key:', error);
    }
  },

  // Notification settings management
  setNotificationSettings(settings: Record<string, unknown>): void {
    try {
      localStorage.setItem(STORAGE_KEYS.NOTIFICATION_SETTINGS, JSON.stringify(settings));
    } catch (error) {
      console.error('Failed to store notification settings:', error);
    }
  },

  getNotificationSettings(): Record<string, unknown> | null {
    try {
      const settings = localStorage.getItem(STORAGE_KEYS.NOTIFICATION_SETTINGS);
      return settings ? JSON.parse(settings) : null;
    } catch (error) {
      console.error('Failed to get notification settings:', error);
      return null;
    }
  },

  removeNotificationSettings(): void {
    try {
      localStorage.removeItem(STORAGE_KEYS.NOTIFICATION_SETTINGS);
    } catch (error) {
      console.error('Failed to remove notification settings:', error);
    }
  },

  // Clear all storage
  clearAll(): void {
    try {
      Object.values(STORAGE_KEYS).forEach(key => {
        localStorage.removeItem(key);
      });
    } catch (error) {
      console.error('Failed to clear storage:', error);
    }
  },

  // Check if localStorage is available
  isAvailable(): boolean {
    try {
      const test = '__storage_test__';
      localStorage.setItem(test, test);
      localStorage.removeItem(test);
      return true;
    } catch {
      return false;
    }
  }
};