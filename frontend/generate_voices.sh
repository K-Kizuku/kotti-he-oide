#!/bin/bash

# VOICEVOX音声生成スクリプト
# S1シーン用の固定音声を生成します

set -e

# 設定
VOICEVOX_HOST="http://127.0.0.1:50021"
SPEAKER_ID=84  # 青山龍星（しっとり）のID（環境に応じて変更してください）
OUTPUT_DIR="./public/audio/voice"

# 色付き出力用
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== VOICEVOX音声生成スクリプト ===${NC}"
echo ""

# 出力ディレクトリを作成
echo -e "${YELLOW}[1/5] ディレクトリを作成...${NC}"
mkdir -p "$OUTPUT_DIR"
echo -e "${GREEN}✓ ディレクトリ作成完了: $OUTPUT_DIR${NC}"
echo ""

# 音声生成関数
generate_voice() {
    local text="$1"
    local output_file="$2"
    local temp_query=$(mktemp)

    echo -e "${YELLOW}生成中: $output_file${NC}"
    echo "  テキスト: $text"

    # Step 1: audio_queryでクエリを生成
    echo "  → クエリ生成中..."
    curl -s -X POST \
        --get \
        --data-urlencode "text=$text" \
        --data-urlencode "speaker=$SPEAKER_ID" \
        "$VOICEVOX_HOST/audio_query" \
        -o "$temp_query"

    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ クエリ生成に失敗しました${NC}"
        rm -f "$temp_query"
        return 1
    fi

    # Step 2: synthesisで音声合成
    echo "  → 音声合成中..."
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -d @"$temp_query" \
        "$VOICEVOX_HOST/synthesis?speaker=$SPEAKER_ID" \
        -o "$OUTPUT_DIR/$output_file"

    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ 音声合成に失敗しました${NC}"
        rm -f "$temp_query"
        return 1
    fi

    # 一時ファイルを削除
    rm -f "$temp_query"

    # ファイルサイズをチェック
    local file_size=$(stat -f%z "$OUTPUT_DIR/$output_file" 2>/dev/null || stat -c%s "$OUTPUT_DIR/$output_file" 2>/dev/null)
    if [ "$file_size" -gt 1000 ]; then
        echo -e "${GREEN}✓ 生成完了: $output_file (${file_size} bytes)${NC}"
    else
        echo -e "${RED}✗ 生成に失敗した可能性があります (ファイルサイズが小さすぎます)${NC}"
        return 1
    fi
    echo ""
}

# VOICEVOXサーバーの疎通確認
echo -e "${YELLOW}[2/5] VOICEVOXサーバーの確認...${NC}"
if curl -s -f "$VOICEVOX_HOST/version" > /dev/null 2>&1; then
    VERSION=$(curl -s "$VOICEVOX_HOST/version")
    echo -e "${GREEN}✓ VOICEVOXサーバーに接続しました (version: $VERSION)${NC}"
else
    echo -e "${RED}✗ VOICEVOXサーバーに接続できません${NC}"
    echo "  $VOICEVOX_HOST が起動していることを確認してください"
    exit 1
fi
echo ""

# 各音声を生成
echo -e "${YELLOW}[3/5] welcome.wav を生成...${NC}"
generate_voice "ようこそ、赤煉瓦文化館へ。私は本日の担当者でございます。お越しいただき誠にありがとうございます。" "welcome.wav"

echo -e "${YELLOW}[4/5] ask_visit_method.wav を生成...${NC}"
generate_voice "本日はどのようにしてこちらまでいらっしゃいましたか？" "ask_visit_method.wav"

echo -e "${YELLOW}[5/5] ask_activities.wav を生成...${NC}"
generate_voice "なるほど。では、普段はどのようなことをされていますか？" "ask_activities.wav"

echo -e "${YELLOW}[6/6] closing.wav を生成...${NC}"
generate_voice "ありがとうございます。それでは、まずこの建物でお気に入りの場所を見つけていただきたいと思います。" "closing.wav"

# 完了
echo ""
echo -e "${GREEN}=== 全ての音声ファイルが生成されました ===${NC}"
echo ""
echo "生成されたファイル:"
ls -lh "$OUTPUT_DIR"/*.wav 2>/dev/null || echo "  (ファイルが見つかりませんでした)"
echo ""
echo -e "${YELLOW}注意: speaker ID ($SPEAKER_ID) が正しいか確認してください${NC}"
echo "  speaker IDの確認: curl -s \"$VOICEVOX_HOST/speakers\" | python3 -m json.tool | grep -B 2 -A 10 \"青山龍星\""
