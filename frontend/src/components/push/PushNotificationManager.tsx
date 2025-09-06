'use client';

import { usePushNotification } from '@/hooks/usePushNotification';
import { useNotificationPermission } from '@/hooks/useNotificationPermission';
import NotificationPermissionButton from './NotificationPermissionButton';

export default function PushNotificationManager() {
  const { isGranted: hasPermission } = useNotificationPermission();
  const { 
    isSubscribed, 
    isLoading, 
    error, 
    subscribe, 
    unsubscribe,
    subscription 
  } = usePushNotification();

  if (!hasPermission) {
    return (
      <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg">
        <h2 className="text-xl font-bold mb-4">プッシュ通知設定</h2>
        <NotificationPermissionButton />
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg">
      <h2 className="text-xl font-bold mb-4">プッシュ通知設定</h2>
      
      {error && (
        <div className="p-3 mb-4 bg-red-100 border border-red-400 text-red-700 rounded">
          エラー: {error}
        </div>
      )}
      
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <span className="text-gray-700">通知状態:</span>
          <span className={`px-2 py-1 rounded text-sm font-medium ${
            isSubscribed ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
          }`}>
            {isSubscribed ? '有効' : '無効'}
          </span>
        </div>
        
        {subscription && (
          <div className="text-xs text-gray-500 p-2 bg-gray-50 rounded">
            <p className="font-medium mb-1">購読情報:</p>
            <p className="break-all">エンドポイント: {subscription.endpoint.substring(0, 50)}...</p>
          </div>
        )}
        
        <div className="flex space-x-2">
          {!isSubscribed ? (
            <button
              onClick={subscribe}
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50 transition-colors font-medium"
            >
              {isLoading ? '設定中...' : '通知を有効化'}
            </button>
          ) : (
            <button
              onClick={unsubscribe}
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 transition-colors font-medium"
            >
              {isLoading ? '解除中...' : '通知を無効化'}
            </button>
          )}
        </div>
        
        {isSubscribed && (
          <div className="text-sm text-gray-600 p-3 bg-blue-50 rounded">
            <p className="font-medium text-blue-800">✅ 通知の準備ができました</p>
            <p className="mt-1">サーバーから送信される通知を受け取ることができます。</p>
          </div>
        )}
      </div>
    </div>
  );
}