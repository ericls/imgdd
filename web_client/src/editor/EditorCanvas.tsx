import React from "react";
import { Anchor } from "~src/__generated__/graphql";

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
  className?: string;
};

export function EditorCanvas({
  baseImageUrl,
  overlay,
  onPositionChange,
  className,
}: EditorCanvasProps) {
  const canvasRef = React.useRef<HTMLCanvasElement>(null);
  const baseImgRef = React.useRef<HTMLImageElement | null>(null);
  const [baseLoaded, setBaseLoaded] = React.useState(false);
  const [dragging, setDragging] = React.useState(false);

  // Load base image
  React.useEffect(() => {
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => {
      baseImgRef.current = img;
      setBaseLoaded(true);
    };
    img.src = baseImageUrl;
  }, [baseImageUrl]);

  // Draw canvas
  React.useEffect(() => {
    const canvas = canvasRef.current;
    const baseImg = baseImgRef.current;
    if (!canvas || !baseImg || !baseLoaded) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    // Size canvas to fit container while maintaining aspect ratio
    const container = canvas.parentElement;
    if (!container) return;
    const maxW = container.clientWidth;
    const ratio = baseImg.naturalWidth / baseImg.naturalHeight;
    const displayW = Math.min(maxW, baseImg.naturalWidth);
    const displayH = displayW / ratio;

    canvas.width = displayW;
    canvas.height = displayH;

    // Draw base
    ctx.clearRect(0, 0, displayW, displayH);
    ctx.drawImage(baseImg, 0, 0, displayW, displayH);

    // Draw overlay
    if (overlay.image) {
      const shortSide = Math.min(displayW, displayH);
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

      const px = overlay.x * displayW;
      const py = overlay.y * displayH;

      // Apply anchor offset to match backend behavior
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
  }, [baseLoaded, overlay]);

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
      if (!overlay.image) return;
      setDragging(true);
      const { x, y } = getCanvasCoords(e);
      onPositionChange(x, y);
    },
    [overlay.image, getCanvasCoords, onPositionChange],
  );

  const handleMouseMove = React.useCallback(
    (e: React.MouseEvent<HTMLCanvasElement>) => {
      if (!dragging) return;
      const { x, y } = getCanvasCoords(e);
      onPositionChange(x, y);
    },
    [dragging, getCanvasCoords, onPositionChange],
  );

  const handleMouseUp = React.useCallback(() => {
    setDragging(false);
  }, []);

  const handleTouchStart = React.useCallback(
    (e: React.TouchEvent<HTMLCanvasElement>) => {
      if (!overlay.image || e.touches.length === 0) return;
      e.preventDefault();
      setDragging(true);
      const { x, y } = getCanvasCoords(e.touches[0]);
      onPositionChange(x, y);
    },
    [overlay.image, getCanvasCoords, onPositionChange],
  );

  const handleTouchMove = React.useCallback(
    (e: React.TouchEvent<HTMLCanvasElement>) => {
      if (!dragging || e.touches.length === 0) return;
      e.preventDefault();
      const { x, y } = getCanvasCoords(e.touches[0]);
      onPositionChange(x, y);
    },
    [dragging, getCanvasCoords, onPositionChange],
  );

  return (
    <canvas
      ref={canvasRef}
      className={className}
      style={{ cursor: overlay.image ? "crosshair" : "default" }}
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
