import { apiClient } from './client';
import { 
  SubscribeRequest, 
  SubscribeResponse, 
  VAPIDResponse, 
  UnsubscribeResponse,
  SendNotificationRequest,
  NotificationClickRequest
} from '@/types/api';

export const pushAPI = {
  async getVAPIDPublicKey(): Promise<string> {
    const data = await apiClient.get<VAPIDResponse>('/api/push/vapid-public-key');
    
    if (!data.success) {
      throw new Error(data.message);
    }
    
    return data.publicKey;
  },

  async subscribe(subscriptionData: SubscribeRequest): Promise<SubscribeResponse> {
    return apiClient.post<SubscribeResponse>('/api/push/subscribe', subscriptionData);
  },

  async unsubscribe(subscriptionId: string): Promise<UnsubscribeResponse> {
    return apiClient.delete<UnsubscribeResponse>(`/api/push/subscriptions/${subscriptionId}`);
  },

  // テスト用の通知送信（管理者用）
  async sendTestNotification(payload: SendNotificationRequest): Promise<{ success: boolean; message: string }> {
    return apiClient.post('/api/push/send', payload);
  },

  // 複数ユーザーへの通知送信（管理者用）
  async sendBatchNotification(payload: SendNotificationRequest): Promise<{ success: boolean; message: string }> {
    return apiClient.post('/api/push/send/batch', payload);
  },

  // 通知クリックトラッキング
  async trackNotificationClick(clickData: NotificationClickRequest): Promise<{ success: boolean; message: string }> {
    return apiClient.post('/api/push/click', clickData);
  }
};