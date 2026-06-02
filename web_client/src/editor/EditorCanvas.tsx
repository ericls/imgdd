import React from "react";
import { Anchor } from "~src/__generated__/graphql";
import { BlurRegion } from "./BlurTool";

export type OverlayState = {
  image: HTMLImageElement | null;
  x: number; // 0-1 normalized
  y: number; // 0-1 normalized
  opacity: number; // 0-1
  scale: number; // 0-1, relative to base short side
  anchor: Anchor;
};

type EditorCanvasProps = {
  baseImageUrl: string;
  overlay: OverlayState;
  onPositionChange: (x: number, y: number) => void;
  blurRegion?: BlurRegion | null;
  blurRadius?: number;
  onBlurRegionChange?: (region: BlurRegion | null) => void;
  mode?: "watermark" | "blur";
  className?: string;
};

export function EditorCanvas({
  baseImageUrl,
  overlay,
  onPositionChange,
  blurRegion,
  blurRadius = 10,
  onBlurRegionChange,
  mode = "watermark",
  className,
}: EditorCanvasProps) {
  const canvasRef = React.useRef<HTMLCanvasElement>(null);
  const baseImgRef = React.useRef<HTMLImageElement | null>(null);
  const [baseLoaded, setBaseLoaded] = React.useState(false);
  const [dragging, setDragging] = React.useState(false);
  const dragStartRef = React.useRef<{ x: number; y: number } | null>(null);

  // Load base image
  React.useEffect(() => {
    let cancelled = false;
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => {
      if (!cancelled) {
        baseImgRef.current = img;
        setBaseLoaded(true);
      }
    };
    img.src = baseImageUrl;
    return () => {
      cancelled = true;
      setBaseLoaded(false);
    };
  }, [baseImageUrl]);

  // Draw canvas
  React.useEffect(() => {
    const canvas = canvasRef.current;
    const baseImg = baseImgRef.current;
    if (!canvas || !baseImg || !baseLoaded) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    const container = canvas.parentElement;
    if (!container) return;
    const dpr = window.devicePixelRatio || 1;
    const maxCssW = container.clientWidth;
    const ratio = baseImg.naturalWidth / baseImg.naturalHeight;
    const cssW = Math.min(maxCssW, Math.round(baseImg.naturalWidth / dpr));
    const cssH = Math.round(cssW / ratio);
    const bufferW = Math.round(cssW * dpr);
    const bufferH = Math.round(cssH * dpr);

    canvas.width = bufferW;
    canvas.height = bufferH;
    canvas.style.width = cssW + "px";
    canvas.style.height = cssH + "px";

    ctx.clearRect(0, 0, bufferW, bufferH);
    ctx.drawImage(baseImg, 0, 0, bufferW, bufferH);

    if (mode === "watermark" && overlay.image) {
      const shortSide = Math.min(bufferW, bufferH);
      const targetSize = shortSide * overlay.scale;
      const overlayRatio =
        overlay.image.naturalWidth / overlay.image.naturalHeight;
      let ow: number, oh: number;
      if (overlayRatio >= 1) {
        ow = targetSize;
        oh = targetSize / overlayRatio;
      } else {
        oh = targetSize;
        ow = targetSize * overlayRatio;
      }

      const px = overlay.x * bufferW;
      const py = overlay.y * bufferH;

      let ox: number, oy: number;
      switch (overlay.anchor) {
        case Anchor.TopLeft:
          ox = 0;
          oy = 0;
          break;
        case Anchor.TopRight:
          ox = -ow;
          oy = 0;
          break;
        case Anchor.BottomLeft:
          ox = 0;
          oy = -oh;
          break;
        case Anchor.BottomRight:
          ox = -ow;
          oy = -oh;
          break;
        case Anchor.Center:
        default:
          ox = -ow / 2;
          oy = -oh / 2;
          break;
      }

      ctx.globalAlpha = overlay.opacity;
      ctx.drawImage(overlay.image, px + ox, py + oy, ow, oh);
      ctx.globalAlpha = 1;
    }

    if (mode === "blur" && blurRegion) {
      const rx1 = blurRegion.x1 * bufferW;
      const ry1 = blurRegion.y1 * bufferH;
      const rw = (blurRegion.x2 - blurRegion.x1) * bufferW;
      const rh = (blurRegion.y2 - blurRegion.y1) * bufferH;
      if (rw > 0 && rh > 0) {
        // Scale the radius from full-res image pixels down to canvas buffer pixels
        const scaledRadius = Math.max(
          1,
          Math.round((blurRadius * bufferW) / baseImg.naturalWidth),
        );

        // Clip to region, apply CSS blur filter, redraw just that part of the
        // base image — this gives a pixel-accurate preview of the blur effect.
        ctx.save();
        ctx.beginPath();
        ctx.rect(rx1, ry1, rw, rh);
        ctx.clip();
        ctx.filter = `blur(${scaledRadius}px)`;
        ctx.drawImage(baseImg, 0, 0, bufferW, bufferH);
        ctx.restore();

        // Draw selection border on top
        ctx.save();
        ctx.strokeStyle = "rgba(99, 102, 241, 0.9)";
        ctx.lineWidth = 2 * dpr;
        ctx.setLineDash([6 * dpr, 3 * dpr]);
        ctx.strokeRect(rx1, ry1, rw, rh);
        ctx.restore();
      }
    }
  }, [baseLoaded, overlay, blurRegion, blurRadius, mode]);

  const getCanvasCoords = React.useCallback(
    (e: React.MouseEvent<HTMLCanvasElement> | React.Touch) => {
      const canvas = canvasRef.current;
      if (!canvas) return { x: 0, y: 0 };
      const rect = canvas.getBoundingClientRect();
      const x = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
      const y = Math.max(0, Math.min(1, (e.clientY - rect.top) / rect.height));
      return { x, y };
    },
    [],
  );

  const handleMouseDown = React.useCallback(
    (e: React.MouseEvent<HTMLCanvasElement>) => {
      if (mode === "watermark") {
        if (!overlay.image) return;
        setDragging(true);
        const { x, y } = getCanvasCoords(e);
        onPositionChange(x, y);
      } else {
        setDragging(true);
        const { x, y } = getCanvasCoords(e);
        dragStartRef.current = { x, y };
        onBlurRegionChange?.(null);
      }
    },
    [
      mode,
      overlay.image,
      getCanvasCoords,
      onPositionChange,
      onBlurRegionChange,
    ],
  );

  const handleMouseMove = React.useCallback(
    (e: React.MouseEvent<HTMLCanvasElement>) => {
      if (!dragging) return;
      const { x, y } = getCanvasCoords(e);
      if (mode === "watermark") {
        onPositionChange(x, y);
      } else if (dragStartRef.current) {
        const start = dragStartRef.current;
        onBlurRegionChange?.({
          x1: Math.min(start.x, x),
          y1: Math.min(start.y, y),
          x2: Math.max(start.x, x),
          y2: Math.max(start.y, y),
        });
      }
    },
    [dragging, mode, getCanvasCoords, onPositionChange, onBlurRegionChange],
  );

  const handleMouseUp = React.useCallback(() => {
    setDragging(false);
    dragStartRef.current = null;
  }, []);

  const handleTouchStart = React.useCallback(
    (e: React.TouchEvent<HTMLCanvasElement>) => {
      if (e.touches.length === 0) return;
      if (mode === "watermark" && !overlay.image) return;
      e.preventDefault();
      setDragging(true);
      const { x, y } = getCanvasCoords(e.touches[0]);
      if (mode === "watermark") {
        onPositionChange(x, y);
      } else {
        dragStartRef.current = { x, y };
        onBlurRegionChange?.(null);
      }
    },
    [
      mode,
      overlay.image,
      getCanvasCoords,
      onPositionChange,
      onBlurRegionChange,
    ],
  );

  const handleTouchMove = React.useCallback(
    (e: React.TouchEvent<HTMLCanvasElement>) => {
      if (!dragging || e.touches.length === 0) return;
      e.preventDefault();
      const { x, y } = getCanvasCoords(e.touches[0]);
      if (mode === "watermark") {
        onPositionChange(x, y);
      } else if (dragStartRef.current) {
        const start = dragStartRef.current;
        onBlurRegionChange?.({
          x1: Math.min(start.x, x),
          y1: Math.min(start.y, y),
          x2: Math.max(start.x, x),
          y2: Math.max(start.y, y),
        });
      }
    },
    [dragging, mode, getCanvasCoords, onPositionChange, onBlurRegionChange],
  );

  const cursor =
    mode === "blur" ? "crosshair" : overlay.image ? "crosshair" : "default";

  return (
    <canvas
      ref={canvasRef}
      className={className}
      style={{ cursor }}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleMouseUp}
    />
  );
}
