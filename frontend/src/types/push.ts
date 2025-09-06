export type NotificationPermissionState = 'default' | 'denied' | 'granted' | 'unsupported';

export interface PushSubscriptionState {
  isSubscribed: boolean;
  subscription: PushSubscription | null;
  subscriptionId: string | null;
  error: string | null;
  isLoading: boolean;
}

export interface NotificationPayload {
  title: string;
  body: string;
  icon?: string;
  badge?: string;
  tag?: string;
  data?: {
    url?: string;
    trackingId?: string;
    [key: string]: unknown;
  };
  actions?: Array<{
    action: string;
    title: string;
    icon?: string;
  }>;
}

export interface WebPushSupportInfo {
  isSupported: boolean;
  missingFeatures: string[];
}

// Service Worker 用の通知オプション（DOM の NotificationOptions とは別名にして衝突を避ける）
export interface SWNotificationOptions {
  body?: string;
  icon?: string;
  badge?: string;
  tag?: string;
  data?: Record<string, unknown>;
  actions?: Array<{
    action: string;
    title: string;
    icon?: string;
  }>;
  requireInteraction?: boolean;
  silent?: boolean;
  timestamp?: number;
}

export interface ServiceWorkerRegistration extends globalThis.ServiceWorkerRegistration {
  showNotification(title: string, options?: SWNotificationOptions): Promise<void>;
}

declare global {
  interface Window {
    swRegistration?: ServiceWorkerRegistration;
  }
}
