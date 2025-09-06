import { WebPushSupportInfo } from '@/types/push';

export function checkWebPushSupport(): WebPushSupportInfo {
  const missingFeatures: string[] = [];

  if (!('serviceWorker' in navigator)) {
    missingFeatures.push('Service Worker');
  }

  if (!('PushManager' in window)) {
    missingFeatures.push('Push Manager');
  }

  if (!('Notification' in window)) {
    missingFeatures.push('Notifications API');
  }

  return {
    isSupported: missingFeatures.length === 0,
    missingFeatures
  };
}

export function arrayBufferToBase64(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
}

export function base64ToArrayBuffer(base64: string): ArrayBuffer {
  const binaryString = window.atob(base64);
  const len = binaryString.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes.buffer;
}

export function urlBase64ToUint8Array(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - base64String.length % 4) % 4);
  const base64 = (base64String + padding)
    .replace(/\-/g, '+')
    .replace(/_/g, '/');

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}

export function isNotificationSupported(): boolean {
  return 'Notification' in window;
}

export function getNotificationPermission(): NotificationPermission {
  if (!isNotificationSupported()) {
    return 'denied';
  }
  return Notification.permission;
}

export async function requestNotificationPermission(): Promise<NotificationPermission> {
  if (!isNotificationSupported()) {
    return 'denied';
  }

  if (Notification.permission === 'granted') {
    return 'granted';
  }

  const permission = await Notification.requestPermission();
  return permission;
}

export function isPushManagerSupported(): boolean {
  return 'serviceWorker' in navigator && 'PushManager' in window;
}

export function isServiceWorkerSupported(): boolean {
  return 'serviceWorker' in navigator;
}