'use client';

import { useNotificationPermission } from '@/hooks/useNotificationPermission';

export default function NotificationPermissionButton() {
  const { isLoading, requestPermission, isSupported, isGranted, isDenied } = useNotificationPermission();

  if (!isSupported) {
    return (
      <div className="p-4 bg-gray-100 rounded-lg">
        <p className="text-gray-600">お使いのブラウザは通知機能をサポートしていません</p>
      </div>
    );
  }

  if (isGranted) {
    return (
      <div className="p-4 bg-green-100 rounded-lg">
        <p className="text-green-800">✅ 通知が許可されています</p>
      </div>
    );
  }

  if (isDenied) {
    return (
      <div className="p-4 bg-red-100 rounded-lg">
        <p className="text-red-800">❌ 通知が拒否されています</p>
        <p className="text-sm text-red-600 mt-2">
          ブラウザの設定から通知を許可してください
        </p>
      </div>
    );
  }

  return (
    <div className="p-4 bg-blue-100 rounded-lg">
      <p className="text-blue-800 mb-3">通知を受け取るには許可が必要です</p>
      <button
        onClick={requestPermission}
        disabled={isLoading}
        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 transition-colors"
      >
        {isLoading ? '確認中...' : '通知を許可する'}
      </button>
    </div>
  );
}