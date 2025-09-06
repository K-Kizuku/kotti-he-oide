# 1. フロントエンド Web Push 通知実装ガイド

## 1.1 目的と要件整理

### 目的
- Next.js 15.5.2 アプリケーションに Web Push 通知機能を統合する
- サーバーサイド API (`/api/push/*`) と連携して通知の受信・表示を行う
- PWA 対応も視野に入れたモダンなフロントエンド実装

### 必須要件
- **ブラウザ通知許可** の取得とプッシュ購読登録
- **Service Worker** によるバックグラウンド通知受信
- **通知表示** とクリックアクションの処理
- **購読管理** (有効化/無効化/更新)
- **HTTPS 必須** (開発環境は localhost で代用可能)

### 推奨要件
- TypeScript による型安全な実装
- React Hooks を活用した状態管理
- エラーハンドリングとユーザーフレンドリーな UX
- 通知設定画面の提供
- PWA マニフェストファイル対応

---

# 2. アーキテクチャ設計

## 2.1 ファイル構成

```
frontend/
├── public/
│   ├── sw.js                          # Service Worker (通知受信処理)
│   ├── manifest.json                  # PWA マニフェスト
│   └── icons/                         # 通知用アイコン
├── src/
│   ├── app/
│   │   ├── components/
│   │   │   ├── PushNotificationManager.tsx  # 通知管理メインコンポーネント
│   │   │   ├── NotificationPermissionButton.tsx  # 許可要求ボタン
│   │   │   ├── NotificationSettings.tsx     # 設定画面
│   │   │   └── NotificationStatus.tsx       # 現在の状態表示
│   │   ├── hooks/
│   │   │   ├── usePushNotification.ts       # 通知管理フック
│   │   │   ├── useNotificationPermission.ts # 許可状態管理
│   │   │   └── useServiceWorker.ts          # SW 管理フック
│   │   ├── utils/
│   │   │   ├── pushApi.ts                   # サーバー API 呼び出し
│   │   │   ├── notificationUtils.ts         # 通知ユーティリティ
│   │   │   └── storage.ts                   # ローカルストレージ管理
│   │   └── types/
│   │       └── push.ts                      # Push 関連型定義
│   └── middleware.ts                        # Next.js ミドルウェア (必要に応じて)
```

## 2.2 データフロー

1. **初期化フロー**
   - Service Worker 登録
   - VAPID 公開鍵取得
   - 通知許可状態確認

2. **購読フロー**
   - ユーザーの通知許可取得
   - Push Manager でブラウザ購読
   - サーバーへ購読情報送信 (`POST /api/push/subscribe`)

3. **通知受信フロー**
   - Service Worker が `push` イベント受信
   - 通知データパース・表示
   - クリック時のアクション実行

4. **購読管理フロー**
   - 購読状態の監視・更新
   - 購読解除 (`DELETE /api/push/subscriptions/{id}`)

---

# 3. 実装手順

## 3.1 Service Worker 実装

### 3.1.1 基本 Service Worker (`public/sw.js`)

```javascript
// Service Worker for Web Push Notifications
const CACHE_NAME = 'push-notifications-v1';

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

  if (event.data) {
    const payload = event.data.json();
    Object.assign(options, {
      body: payload.body || options.body,
      icon: payload.icon || options.icon,
      tag: payload.tag || options.tag,
      data: payload.data || options.data,
      actions: payload.actions || options.actions
    });
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
});
```

### 3.1.2 Service Worker 登録ユーティリティ

```typescript
// src/app/utils/serviceWorker.ts
export async function registerServiceWorker(): Promise<ServiceWorkerRegistration | null> {
  if (!('serviceWorker' in navigator)) {
    console.warn('Service Worker not supported');
    return null;
  }

  try {
    const registration = await navigator.serviceWorker.register('/sw.js', {
      scope: '/'
    });
    
    console.log('Service Worker registered:', registration);
    return registration;
  } catch (error) {
    console.error('Service Worker registration failed:', error);
    return null;
  }
}

export function unregisterServiceWorker(): Promise<boolean> {
  if (!('serviceWorker' in navigator)) {
    return Promise.resolve(false);
  }

  return navigator.serviceWorker.getRegistration()
    .then((registration) => {
      if (registration) {
        return registration.unregister();
      }
      return false;
    });
}
```

## 3.2 API 通信ユーティリティ

### 3.2.1 サーバー API 呼び出し (`src/app/utils/pushApi.ts`)

```typescript
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

class PushAPI {
  private baseURL: string;

  constructor(baseURL: string = 'http://localhost:8080') {
    this.baseURL = baseURL;
  }

  async getVAPIDPublicKey(): Promise<string> {
    const response = await fetch(`${this.baseURL}/api/push/vapid-public-key`);
    const data: VAPIDResponse = await response.json();
    
    if (!data.success) {
      throw new Error(data.message);
    }
    
    return data.publicKey;
  }

  async subscribe(subscriptionData: SubscribeRequest): Promise<SubscribeResponse> {
    const response = await fetch(`${this.baseURL}/api/push/subscribe`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(subscriptionData),
    });

    const data: SubscribeResponse = await response.json();
    return data;
  }

  async unsubscribe(subscriptionId: string): Promise<{success: boolean, message: string}> {
    const response = await fetch(`${this.baseURL}/api/push/subscriptions/${subscriptionId}`, {
      method: 'DELETE',
    });

    return response.json();
  }
}

export const pushAPI = new PushAPI();
```

## 3.3 React Hooks 実装

### 3.3.1 通知許可管理フック (`src/app/hooks/useNotificationPermission.ts`)

```typescript
import { useState, useEffect } from 'react';

export type NotificationPermissionState = 'default' | 'denied' | 'granted' | 'unsupported';

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
```

### 3.3.2 プッシュ通知管理フック (`src/app/hooks/usePushNotification.ts`)

```typescript
import { useState, useEffect, useCallback } from 'react';
import { pushAPI } from '../utils/pushApi';
import { registerServiceWorker } from '../utils/serviceWorker';

interface PushSubscriptionState {
  isSubscribed: boolean;
  subscription: PushSubscription | null;
  subscriptionId: string | null;
  error: string | null;
  isLoading: boolean;
}

export function usePushNotification() {
  const [state, setState] = useState<PushSubscriptionState>({
    isSubscribed: false,
    subscription: null,
    subscriptionId: null,
    error: null,
    isLoading: true
  });

  // Initialize: Check existing subscription
  useEffect(() => {
    initializePushSubscription();
  }, []);

  const initializePushSubscription = async () => {
    try {
      const registration = await registerServiceWorker();
      if (!registration) {
        throw new Error('Service Worker registration failed');
      }

      const subscription = await registration.pushManager.getSubscription();
      const subscriptionId = subscription ? localStorage.getItem('pushSubscriptionId') : null;

      setState(prev => ({
        ...prev,
        isSubscribed: !!subscription,
        subscription,
        subscriptionId,
        isLoading: false
      }));
    } catch (error) {
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Unknown error',
        isLoading: false
      }));
    }
  };

  const subscribe = useCallback(async (): Promise<boolean> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }));

    try {
      // Get VAPID public key from server
      const vapidPublicKey = await pushAPI.getVAPIDPublicKey();

      // Get service worker registration
      const registration = await navigator.serviceWorker.getRegistration();
      if (!registration) {
        throw new Error('Service Worker not registered');
      }

      // Subscribe to push manager
      const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: vapidPublicKey
      });

      // Send subscription to server
      const subscriptionData = {
        endpoint: subscription.endpoint,
        keys: {
          p256dh: arrayBufferToBase64(subscription.getKey('p256dh')!),
          auth: arrayBufferToBase64(subscription.getKey('auth')!)
        },
        ua: navigator.userAgent
      };

      const response = await pushAPI.subscribe(subscriptionData);
      
      if (response.success) {
        // Save subscription ID to localStorage
        localStorage.setItem('pushSubscriptionId', response.id);
        
        setState(prev => ({
          ...prev,
          isSubscribed: true,
          subscription,
          subscriptionId: response.id,
          isLoading: false
        }));
        
        return true;
      } else {
        throw new Error(response.message);
      }
    } catch (error) {
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Subscribe failed',
        isLoading: false
      }));
      return false;
    }
  }, []);

  const unsubscribe = useCallback(async (): Promise<boolean> => {
    if (!state.subscription || !state.subscriptionId) {
      return false;
    }

    setState(prev => ({ ...prev, isLoading: true, error: null }));

    try {
      // Unsubscribe from push manager
      await state.subscription.unsubscribe();

      // Notify server
      await pushAPI.unsubscribe(state.subscriptionId);

      // Clear localStorage
      localStorage.removeItem('pushSubscriptionId');

      setState(prev => ({
        ...prev,
        isSubscribed: false,
        subscription: null,
        subscriptionId: null,
        isLoading: false
      }));

      return true;
    } catch (error) {
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Unsubscribe failed',
        isLoading: false
      }));
      return false;
    }
  }, [state.subscription, state.subscriptionId]);

  return {
    ...state,
    subscribe,
    unsubscribe,
    refresh: initializePushSubscription
  };
}

// Helper function to convert ArrayBuffer to Base64
function arrayBufferToBase64(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary)
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, ''); // Remove padding
}
```

## 3.4 React コンポーネント実装

### 3.4.1 通知許可ボタン (`src/app/components/NotificationPermissionButton.tsx`)

```tsx
'use client';

import { useNotificationPermission } from '../hooks/useNotificationPermission';

export default function NotificationPermissionButton() {
  const { permission, isLoading, requestPermission, isSupported, isGranted, isDenied } = useNotificationPermission();

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
        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
      >
        {isLoading ? '確認中...' : '通知を許可する'}
      </button>
    </div>
  );
}
```

### 3.4.2 プッシュ通知管理コンポーネント (`src/app/components/PushNotificationManager.tsx`)

```tsx
'use client';

import { usePushNotification } from '../hooks/usePushNotification';
import { useNotificationPermission } from '../hooks/useNotificationPermission';
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
          <span className={`px-2 py-1 rounded text-sm ${
            isSubscribed ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
          }`}>
            {isSubscribed ? '有効' : '無効'}
          </span>
        </div>
        
        {subscription && (
          <div className="text-xs text-gray-500">
            <p>エンドポイント: {subscription.endpoint.substring(0, 50)}...</p>
          </div>
        )}
        
        <div className="flex space-x-2">
          {!isSubscribed ? (
            <button
              onClick={subscribe}
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
            >
              {isLoading ? '設定中...' : '通知を有効化'}
            </button>
          ) : (
            <button
              onClick={unsubscribe}
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50"
            >
              {isLoading ? '解除中...' : '通知を無効化'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
```

## 3.5 PWA 対応

### 3.5.1 マニフェストファイル (`public/manifest.json`)

```json
{
  "name": "Kotti He Oide App",
  "short_name": "KottiApp",
  "description": "Web Push通知対応アプリケーション",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#000000",
  "icons": [
    {
      "src": "/icons/icon-192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-512.png", 
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

### 3.5.2 アイコンファイル準備

```bash
# public/icons/ に以下のファイルを配置
icon-192.png    # 192x192px の通知アイコン
icon-512.png    # 512x512px の PWA アイコン
badge-72.png    # 72x72px の通知バッジアイコン
```

### 3.5.3 Layout への PWA メタタグ追加

```tsx
// src/app/layout.tsx に追加
export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <head>
        <link rel="manifest" href="/manifest.json" />
        <meta name="theme-color" content="#000000" />
        <link rel="icon" href="/icons/icon-192.png" />
        <link rel="apple-touch-icon" href="/icons/icon-192.png" />
      </head>
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        {children}
      </body>
    </html>
  );
}
```

---

# 4. 開発・テスト手順

## 4.1 開発環境セットアップ

```bash
cd frontend

# 依存関係の追加（必要に応じて）
pnpm add @types/web

# 開発サーバー起動
pnpm dev
```

## 4.2 テストシナリオ

### 4.2.1 基本フロー
1. ブラウザで https://localhost:3000 にアクセス
2. 通知許可を求められた際に「許可」を選択
3. 「通知を有効化」ボタンをクリック
4. サーバー側から `POST /api/push/send` でテスト通知送信
5. 通知が表示されることを確認
6. 通知をクリックしてアプリが開くことを確認

### 4.2.2 エラーケース
- 通知拒否時の表示確認
- ネットワークエラー時の挙動確認
- Service Worker 登録失敗時の処理確認

## 4.3 ブラウザ別対応確認

- **Chrome/Edge**: 完全サポート
- **Firefox**: 完全サポート  
- **Safari**: iOS 16.4+ / macOS 13+ で対応
- **Safari (旧バージョン)**: 非対応表示

---

# 5. 本番環境対応

## 5.1 HTTPS 必須対応
- 本番環境では必ず HTTPS での配信が必要
- Let's Encrypt 等での SSL 証明書設定

## 5.2 セキュリティ対応
- VAPID キーの適切な管理
- CSP (Content Security Policy) 設定
- Service Worker のスコープ制限

## 5.3 パフォーマンス最適化
- Service Worker のキャッシュ戦略
- 通知表示の最適化
- バックグラウンド同期対応（将来的）

---

# 6. 運用・監視

## 6.1 メトリクス収集
- 通知許可率の測定
- 購読数の監視
- 通知クリック率の計測

## 6.2 エラー監視
- Service Worker エラーの監視
- API 呼び出し失敗の記録
- ユーザーエージェント別の対応状況確認

---

# 7. 今後の拡張案

## 7.1 機能拡張
- 通知設定のカスタマイズ（トピック別 ON/OFF）
- 通知履歴表示
- 通知スケジューリング機能
- リッチ通知（画像・アクション付き）

## 7.2 UX 改善  
- 通知のプレビュー機能
- 設定画面の充実
- オンボーディングフローの改善

---

この実装により、サーバーサイドと連携したフル機能のWeb Push通知システムが完成します。Next.js 15.5.2の最新機能を活用し、TypeScriptによる型安全性を確保した実装となっています。