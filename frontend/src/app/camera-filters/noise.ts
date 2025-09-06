/*
 * フィルターとは独立して適用できる「一部分のみ瞬間的に損傷」系ノイズ群。
 * - 旧テレビ/ビデオテープ風の局所的な破損を想定。
 * - strength: 0..100（頻度と規模をざっくり兼用）。
 */

export type NoiseId =
  | "dropout" // 細い水平ドロップアウト
  | "block" // ブロック欠損/ずれ
  | "tear" // 横ずれ（縦帯の時間軸ぶれ）
  | "snow" // 局所的な砂嵐（白黒ノイズ）
  | "headswitch"; // 底部のヘッドスイッチングノイズ

export type NoiseSpec = {
  id: NoiseId;
  label: string;
  description: string;
};

export const NOISES: NoiseSpec[] = [
  { id: "dropout", label: "ドロップアウト（細線の欠損）", description: "水平の明滅する細い線" },
  { id: "block", label: "ブロック欠損/ずれ", description: "矩形領域の破損やコピーずれ" },
  { id: "tear", label: "ティア（帯の横ずれ）", description: "一部帯域の時間軸エラー風" },
  { id: "snow", label: "スノーノイズ（砂嵐）", description: "局所白黒ノイズのバースト" },
  { id: "headswitch", label: "ヘッドスイッチング（底部）", description: "下端のカラーバー/乱れ" },
];

const clampByte = (v: number) => (v < 0 ? 0 : v > 255 ? 255 : v) | 0;

function randInt(min: number, max: number) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

// strength を 0..1 に正規化
const nrm = (s: number) => Math.max(0, Math.min(1, s / 100));

// 1) 細い水平のドロップアウトライン（明滅）
function noiseDropout(image: ImageData, strength: number) {
  const { data, width: w, height: h } = image;
  const s = nrm(strength);
  const chance = 0.25 + s * 0.55; // フレーム毎の発生確率
  if (Math.random() > chance) return;
  const lines = 1 + Math.floor(s * 3);
  for (let n = 0; n < lines; n++) {
    const y = randInt(0, h - 1);
    const thick = randInt(1, Math.max(1, Math.floor(1 + s * 3)));
    const mode = Math.random();
    for (let t = 0; t < thick; t++) {
      const yy = Math.min(h - 1, y + t);
      for (let x = 0; x < w; x++) {
        const i = (yy * w + x) * 4;
        if (mode < 0.33) {
          // ほぼ白飛び
          data[i] = 240 + randInt(-8, 8);
          data[i + 1] = 240 + randInt(-8, 8);
          data[i + 2] = 240 + randInt(-8, 8);
        } else if (mode < 0.66) {
          // 黒寄り
          data[i] = randInt(0, 30);
          data[i + 1] = randInt(0, 30);
          data[i + 2] = randInt(0, 30);
        } else {
          // 近傍ラインをコピーして微ずれ
          const sy = Math.max(0, Math.min(h - 1, yy + (Math.random() < 0.5 ? -1 : 1)));
          const si = (sy * w + x) * 4;
          data[i] = data[si];
          data[i + 1] = data[si + 1];
          data[i + 2] = data[si + 2];
        }
      }
    }
  }
}

// 2) ブロック欠損/コピーずれ
function noiseBlock(image: ImageData, strength: number) {
  const { data, width: w, height: h } = image;
  const s = nrm(strength);
  const chance = 0.15 + s * 0.45;
  if (Math.random() > chance) return;
  const blocks = 1 + Math.floor(s * 3);
  const src = new Uint8ClampedArray(data);
  for (let b = 0; b < blocks; b++) {
    const bw = Math.max(8, Math.floor(w * (0.06 + s * 0.18)));
    const bh = Math.max(6, Math.floor(h * (0.04 + s * 0.14)));
    const dx = randInt(0, Math.max(0, w - bw));
    const dy = randInt(0, Math.max(0, h - bh));
    const sx = Math.max(0, Math.min(w - bw, dx + randInt(-Math.floor(20 + s * 80), Math.floor(20 + s * 80))));
    const sy = Math.max(0, Math.min(h - bh, dy + randInt(-Math.floor(10 + s * 60), Math.floor(10 + s * 60))));
    for (let y = 0; y < bh; y++) {
      for (let x = 0; x < bw; x++) {
        const di = ((dy + y) * w + (dx + x)) * 4;
        const si = ((sy + y) * w + (sx + x)) * 4;
        data[di] = src[si];
        data[di + 1] = src[si + 1];
        data[di + 2] = src[si + 2];
      }
    }
  }
}

// 3) 横ティア：水平帯の横方向ずれ（波形）
function noiseTear(image: ImageData, strength: number) {
  const { data, width: w, height: h } = image;
  const s = nrm(strength);
  const chance = 0.18 + s * 0.5;
  if (Math.random() > chance) return;
  const bandH = Math.max(6, Math.floor(h * (0.06 + s * 0.18)));
  const y0 = randInt(0, Math.max(0, h - bandH));
  const maxShift = Math.max(2, Math.floor(4 + s * 24));
  const src = new Uint8ClampedArray(data);
  for (let y = 0; y < bandH; y++) {
    const yy = y0 + y;
    const phase = (y / bandH) * Math.PI;
    const dx = Math.floor(Math.sin(phase) * maxShift * (Math.random() < 0.5 ? 1 : -1));
    for (let x = 0; x < w; x++) {
      const sx = Math.max(0, Math.min(w - 1, x + dx));
      const di = (yy * w + x) * 4;
      const si = (yy * w + sx) * 4;
      data[di] = src[si];
      data[di + 1] = src[si + 1];
      data[di + 2] = src[si + 2];
    }
  }
}

// 4) 局所スノーノイズ（砂嵐パッチ）
function noiseSnow(image: ImageData, strength: number) {
  const { data, width: w, height: h } = image;
  const s = nrm(strength);
  const chance = 0.25 + s * 0.6;
  if (Math.random() > chance) return;
  const patches = 1 + Math.floor(s * 4);
  for (let p = 0; p < patches; p++) {
    const pw = Math.max(6, Math.floor(w * (0.04 + s * 0.12)));
    const ph = Math.max(6, Math.floor(h * (0.04 + s * 0.12)));
    const x0 = randInt(0, Math.max(0, w - pw));
    const y0 = randInt(0, Math.max(0, h - ph));
    for (let y = 0; y < ph; y++) {
      for (let x = 0; x < pw; x++) {
        const i = ((y0 + y) * w + (x0 + x)) * 4;
        const v = randInt(0, 255);
        data[i] = v;
        data[i + 1] = v + randInt(-20, 20);
        data[i + 2] = v + randInt(-20, 20);
      }
    }
  }
}

// 5) 下端のヘッドスイッチングノイズ（色つきストライプ + 乱れ）
function noiseHeadSwitch(image: ImageData, strength: number) {
  const { data, width: w, height: h } = image;
  const s = nrm(strength);
  const chance = 0.15 + s * 0.5;
  if (Math.random() > chance) return;
  const bandH = Math.max(2, Math.floor(h * (0.02 + s * 0.06)));
  const y0 = h - bandH;
  for (let y = 0; y < bandH; y++) {
    const yy = y0 + y;
    const colorPhase = (y / bandH) * Math.PI * 2;
    const rBias = 128 + Math.floor(127 * Math.sin(colorPhase));
    const gBias = 128 + Math.floor(127 * Math.sin(colorPhase + 2.1));
    const bBias = 128 + Math.floor(127 * Math.sin(colorPhase + 4.2));
    for (let x = 0; x < w; x++) {
      const i = (yy * w + x) * 4;
      // 元の画素に強いカラーバイアス + 細かな乱れ
      data[i] = clampByte(rBias + randInt(-40, 40));
      data[i + 1] = clampByte(gBias + randInt(-40, 40));
      data[i + 2] = clampByte(bBias + randInt(-40, 40));
    }
  }
}

export function applyNoise(image: ImageData, id: NoiseId, strength: number) {
  switch (id) {
    case "dropout":
      return noiseDropout(image, strength);
    case "block":
      return noiseBlock(image, strength);
    case "tear":
      return noiseTear(image, strength);
    case "snow":
      return noiseSnow(image, strength);
    case "headswitch":
      return noiseHeadSwitch(image, strength);
  }
}

