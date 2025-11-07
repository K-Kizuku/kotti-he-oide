/**
 * ブラウザの戻る/進むボタンを完全にブロックするユーティリティ
 * ゲーム進行中は正しい順序でシーン遷移を強制するため、履歴操作を無効化します
 */

type PopStateHandler = (event: PopStateEvent) => void;

class BrowserHistoryBlocker {
  private isBlocking = false;
  private popStateHandler: PopStateHandler | null = null;

  /**
   * ブラウザバックのブロックを有効化
   */
  enable(): void {
    if (this.isBlocking) {
      return;
    }

    // 現在の状態を履歴にプッシュ（これにより「戻る」を無効化）
    window.history.pushState(null, '', window.location.href);

    // popstateイベントをリッスンし、戻る操作を阻止
    this.popStateHandler = (event: PopStateEvent) => {
      event.preventDefault();
      // 即座に同じ状態を再度プッシュして、履歴スタックを維持
      window.history.pushState(null, '', window.location.href);
    };

    window.addEventListener('popstate', this.popStateHandler);
    this.isBlocking = true;

    console.debug('[BrowserHistoryBlocker] 有効化: ブラウザバックがブロックされました');
  }

  /**
   * ブラウザバックのブロックを無効化
   */
  disable(): void {
    if (!this.isBlocking || !this.popStateHandler) {
      return;
    }

    window.removeEventListener('popstate', this.popStateHandler);
    this.popStateHandler = null;
    this.isBlocking = false;

    console.debug('[BrowserHistoryBlocker] 無効化: ブラウザバックのブロックが解除されました');
  }

  /**
   * 現在ブロック中かどうかを返す
   */
  isEnabled(): boolean {
    return this.isBlocking;
  }
}

// シングルトンインスタンス
const browserHistoryBlocker = new BrowserHistoryBlocker();

export default browserHistoryBlocker;
