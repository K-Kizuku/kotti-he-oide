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
const nrm01 = (v: number) => Math.max(0, Math.min(1, v));
// 旧API用: 0..100 -> 0..1
const nrm = (s: number) => Math.max(0, Math.min(1, s / 100));

// === 高度設定対応: エンジン方式 ===
export type NoiseParams = {
  enabled: boolean;
  frequency: number; // 0..1 (フレーム毎発生確率)
  magnitude: number; // 0..1 (影響の強さ/振幅)
  durationFrames: number; // 1..60 程度
  size: number; // 0..1 (幾何サイズ: 帯幅/ブロック/パッチ)
};

export type NoiseConfig = Record<NoiseId, NoiseParams>;

export const defaultNoiseConfig = (): NoiseConfig => ({
  dropout: { enabled: false, frequency: 0.2, magnitude: 0.6, durationFrames: 1, size: 0.25 },
  block: { enabled: false, frequency: 0.12, magnitude: 0.6, durationFrames: 2, size: 0.25 },
  tear: { enabled: false, frequency: 0.12, magnitude: 0.7, durationFrames: 3, size: 0.3 },
  snow: { enabled: false, frequency: 0.2, magnitude: 0.6, durationFrames: 1, size: 0.2 },
  headswitch: { enabled: false, frequency: 0.08, magnitude: 0.8, durationFrames: 2, size: 0.15 },
});

// イベント表現（各ノイズにつき同時に 0..1 個発生を想定）
type DropoutEvent = { y: number; thick: number; mode: number; ttl: number };
type BlockEvent = { dx: number; dy: number; bw: number; bh: number; sx: number; sy: number; ttl: number };
type TearEvent = { y0: number; bandH: number; maxShift: number; dir: 1 | -1; ttl: number };
type SnowEvent = { x0: number; y0: number; pw: number; ph: number; ttl: number };
type HeadSwitchEvent = { bandH: number; ttl: number };

type EngineState = {
  dropout?: DropoutEvent;
  block?: BlockEvent;
  tear?: TearEvent;
  snow?: SnowEvent;
  headswitch?: HeadSwitchEvent;
};

export function createNoiseEngine(init: NoiseConfig) {
  let config: NoiseConfig = JSON.parse(JSON.stringify(init));
  const state: EngineState = {};

  function spawn(width: number, height: number) {
    // ドロップアウト
    if (!state.dropout && config.dropout.enabled && Math.random() < nrm01(config.dropout.frequency)) {
      const thick = Math.max(1, Math.floor(1 + config.dropout.size * 10));
      state.dropout = {
        y: randInt(0, height - 1),
        thick,
        mode: Math.random(),
        ttl: Math.max(1, Math.floor(config.dropout.durationFrames)),
      };
    }
    // ブロック
    if (!state.block && config.block.enabled && Math.random() < nrm01(config.block.frequency)) {
      const bw = Math.max(6, Math.floor(width * (0.05 + config.block.size * 0.3)));
      const bh = Math.max(6, Math.floor(height * (0.04 + config.block.size * 0.25)));
      const dx = randInt(0, Math.max(0, width - bw));
      const dy = randInt(0, Math.max(0, height - bh));
      const shift = Math.floor((10 + config.block.magnitude * 90));
      const sx = Math.max(0, Math.min(width - bw, dx + randInt(-shift, shift)));
      const sy = Math.max(0, Math.min(height - bh, dy + randInt(-Math.floor(shift * 0.6), Math.floor(shift * 0.6))));
      state.block = { dx, dy, bw, bh, sx, sy, ttl: Math.max(1, Math.floor(config.block.durationFrames)) };
    }
    // ティア
    if (!state.tear && config.tear.enabled && Math.random() < nrm01(config.tear.frequency)) {
      const bandH = Math.max(6, Math.floor(height * (0.04 + config.tear.size * 0.25)));
      const y0 = randInt(0, Math.max(0, height - bandH));
      const maxShift = Math.max(2, Math.floor(4 + config.tear.magnitude * 28));
      const dir: 1 | -1 = Math.random() < 0.5 ? 1 : -1;
      state.tear = { y0, bandH, maxShift, dir, ttl: Math.max(1, Math.floor(config.tear.durationFrames)) };
    }
    // スノー
    if (!state.snow && config.snow.enabled && Math.random() < nrm01(config.snow.frequency)) {
      const pw = Math.max(6, Math.floor(width * (0.03 + config.snow.size * 0.2)));
      const ph = Math.max(6, Math.floor(height * (0.03 + config.snow.size * 0.2)));
      const x0 = randInt(0, Math.max(0, width - pw));
      const y0 = randInt(0, Math.max(0, height - ph));
      state.snow = { x0, y0, pw, ph, ttl: Math.max(1, Math.floor(config.snow.durationFrames)) };
    }
    // ヘッドスイッチング
    if (!state.headswitch && config.headswitch.enabled && Math.random() < nrm01(config.headswitch.frequency)) {
      const bandH = Math.max(2, Math.floor(height * (0.015 + config.headswitch.size * 0.12)));
      state.headswitch = { bandH, ttl: Math.max(1, Math.floor(config.headswitch.durationFrames)) };
    }
  }

  function applyEvents(image: ImageData) {
    const { data, width: w, height: h } = image;
    // ドロップアウト
    if (state.dropout && state.dropout.ttl > 0) {
      const { y, thick, mode } = state.dropout;
      for (let t = 0; t < thick; t++) {
        const yy = Math.min(h - 1, y + t);
        for (let x = 0; x < w; x++) {
          const i = (yy * w + x) * 4;
          if (mode < 0.33) {
            data[i] = 240 + randInt(-8, 8);
            data[i + 1] = 240 + randInt(-8, 8);
            data[i + 2] = 240 + randInt(-8, 8);
          } else if (mode < 0.66) {
            data[i] = randInt(0, 30);
            data[i + 1] = randInt(0, 30);
            data[i + 2] = randInt(0, 30);
          } else {
            const sy = Math.max(0, Math.min(h - 1, yy + (Math.random() < 0.5 ? -1 : 1)));
            const si = (sy * w + x) * 4;
            data[i] = data[si];
            data[i + 1] = data[si + 1];
            data[i + 2] = data[si + 2];
          }
        }
      }
      state.dropout.ttl -= 1;
      if (state.dropout.ttl <= 0) state.dropout = undefined;
    }
    // ブロック
    if (state.block && state.block.ttl > 0) {
      const { dx, dy, bw, bh, sx, sy } = state.block;
      const src = new Uint8ClampedArray(data);
      for (let y = 0; y < bh; y++) {
        for (let x = 0; x < bw; x++) {
          const di = ((dy + y) * w + (dx + x)) * 4;
          const si = ((sy + y) * w + (sx + x)) * 4;
          data[di] = src[si];
          data[di + 1] = src[si + 1];
          data[di + 2] = src[si + 2];
        }
      }
      state.block.ttl -= 1;
      if (state.block.ttl <= 0) state.block = undefined;
    }
    // ティア
    if (state.tear && state.tear.ttl > 0) {
      const { y0, bandH, maxShift, dir } = state.tear;
      const src = new Uint8ClampedArray(data);
      for (let y = 0; y < bandH; y++) {
        const yy = y0 + y;
        const phase = (y / bandH) * Math.PI;
        const dx = Math.floor(Math.sin(phase) * maxShift) * dir;
        for (let x = 0; x < w; x++) {
          const sx = Math.max(0, Math.min(w - 1, x + dx));
          const di = (yy * w + x) * 4;
          const si = (yy * w + sx) * 4;
          data[di] = src[si];
          data[di + 1] = src[si + 1];
          data[di + 2] = src[si + 2];
        }
      }
      state.tear.ttl -= 1;
      if (state.tear.ttl <= 0) state.tear = undefined;
    }
    // スノー
    if (state.snow && state.snow.ttl > 0) {
      const { x0, y0, pw, ph } = state.snow;
      for (let y = 0; y < ph; y++) {
        for (let x = 0; x < pw; x++) {
          const i = ((y0 + y) * w + (x0 + x)) * 4;
          const v = randInt(0, 255);
          data[i] = v;
          data[i + 1] = v + randInt(-20, 20);
          data[i + 2] = v + randInt(-20, 20);
        }
      }
      state.snow.ttl -= 1;
      if (state.snow.ttl <= 0) state.snow = undefined;
    }
    // ヘッドスイッチング
    if (state.headswitch && state.headswitch.ttl > 0) {
      const { bandH } = state.headswitch;
      const y0 = h - bandH;
      for (let y = 0; y < bandH; y++) {
        const yy = y0 + y;
        const colorPhase = (y / bandH) * Math.PI * 2;
        const amp = 127;
        const rBias = 128 + Math.floor(amp * Math.sin(colorPhase));
        const gBias = 128 + Math.floor(amp * Math.sin(colorPhase + 2.1));
        const bBias = 128 + Math.floor(amp * Math.sin(colorPhase + 4.2));
        for (let x = 0; x < w; x++) {
          const i = (yy * w + x) * 4;
          data[i] = clampByte(rBias + randInt(-40, 40));
          data[i + 1] = clampByte(gBias + randInt(-40, 40));
          data[i + 2] = clampByte(bBias + randInt(-40, 40));
        }
      }
      state.headswitch.ttl -= 1;
      if (state.headswitch.ttl <= 0) state.headswitch = undefined;
    }
  }

  return {
    updateConfig(next: NoiseConfig) {
      config = JSON.parse(JSON.stringify(next));
    },
    step(image: ImageData) {
      spawn(image.width, image.height);
      applyEvents(image);
    },
  };
}

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

// 既存の簡易API（後方互換）。単発・簡易ノイズの付加。
export function applyNoise(image: ImageData, id: NoiseId, strength: number) {
  const s = Math.max(0, Math.min(100, strength));
  const imageCopy = new ImageData(new Uint8ClampedArray(image.data), image.width, image.height);
  switch (id) {
    case "dropout":
      noiseDropout(imageCopy, s);
      break;
    case "block":
      noiseBlock(imageCopy, s);
      break;
    case "tear":
      noiseTear(imageCopy, s);
      break;
    case "snow":
      noiseSnow(imageCopy, s);
      break;
    case "headswitch":
      noiseHeadSwitch(imageCopy, s);
      break;
  }
  image.data.set(imageCopy.data);
}
