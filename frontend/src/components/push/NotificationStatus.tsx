'use client';

import { useNotificationPermission } from '@/hooks/useNotificationPermission';
import { usePushNotification } from '@/hooks/usePushNotification';
import { checkWebPushSupport } from '@/lib/utils/notificationUtils';
import { useEffect, useState } from 'react';

interface WebPushSupportInfo {
  isSupported: boolean;
  missingFeatures: string[];
}

export default function NotificationStatus() {
  const { permission } = useNotificationPermission();
  const { isSubscribed, subscription, subscriptionId, isLoading } = usePushNotification();
  const [supportInfo, setSupportInfo] = useState<WebPushSupportInfo | null>(null);

  useEffect(() => {
    setSupportInfo(checkWebPushSupport());
  }, []);

  if (isLoading) {
    return (
      <div className="p-4 bg-gray-50 rounded-lg">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/3 mb-2"></div>
          <div className="h-3 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {/* Web Push サポート状況 */}
      <div className="p-3 bg-gray-50 rounded">
        <h3 className="font-medium text-sm text-gray-700 mb-2">ブラウザサポート状況</h3>
        <div className="space-y-1 text-xs">
          <div className="flex justify-between">
            <span>Web Push API:</span>
            <span className={supportInfo?.isSupported ? 'text-green-600' : 'text-red-600'}>
              {supportInfo?.isSupported ? '✅ 対応' : '❌ 非対応'}
            </span>
          </div>
          {!supportInfo?.isSupported && supportInfo?.missingFeatures && (
            <p className="text-red-600 text-xs">
              不足機能: {supportInfo.missingFeatures.join(', ')}
            </p>
          )}
        </div>
      </div>

      {/* 通知許可状況 */}
      <div className="p-3 bg-gray-50 rounded">
        <h3 className="font-medium text-sm text-gray-700 mb-2">通知許可状況</h3>
        <div className="flex justify-between text-xs">
          <span>許可状態:</span>
          <span className={`font-medium ${
            permission === 'granted' ? 'text-green-600' : 
            permission === 'denied' ? 'text-red-600' : 
            permission === 'default' ? 'text-yellow-600' : 'text-gray-600'
          }`}>
            {permission === 'granted' && '✅ 許可済み'}
            {permission === 'denied' && '❌ 拒否済み'}
            {permission === 'default' && '⏳ 未設定'}
            {permission === 'unsupported' && '❌ 非対応'}
          </span>
        </div>
      </div>

      {/* 購読状況 */}
      <div className="p-3 bg-gray-50 rounded">
        <h3 className="font-medium text-sm text-gray-700 mb-2">プッシュ購読状況</h3>
        <div className="space-y-1 text-xs">
          <div className="flex justify-between">
            <span>購読状態:</span>
            <span className={`font-medium ${isSubscribed ? 'text-green-600' : 'text-gray-600'}`}>
              {isSubscribed ? '✅ 有効' : '❌ 無効'}
            </span>
          </div>
          {subscriptionId && (
            <div className="flex justify-between">
              <span>購読ID:</span>
              <span className="text-gray-500 break-all">
                {subscriptionId.substring(0, 8)}...
              </span>
            </div>
          )}
          {subscription && (
            <div className="mt-2 p-2 bg-white rounded text-xs">
              <p className="text-gray-600 mb-1">エンドポイント:</p>
              <p className="text-gray-500 break-all font-mono text-[10px]">
                {subscription.endpoint}
              </p>
            </div>
          )}
        </div>
      </div>

      {/* 全体的なステータス */}
      <div className="p-3 rounded border-l-4 border-l-blue-500 bg-blue-50">
        <div className="flex items-center space-x-2">
          <div className={`w-3 h-3 rounded-full ${
            supportInfo?.isSupported && permission === 'granted' && isSubscribed 
              ? 'bg-green-500' : 'bg-yellow-500'
          }`}></div>
          <span className="text-sm font-medium text-blue-800">
            {supportInfo?.isSupported && permission === 'granted' && isSubscribed 
              ? 'Web Push通知が利用可能です' 
              : 'Web Push通知の設定が必要です'}
          </span>
        </div>
      </div>
    </div>
  );
}