'use client';

import PushNotificationManager from '@/components/push/PushNotificationManager';
import NotificationStatus from '@/components/push/NotificationStatus';
import { checkWebPushSupport } from '@/lib/utils/notificationUtils';
import { useEffect, useState } from 'react';
import Link from 'next/link';

interface WebPushSupportInfo {
  isSupported: boolean;
  missingFeatures: string[];
}

export default function NotificationsPage() {
  const [supportInfo, setSupportInfo] = useState<WebPushSupportInfo | null>(null);
  const [showStatus, setShowStatus] = useState(false);

  useEffect(() => {
    setSupportInfo(checkWebPushSupport());
  }, []);

  if (!supportInfo) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-md mx-auto">
          <div className="animate-pulse bg-gray-200 h-8 rounded mb-4"></div>
          <div className="animate-pulse bg-gray-200 h-32 rounded"></div>
        </div>
      </div>
    );
  }

  if (!supportInfo.isSupported) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold mb-6 text-center">通知設定</h1>
        <div className="max-w-md mx-auto bg-red-50 border border-red-200 text-red-700 px-4 py-6 rounded-lg">
          <div className="text-center">
            <div className="text-4xl mb-3">❌</div>
            <p className="font-bold text-lg mb-2">お使いのブラウザは Web Push 通知をサポートしていません</p>
            <p className="text-sm">不足している機能: {supportInfo.missingFeatures.join(', ')}</p>
            <div className="mt-4 p-3 bg-red-100 rounded text-xs text-left">
              <p className="font-medium mb-2">推奨ブラウザ:</p>
              <ul className="space-y-1">
                <li>• Chrome 42+</li>
                <li>• Firefox 44+</li>
                <li>• Safari 16.4+ (macOS 13+, iOS 16.4+)</li>
                <li>• Edge 79+</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold mb-2 text-center">通知設定</h1>
        <p className="text-gray-600 text-center mb-8">
          Web Push通知を設定して、重要な更新情報を受け取りましょう
        </p>
        
        <div className="space-y-6">
          {/* メイン設定コンポーネント */}
          <PushNotificationManager />
          
          {/* ステータス表示の切り替え */}
          <div className="max-w-md mx-auto">
            <button
              onClick={() => setShowStatus(!showStatus)}
              className="w-full px-4 py-2 text-sm text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded transition-colors"
            >
              {showStatus ? '詳細情報を非表示' : '詳細情報を表示'} 
              <span className="ml-1">{showStatus ? '▼' : '▶'}</span>
            </button>
            
            {showStatus && (
              <div className="mt-4">
                <NotificationStatus />
              </div>
            )}
          </div>
          
          {/* 説明セクション */}
          <div className="max-w-md mx-auto bg-gray-50 p-4 rounded-lg">
            <h3 className="font-medium text-gray-800 mb-3">Web Push通知について</h3>
            <div className="space-y-2 text-sm text-gray-600">
              <p>• ブラウザを閉じていても通知を受け取れます</p>
              <p>• 通知の許可はいつでも変更できます</p>
              <p>• 個人情報は送信されません</p>
              <p>• HTTPS接続が必要です</p>
            </div>
            
            <div className="mt-4 pt-3 border-t border-gray-200">
              <p className="text-xs text-gray-500">
                通知が不要になった場合は、いつでもこのページで無効化できます。
              </p>
            </div>
          </div>
          
          {/* ナビゲーション */}
          <div className="max-w-md mx-auto text-center">
            <Link
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-800 transition-colors"
            >
              ← ホームページに戻る
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}