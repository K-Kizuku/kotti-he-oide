export interface SubscribeRequest {
  endpoint: string;
  keys: {
    p256dh: string;
    auth: string;
  };
  ua?: string;
  expirationTime?: number;
}

export interface SubscribeResponse {
  id: string;
  success: boolean;
  message: string;
}

export interface VAPIDResponse {
  publicKey: string;
  success: boolean;
  message: string;
}

export interface UnsubscribeResponse {
  success: boolean;
  message: string;
}

export interface SendNotificationRequest {
  title: string;
  body: string;
  icon?: string;
  badge?: string;
  tag?: string;
  data?: Record<string, unknown>;
  actions?: NotificationAction[];
  url?: string;
}

export interface NotificationAction {
  action: string;
  title: string;
  icon?: string;
}

export interface NotificationClickRequest {
  trackingId: string;
  action: string;
  timestamp: string;
}