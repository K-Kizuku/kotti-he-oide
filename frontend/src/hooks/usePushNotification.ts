import { useState, useEffect, useCallback } from 'react';
import { pushAPI } from '@/lib/api/pushApi';
import { registerServiceWorker } from '@/lib/utils/serviceWorker';
import { arrayBufferToBase64, urlBase64ToUint8Array } from '@/lib/utils/notificationUtils';
import { storage } from '@/lib/utils/storage';
import { PushSubscriptionState } from '@/types/push';

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
      const subscriptionId = subscription ? storage.getPushSubscriptionId() : null;

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

      // Convert VAPID key to the proper format
      const applicationServerKey = urlBase64ToUint8Array(vapidPublicKey);

      // Subscribe to push manager
      const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: applicationServerKey
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
        // Save subscription ID to storage
        storage.setPushSubscriptionId(response.id);
        storage.setVapidPublicKey(vapidPublicKey);
        
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

      // Clear storage
      storage.removePushSubscriptionId();
      storage.removeVapidPublicKey();

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

  const refresh = useCallback(async () => {
    await initializePushSubscription();
  }, []);

  return {
    ...state,
    subscribe,
    unsubscribe,
    refresh
  };
}