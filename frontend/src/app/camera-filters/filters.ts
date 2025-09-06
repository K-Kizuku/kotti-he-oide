/*
 * リアルタイム用の軽量フィルター群（Canvas 2D ImageData ベース）
 * - 画素処理は基本的に in-place（必要に応じて元画像をコピー）
 * - 目安: 640x480 / 24fps 程度を想定
 */

export type FilterId =
  | "retro"
  | "horror"
  | "serious"
  | "vhs"
  | "comic";

export type FilterSpec = {
  id: FilterId;
  label: string;
  description: string;
};

export const FILTERS: FilterSpec[] = [
  { id: "retro", label: "レトロ（セピア/ビネット/粒子）", description: "セピア + 周辺減光 + フィルム粒子" },
  { id: "horror", label: "ホラー（低彩度/緑被り/強コントラスト）", description: "寒色寄りの恐怖演出" },
  { id: "serious", label: "シリアス（ノワール調）", description: "低彩度 + 強コントラスト + 軽いビネット" },
  { id: "vhs", label: "VHS/グリッチ（色ずれ/走査線）", description: "色チャンネルのずれ + 走査線" },
  { id: "comic", label: "コミック（ポスタライズ/輪郭）", description: "色段階化 + Sobel輪郭" },
];

const clampByte = (v: number) => (v < 0 ? 0 : v > 255 ? 255 : v) | 0;

const toGray = (r: number, g: number, b: number) => 0.2126 * r + 0.7152 * g + 0.0722 * b;

function adjustContrast(v: number, factor: number) {
  return clampByte((v - 128) * factor + 128);
}

function applyVignette(data: Uint8ClampedArray, w: number, h: number, strength = 0.4) {
  const cx = w * 0.5;
  const cy = h * 0.5;
  const maxD = Math.sqrt(cx * cx + cy * cy);
  for (let y = 0; y < h; y++) {
    for (let x = 0; x < w; x++) {
      const i = (y * w + x) * 4;
      const d = Math.sqrt((x - cx) * (x - cx) + (y - cy) * (y - cy));
      const v = 1 - strength * (d / maxD) ** 2; // 中心=1, 周辺で減光
      data[i] = clampByte(data[i] * v);
      data[i + 1] = clampByte(data[i + 1] * v);
      data[i + 2] = clampByte(data[i + 2] * v);
    }
  }
}

function addScanlines(data: Uint8ClampedArray, w: number, h: number, darkness = 0.12) {
  for (let y = 0; y < h; y += 2) {
    const factor = 1 - darkness;
    for (let x = 0; x < w; x++) {
      const i = (y * w + x) * 4;
      data[i] = clampByte(data[i] * factor);
      data[i + 1] = clampByte(data[i + 1] * factor);
      data[i + 2] = clampByte(data[i + 2] * factor);
    }
  }
}

export function applyRetro(image: ImageData) {
  const { data, width: w, height: h } = image;
  // セピア調 + 粒子
  for (let i = 0; i < data.length; i += 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    let nr = 0.393 * r + 0.769 * g + 0.189 * b;
    let ng = 0.349 * r + 0.686 * g + 0.168 * b;
    let nb = 0.272 * r + 0.534 * g + 0.131 * b;
    // 粒子（ランダム微小ノイズ）
    const n = (Math.random() - 0.5) * 18; // ±9程度
    data[i] = clampByte(nr + n);
    data[i + 1] = clampByte(ng + n);
    data[i + 2] = clampByte(nb + n);
  }
  applyVignette(data, w, h, 0.45);
  addScanlines(data, w, h, 0.15);
}

export function applyHorror(image: ImageData) {
  const { data, width: w, height: h } = image;
  // 低彩度 + 強コントラスト + 緑被り + シャドウノイズ
  for (let i = 0; i < data.length; i += 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    const gray = toGray(r, g, b);
    // 低彩度へ寄せる
    let nr = r * 0.3 + gray * 0.7;
    let ng = g * 0.3 + gray * 0.7;
    let nb = b * 0.3 + gray * 0.7;
    // 緑被り（寒色寄り）
    nr *= 0.92;
    ng *= 1.08;
    nb *= 0.95;
    // コントラスト強調
    nr = adjustContrast(nr, 1.45);
    ng = adjustContrast(ng, 1.45);
    nb = adjustContrast(nb, 1.45);
    // シャドウにノイズ
    const v = (nr + ng + nb) / 3;
    const noise = v < 80 ? (Math.random() - 0.5) * 24 : 0;
    data[i] = clampByte(nr + noise);
    data[i + 1] = clampByte(ng + noise * 0.6);
    data[i + 2] = clampByte(nb + noise * 0.4);
  }
  applyVignette(image.data, image.width, image.height, 0.6);
  addScanlines(image.data, image.width, image.height, 0.1);
}

export function applySerious(image: ImageData) {
  const { data } = image;
  // ノワール: 低彩度 + しっかりコントラスト + 弱めの暗部締め
  for (let i = 0; i < data.length; i += 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    const gray = toGray(r, g, b);
    let nr = r * 0.4 + gray * 0.6;
    let ng = g * 0.4 + gray * 0.6;
    let nb = b * 0.4 + gray * 0.6;
    // コントラスト
    nr = adjustContrast(nr, 1.25);
    ng = adjustContrast(ng, 1.25);
    nb = adjustContrast(nb, 1.25);
    // ガンマでやや暗めに（>1 で暗く）
    const gamma = 1.08;
    nr = clampByte(255 * Math.pow(nr / 255, gamma));
    ng = clampByte(255 * Math.pow(ng / 255, gamma));
    nb = clampByte(255 * Math.pow(nb / 255, gamma));
    data[i] = nr;
    data[i + 1] = ng;
    data[i + 2] = nb;
  }
  applyVignette(data, image.width, image.height, 0.35);
}

export function applyVhs(image: ImageData) {
  const { data, width: w, height: h } = image;
  const src = new Uint8ClampedArray(data); // 元を保持
  const offset = 1; // R/B の水平シフト(px)
  // 走査線
  addScanlines(data, w, h, 0.18);
  // ランダムの横ノイズ行
  const glitchRow = Math.floor(Math.random() * h);
  for (let x = 0; x < w; x++) {
    const i = (glitchRow * w + x) * 4;
    data[i] = clampByte(data[i] * 1.3);
    data[i + 1] = clampByte(data[i + 1] * 1.1);
    data[i + 2] = clampByte(data[i + 2] * 1.4);
  }
  // 色ずれ（R 左、B 右）
  for (let y = 0; y < h; y++) {
    for (let x = 0; x < w; x++) {
      const i = (y * w + x) * 4;
      const rx = Math.max(0, x - offset);
      const bx = Math.min(w - 1, x + offset);
      const ri = (y * w + rx) * 4;
      const bi = (y * w + bx) * 4;
      data[i] = src[ri]; // R
      // G は元のまま
      data[i + 2] = src[bi + 2]; // B
    }
  }
  // 彩度わずかに低下
  for (let i = 0; i < data.length; i += 4) {
    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];
    const gray = toGray(r, g, b);
    data[i] = clampByte(r * 0.8 + gray * 0.2);
    data[i + 1] = clampByte(g * 0.8 + gray * 0.2);
    data[i + 2] = clampByte(b * 0.8 + gray * 0.2);
  }
}

export function applyComic(image: ImageData) {
  const { data, width: w, height: h } = image;
  const src = new Uint8ClampedArray(data); // 元を保持
  // ポスタライズ
  const levels = 6;
  const step = 255 / (levels - 1);
  for (let i = 0; i < data.length; i += 4) {
    data[i] = Math.round(data[i] / step) * step;
    data[i + 1] = Math.round(data[i + 1] / step) * step;
    data[i + 2] = Math.round(data[i + 2] / step) * step;
  }
  // Sobel エッジ検出をグレイスケールで実施し、閾値を超えた箇所を黒で描く
  const threshold = 100; // 0..~360 程度
  const gray = new Uint8ClampedArray(w * h);
  for (let y = 0; y < h; y++) {
    for (let x = 0; x < w; x++) {
      const i = (y * w + x) * 4;
      gray[y * w + x] = toGray(src[i], src[i + 1], src[i + 2]);
    }
  }
  const sobelX = [-1, 0, 1, -2, 0, 2, -1, 0, 1];
  const sobelY = [-1, -2, -1, 0, 0, 0, 1, 2, 1];
  for (let y = 1; y < h - 1; y++) {
    for (let x = 1; x < w - 1; x++) {
      let gx = 0;
      let gy = 0;
      let k = 0;
      for (let j = -1; j <= 1; j++) {
        for (let i = -1; i <= 1; i++) {
          const v = gray[(y + j) * w + (x + i)];
          gx += v * sobelX[k];
          gy += v * sobelY[k];
          k++;
        }
      }
      const mag = Math.sqrt(gx * gx + gy * gy);
      if (mag > threshold) {
        const oi = (y * w + x) * 4;
        data[oi] = 0;
        data[oi + 1] = 0;
        data[oi + 2] = 0;
      }
    }
  }
}

export function applyFilter(image: ImageData, id: FilterId) {
  switch (id) {
    case "retro":
      return applyRetro(image);
    case "horror":
      return applyHorror(image);
    case "serious":
      return applySerious(image);
    case "vhs":
      return applyVhs(image);
    case "comic":
      return applyComic(image);
    default:
      return;
  }
}

