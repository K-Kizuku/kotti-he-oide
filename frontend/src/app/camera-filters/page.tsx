"use client";

import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";
import styles from "./page.module.css";
import { FILTERS, FilterId, applyFilter } from "./filters";
import { NOISES, createNoiseEngine, defaultNoiseConfig, type NoiseConfig } from "./noise";

type StreamState = "idle" | "starting" | "running" | "stopped" | "error";

const MAX_WIDTH = 640; // パフォーマンスと画質のバランス
const TARGET_FPS = 24;

export default function CameraFiltersPage() {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const rafRef = useRef<number | null>(null);
  const lastFrameRef = useRef<number>(0);

  const [state, setState] = useState<StreamState>("idle");
  const [error, setError] = useState<string>("");
  const [filter, setFilter] = useState<FilterId>("retro");
  const filterRef = useRef<FilterId>("retro");
  const [facingMode, setFacingMode] = useState<"user" | "environment">("user");
  // 複数ノイズの詳細設定
  const [noiseConfig, setNoiseConfig] = useState<NoiseConfig>(() => defaultNoiseConfig());

  const filterOptions = useMemo(() => FILTERS, []);
  const noiseOptions = useMemo(() => NOISES, []);
  const noiseEngineRef = useRef<ReturnType<typeof createNoiseEngine> | null>(null);

  const stopCamera = useCallback(() => {
    if (rafRef.current) {
      cancelAnimationFrame(rafRef.current);
      rafRef.current = null;
    }
    const stream = streamRef.current;
    if (stream) {
      stream.getTracks().forEach((t) => t.stop());
      streamRef.current = null;
    }
    setState("stopped");
  }, []);

  useEffect(() => {
    return () => stopCamera();
  }, [stopCamera]);

  const startCamera = useCallback(async () => {
    setError("");
    setState("starting");
    try {
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        throw new Error(
          "このブラウザは getUserMedia に対応していません。最新のブラウザをご利用ください。"
        );
      }
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode },
        audio: false,
      });
      streamRef.current = stream;
      const v = videoRef.current!;
      v.srcObject = stream;
      if (v.readyState < 1 || v.videoWidth === 0) {
        await new Promise<void>((resolve) => {
          const onMeta = () => {
            v.removeEventListener("loadedmetadata", onMeta);
            resolve();
          };
          v.addEventListener("loadedmetadata", onMeta);
        });
      }
      await v.play();

      // サイズ調整
      const vw = v.videoWidth || 640;
      const vh = v.videoHeight || 480;
      const scale = Math.min(1, MAX_WIDTH / vw);
      const cw = Math.round(vw * scale);
      const ch = Math.round(vh * scale);
      const c = canvasRef.current!;
      c.width = cw;
      c.height = ch;

      // 描画ループ
      const ctx = c.getContext("2d", { willReadFrequently: true });
      if (!ctx) throw new Error("Canvas 2D コンテキストを取得できませんでした");

      setState("running");
      lastFrameRef.current = performance.now();
      // ノイズエンジンを作成
      noiseEngineRef.current = createNoiseEngine(noiseConfig);

      const loop = () => {
        const now = performance.now();
        const elapsed = now - lastFrameRef.current;
        const interval = 1000 / TARGET_FPS;
        if (elapsed >= interval) {
          lastFrameRef.current = now - (elapsed % interval);
          // ミラー無しでそのまま描画
          ctx.drawImage(v, 0, 0, cw, ch);
          const frame = ctx.getImageData(0, 0, cw, ch);
          applyFilter(frame, filterRef.current);
          // 複数ノイズ適用（発生確率/持続/サイズ等はエンジンが管理）
          noiseEngineRef.current?.step(frame);
          ctx.putImageData(frame, 0, 0);
        }
        rafRef.current = requestAnimationFrame(loop);
      };
      rafRef.current = requestAnimationFrame(loop);
    } catch (e: unknown) {
      // 例外を安全にメッセージへ変換
      const message = e instanceof Error ? e.message : String(e);
      // ネイティブの例外メッセージが空の場合にフォールバック
      setError(
        message ||
          "カメラの起動に失敗しました。ブラウザの権限設定をご確認ください。"
      );
      setState("error");
    }
  }, [facingMode, noiseConfig]);

  const onChangeFilter = useCallback((e: React.ChangeEvent<HTMLSelectElement>) => {
    const id = e.target.value as FilterId;
    setFilter(id);
    filterRef.current = id;
  }, []);

  useEffect(() => {
    filterRef.current = filter;
  }, [filter]);

  // ノイズ設定の変更をエンジンへ反映
  useEffect(() => {
    noiseEngineRef.current?.updateConfig(noiseConfig);
  }, [noiseConfig]);

  const onToggle = useCallback(() => {
    if (state === "running" || state === "starting") {
      stopCamera();
    } else {
      startCamera();
    }
  }, [startCamera, stopCamera, state]);

  const isBusy = state === "starting";
  const isRunning = state === "running";

  return (
    <div className={styles.container}>
      <div className={styles.controls}>
        <h1 style={{ marginRight: 8, fontSize: 18 }}>カメラ・フィルターデモ</h1>
        <button
          className={`${styles.button} ${styles.primary}`}
          onClick={onToggle}
          disabled={isBusy}
        >
          {isRunning || isBusy ? "停止" : "カメラ開始"}
        </button>
        <select
          className={styles.select}
          value={filter}
          onChange={onChangeFilter}
          disabled={!isRunning && !isBusy}
          aria-label="フィルター"
        >
          {filterOptions.map((f) => (
            <option key={f.id} value={f.id}>
              {f.label}
            </option>
          ))}
        </select>
        <div className={styles.noisePanel} aria-label="ノイズ設定">
          {noiseOptions.map((n) => {
            const cfg = noiseConfig[n.id];
            return (
              <div key={n.id} className={styles.noiseRow} title={n.description}>
                <label className={styles.checkbox}>
                  <input
                    type="checkbox"
                    checked={cfg.enabled}
                    onChange={(e) =>
                      setNoiseConfig((prev) => ({
                        ...prev,
                        [n.id]: { ...prev[n.id], enabled: e.target.checked },
                      }))
                    }
                    disabled={!isRunning && !isBusy}
                  />
                  {n.label}
                </label>
                <label className={styles.noiseCtl}>
                  頻度
                  <input
                    className={styles.range}
                    type="range"
                    min={0}
                    max={100}
                    value={Math.round(cfg.frequency * 100)}
                    onChange={(e) =>
                      setNoiseConfig((prev) => ({
                        ...prev,
                        [n.id]: { ...prev[n.id], frequency: Number(e.target.value) / 100 },
                      }))
                    }
                    disabled={!cfg.enabled || (!isRunning && !isBusy)}
                  />
                </label>
                <label className={styles.noiseCtl}>
                  規模
                  <input
                    className={styles.range}
                    type="range"
                    min={0}
                    max={100}
                    value={Math.round(cfg.magnitude * 100)}
                    onChange={(e) =>
                      setNoiseConfig((prev) => ({
                        ...prev,
                        [n.id]: { ...prev[n.id], magnitude: Number(e.target.value) / 100 },
                      }))
                    }
                    disabled={!cfg.enabled || (!isRunning && !isBusy)}
                  />
                </label>
                <label className={styles.noiseCtl}>
                  持続f
                  <input
                    className={styles.range}
                    type="range"
                    min={1}
                    max={30}
                    value={cfg.durationFrames}
                    onChange={(e) =>
                      setNoiseConfig((prev) => ({
                        ...prev,
                        [n.id]: { ...prev[n.id], durationFrames: Number(e.target.value) },
                      }))
                    }
                    disabled={!cfg.enabled || (!isRunning && !isBusy)}
                  />
                </label>
                <label className={styles.noiseCtl}>
                  サイズ
                  <input
                    className={styles.range}
                    type="range"
                    min={0}
                    max={100}
                    value={Math.round(cfg.size * 100)}
                    onChange={(e) =>
                      setNoiseConfig((prev) => ({
                        ...prev,
                        [n.id]: { ...prev[n.id], size: Number(e.target.value) / 100 },
                      }))
                    }
                    disabled={!cfg.enabled || (!isRunning && !isBusy)}
                  />
                </label>
              </div>
            );
          })}
        </div>
        <select
          className={styles.select}
          value={facingMode}
          onChange={(e) =>
            setFacingMode(e.target.value as "user" | "environment")
          }
          disabled={isRunning || isBusy}
          aria-label="カメラ"
        >
          <option value="user">フロント</option>
          <option value="environment">バック</option>
        </select>
        <div className={styles.spacer} />
        <div className={styles.status}>
          状態: {state} {isBusy ? "（起動中）" : ""}
        </div>
      </div>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.grid}>
        <section className={styles.panel}>
          <div className={styles.title}>元映像（カメラ）</div>
          <div className={styles.mediaWrap}>
            <video
              ref={videoRef}
              className={styles.video}
              playsInline
              muted
              aria-hidden="true"
            />
          </div>
          <div className={styles.hint}>許可ダイアログが表示されたら「許可」を押してください。</div>
        </section>

        <section className={styles.panel}>
          <div className={styles.title}>フィルター適用後（{filterOptions.find((f) => f.id === filter)?.label}）</div>
          <div className={styles.mediaWrap}>
            <canvas
              ref={canvasRef}
              className={styles.canvas}
              role="img"
              aria-label="フィルター適用後の出力"
            />
          </div>
          <div className={styles.hint}>
            目安: 最大幅 {MAX_WIDTH}px / {TARGET_FPS}fps（環境により変動）
          </div>
        </section>
      </div>
    </div>
  );
}
