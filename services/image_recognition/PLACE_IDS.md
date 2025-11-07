# place_id 語彙（共有定義）

本サービスおよびフロント/Goサーバ間で共通に利用する `place_id` の一覧です。S3の参照画像は `REFERENCE_S3_PREFIX/{place_id}/...` 配下に配置してください。

| place_id | 説明 (日本語) | 推奨S3パス例 |
|---|---|---|
| `spiral_stairs` | 螺旋階段を見上げる高い天井 | `<prefix>/spiral_stairs/*.jpg` |
| `main_hall_fireplace` | メインホールの暖炉のレンガ | `<prefix>/main_hall_fireplace/*.jpg` |
| `back_entrance_hinge` | 裏玄関の扉の蝶番 | `<prefix>/back_entrance_hinge/*.jpg` |
| `entrance_door` | 入口エントランスの扉 | `<prefix>/entrance_door/*.jpg` |
| `upstairs_parlor_piano` | 階上応接室のピアノ | `<prefix>/upstairs_parlor_piano/*.jpg` |

備考:
- 参照カテゴリはS3キーの先頭セグメントで判定します（例: `spiral_stairs/ref1.jpg` → `place_id=spiral_stairs`）。
- 画像は複数枚配置可能です。撮影条件の違い（角度/距離/明るさ）を含めた方が精度が安定します。

