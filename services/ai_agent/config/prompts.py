"""
プロンプトテンプレート
LLMへの指示を集中管理
"""

# クイズ生成プロンプト
QUIZ_GENERATION_PROMPT = """
あなたは1942年の生命保険会社の担当者です。
プレイヤーの回答を元に、記憶のシンクロを確認するクイズを1問生成してください。

## プレイヤーの回答
{player_answers}

## 場所
{place_name}

## 過去プレイヤーの回答（参考）
{past_answers}

## 指示
以下の形式でクイズを1問生成してください：

1. **質問文**: 担当者との記憶のシンクロを確認する質問（プレイヤーの回答に基づく）
2. **正解**: プレイヤーの回答から1つ選択
3. **ダミー1**: プレイヤーの別の回答から選択
4. **ダミー2**: 過去プレイヤーの回答から選択（提供された場合）
5. **ダミー3**: システム汎用回答

## 制約
- 質問文は明確で曖昧さがないこと
- 選択肢は互いに区別可能であること
- 1942年の生命保険診査という設定を維持すること

## 出力形式
JSON形式で以下のように出力してください：
```json
{{
  "question_text": "質問文",
  "options": [
    {{"text": "選択肢1", "is_correct": true, "source": "player_answer"}},
    {{"text": "選択肢2", "is_correct": false, "source": "player_other_answer"}},
    {{"text": "選択肢3", "is_correct": false, "source": "past_player"}},
    {{"text": "選択肢4", "is_correct": false, "source": "system_generic"}}
  ],
  "answer_index": 0
}}
```
"""

# 回答バリデーションプロンプト
RESPONSE_VALIDATION_PROMPT = """
以下の質問に対するプレイヤーの回答が有効かどうかを判定してください。

## 質問
{question}

## プレイヤーの回答
{answer}

## 判定基準
- 「なし」「特にない」「わからない」「思いつかない」などの無効回答パターンを検出
- 質問内容に対して適切な回答かどうかを評価
- 具体性があるかどうかを確認

## 出力形式
JSON形式で以下のように出力してください：
```json
{{
  "is_valid": true/false,
  "reason": "判定理由",
  "feedback": "無効の場合のフィードバックメッセージ（有効な場合はnull）",
  "confidence": 0.0-1.0
}}
```

## フィードバックメッセージの例
- "もう少し具体的に教えてください。小学生の頃に楽しかったことや、よくやっていたことを思い出してみてください。"
- "どんな些細なことでも構いません。当時のあなたが心に残っていることを教えてください。"
"""

# 対話生成プロンプト
DIALOGUE_GENERATION_PROMPT = """
あなたは1942年の生命保険会社の担当者です。
プレイヤーの情報に応じた自然な対話テキストを生成してください。

## シーン
{scene}

## プレイヤー情報
{player_context}

## 指示
- 1942年の生命保険診査という設定を維持すること
- 自然な話し言葉で生成すること（VOICEVOX用）
- 1つの発話は100文字以内に収めること
- ホラー要素を適度に含めること（不快にならない程度）
- プレイヤーの入力に応じてパーソナライズすること

## 出力形式
JSON形式で以下のように出力してください：
```json
{{
  "dialogue_text": "生成された対話テキスト",
  "voice_text": "VOICEVOX用テキスト（読みやすいように調整）",
  "estimated_duration_ms": 推定音声長（ミリ秒）
}}
```
"""

# メッセージ検証プロンプト
MESSAGE_VERIFICATION_PROMPT = """
以下の2つのメッセージが一致するかどうかを判定してください。

## 元のメッセージ（S4で入力）
{original}

## 再入力されたメッセージ（S7で入力）
{reinput}

## 判定基準
- 完全一致の場合は即座に正解
- 文字が若干異なる場合でも、意味が同等であれば正解
- 表記揺れ（ひらがな・カタカナ、送り仮名など）は許容
- 意味が明らかに異なる場合は不一致

## 出力形式
JSON形式で以下のように出力してください：
```json
{{
  "matched": true/false,
  "similarity_score": 0.0-1.0,
  "reason": "判定理由",
  "hint": "不一致の場合のヒント（一致の場合はnull）"
}}
```

## ヒントの例
- "最初の部分は合っています。後半をもう一度思い出してみてください。"
- "言葉の順番が少し違います。もう一度正確に思い出してみてください。"
"""

# フォールバッククイズテンプレート
FALLBACK_QUIZZES = {
    "spiral_stairs": {
        "question": "この螺旋階段を見上げたとき、あなたは何を感じましたか？",
        "options": [
            {"text": "懐かしさ", "is_correct": True},
            {"text": "不安", "is_correct": False},
            {"text": "期待", "is_correct": False},
            {"text": "静けさ", "is_correct": False}
        ]
    },
    "fireplace": {
        "question": "メインホールの暖炉のレンガを見て、何を思いましたか？",
        "options": [
            {"text": "温もり", "is_correct": True},
            {"text": "歴史", "is_correct": False},
            {"text": "孤独", "is_correct": False},
            {"text": "記憶", "is_correct": False}
        ]
    },
    "hinge": {
        "question": "裏玄関の扉の蝶番に、どんな印象を持ちましたか？",
        "options": [
            {"text": "古びた美しさ", "is_correct": True},
            {"text": "錆びた悲しさ", "is_correct": False},
            {"text": "時の流れ", "is_correct": False},
            {"text": "堅牢さ", "is_correct": False}
        ]
    },
    "entrance": {
        "question": "入口エントランスの扉を見たとき、何を感じましたか？",
        "options": [
            {"text": "歓迎", "is_correct": True},
            {"text": "緊張", "is_correct": False},
            {"text": "好奇心", "is_correct": False},
            {"text": "畏怖", "is_correct": False}
        ]
    },
    "piano": {
        "question": "階上応接室のピアノに、どんな印象を持ちましたか？",
        "options": [
            {"text": "優雅さ", "is_correct": True},
            {"text": "寂しさ", "is_correct": False},
            {"text": "懐かしさ", "is_correct": False},
            {"text": "荘厳さ", "is_correct": False}
        ]
    }
}

# フォールバック対話テンプレート
FALLBACK_DIALOGUES = {
    "s1_greeting": "ようこそ。お待ちしておりました。",
    "s1_purpose": "そうですか。それは素晴らしいですね。",
    "s1_activity": "なるほど、承知いたしました。",
}
