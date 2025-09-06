import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <main className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center">
          <div className="mb-8">
            <h1 className="text-5xl font-bold text-gray-900 mb-4">
              Kotti He Oide App
            </h1>
            <p className="text-xl text-gray-600 mb-8">
              Web Push通知とカメラフィルター機能を備えたPWAアプリ
            </p>
          </div>

          <div className="grid md:grid-cols-2 gap-8 mb-12">
            <div className="bg-white rounded-xl p-6 shadow-lg hover:shadow-xl transition-shadow">
              <div className="text-4xl mb-4">🔔</div>
              <h2 className="text-2xl font-semibold mb-3">Web Push通知</h2>
              <p className="text-gray-600 mb-4">
                ブラウザを閉じていても重要な通知を受け取ることができます。リアルタイムで情報をお届けします。
              </p>
              <Link 
                href="/notifications" 
                className="inline-block bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors"
              >
                通知設定
              </Link>
            </div>

            <div className="bg-white rounded-xl p-6 shadow-lg hover:shadow-xl transition-shadow">
              <div className="text-4xl mb-4">📷</div>
              <h2 className="text-2xl font-semibold mb-3">カメラフィルター</h2>
              <p className="text-gray-600 mb-4">
                5種類のリアルタイムカメラフィルター（レトロ、ホラー、シリアス、VHS、コミック）で写真を楽しく撮影。
              </p>
              <Link 
                href="/camera-filters" 
                className="inline-block bg-green-600 text-white px-6 py-2 rounded-lg hover:bg-green-700 transition-colors"
              >
                フィルターを試す
              </Link>
            </div>
          </div>

          <div className="bg-white/80 backdrop-blur rounded-xl p-6">
            <h3 className="text-2xl font-semibold mb-4">PWA機能</h3>
            <div className="grid md:grid-cols-3 gap-4 text-sm">
              <div className="flex items-center space-x-2">
                <span className="text-green-500">✓</span>
                <span>オフライン対応</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-green-500">✓</span>
                <span>ホーム画面に追加</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-green-500">✓</span>
                <span>プッシュ通知</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-green-500">✓</span>
                <span>バックグラウンド同期</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-green-500">✓</span>
                <span>アプリライクな体験</span>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-green-500">✓</span>
                <span>高速起動</span>
              </div>
            </div>
          </div>
        </div>
      </main>

      <footer className="container mx-auto px-4 py-8 text-center text-gray-600">
        <div className="flex justify-center space-x-8 text-sm">
          <Link href="/notifications" className="hover:text-blue-600 transition-colors">
            通知設定
          </Link>
          <Link href="/camera-filters" className="hover:text-blue-600 transition-colors">
            カメラフィルター
          </Link>
        </div>
        <p className="mt-4 text-xs">
          © 2024 Kotti He Oide App. Web Push通知対応PWAアプリケーション
        </p>
      </footer>
    </div>
  );
}
