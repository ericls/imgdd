import React from "react";
import { Anchor } from "~src/__generated__/graphql";
import { Button } from "~src/ui/button";
import { Input, InputWithLabel } from "~src/ui/input";
import { SelectWithLabel } from "~src/ui/select";
import { ImageGallery } from "~src/common/ImageGallery/render";
import { useAuth } from "~src/lib/auth";
import { RenderingImageItem } from "~src/common/ImageGallery/types";
import classNames from "classnames";
import { TEXT_COLOR, SECONDARY_TEXT_COLOR_DIM } from "~src/ui/classNames";
import { useTranslation } from "react-i18next";

export type WatermarkSettings = {
  overlayImageId: string;
  overlayImageUrl: string;
  opacity: number;
  scale: number;
  anchor: Anchor;
  positionX: number;
  positionY: number;
};

type WatermarkToolProps = {
  baseImageId: string;
  settings: WatermarkSettings;
  onSettingsChange: (settings: WatermarkSettings) => void;
  onApply: () => void;
  applying: boolean;
};

const ANCHOR_KEYS: Record<Anchor, string> = {
  [Anchor.TopLeft]: "watermarkTool.anchorOption.topLeft",
  [Anchor.TopRight]: "watermarkTool.anchorOption.topRight",
  [Anchor.BottomLeft]: "watermarkTool.anchorOption.bottomLeft",
  [Anchor.BottomRight]: "watermarkTool.anchorOption.bottomRight",
  [Anchor.Center]: "watermarkTool.anchorOption.center",
};

export function WatermarkTool({
  baseImageId,
  settings,
  onSettingsChange,
  onApply,
  applying,
}: WatermarkToolProps) {
  const { t } = useTranslation();
  const [showPicker, setShowPicker] = React.useState(false);
  const { data: authData } = useAuth();
  const userId = authData?.viewer.organizationUser?.id;

  const handleSelectOverlay = React.useCallback(
    (image: RenderingImageItem) => {
      if (image.id === baseImageId) return;
      onSettingsChange({
        ...settings,
        overlayImageId: image.id,
        overlayImageUrl: image.url,
      });
      setShowPicker(false);
    },
    [baseImageId, settings, onSettingsChange],
  );

  const overlayItemRenderer = React.useCallback(
    (image: RenderingImageItem) => {
      const isBase = image.id === baseImageId;
      const isSelected = image.id === settings.overlayImageId;
      return (
        <button
          type="button"
          disabled={isBase}
          onClick={() => handleSelectOverlay(image)}
          className={classNames(
            "rounded-md overflow-hidden cursor-pointer transition-all",
            {
              "opacity-30 cursor-not-allowed": isBase,
              "ring-2 ring-indigo-500": isSelected,
            },
          )}
        >
          <div className="relative w-full pb-[80%] overflow-hidden bg-transparent rounded-md">
            <img
              src={image.url}
              alt={image.name}
              className="absolute top-0 left-0 w-full h-full object-cover"
            />
          </div>
          <div className="p-1 text-xs truncate">{image.name}</div>
        </button>
      );
    },
    [baseImageId, settings.overlayImageId, handleSelectOverlay],
  );

  return (
    <div className="space-y-4">
      <div>
        <label className={classNames("block mb-1 font-medium", TEXT_COLOR)}>
          {t("watermarkTool.overlayImage")}
        </label>
        {settings.overlayImageUrl ? (
          <div className="flex items-center gap-3">
            <img
              src={settings.overlayImageUrl}
              alt="overlay"
              className="w-16 h-16 object-cover rounded-md"
            />
            <Button
              variant="secondary"
              onClick={() => setShowPicker(!showPicker)}
            >
              {t("watermarkTool.change")}
            </Button>
          </div>
        ) : (
          <Button
            variant="secondary"
            onClick={() => setShowPicker(!showPicker)}
          >
            {t("watermarkTool.selectOverlay")}
          </Button>
        )}
      </div>

      {showPicker && userId && (
        <OverlayPicker userId={userId} itemRenderer={overlayItemRenderer} />
      )}

      {settings.overlayImageId && (
        <>
          <InputWithLabel
            label={t("watermarkTool.opacity", {
              value: Math.round(settings.opacity * 100),
            })}
            type="range"
            min="0"
            max="1"
            step="0.01"
            value={settings.opacity}
            onChange={(e) =>
              onSettingsChange({
                ...settings,
                opacity: parseFloat(e.target.value),
              })
            }
          />

          <InputWithLabel
            label={t("watermarkTool.scale", {
              value: Math.round(settings.scale * 100),
            })}
            type="range"
            min="0.01"
            max="1"
            step="0.01"
            value={settings.scale}
            onChange={(e) =>
              onSettingsChange({
                ...settings,
                scale: parseFloat(e.target.value),
              })
            }
          />

          <SelectWithLabel
            label={t("watermarkTool.anchor")}
            value={settings.anchor}
            onChange={(e) =>
              onSettingsChange({
                ...settings,
                anchor: e.target.value as Anchor,
              })
            }
          >
            {Object.entries(ANCHOR_KEYS).map(([value, key]) => (
              <option key={value} value={value}>
                {t(key)}
              </option>
            ))}
          </SelectWithLabel>

          <Button
            onClick={onApply}
            disabled={applying || !settings.overlayImageId}
            className="w-full"
          >
            {applying
              ? t("watermarkTool.applying")
              : t("watermarkTool.applyWatermark")}
          </Button>
        </>
      )}
    </div>
  );
}

function OverlayPicker({
  userId,
  itemRenderer,
}: {
  userId: string;
  itemRenderer: (image: RenderingImageItem) => React.ReactNode;
}) {
  const { t } = useTranslation();
  const [search, setSearch] = React.useState("");
  const [debouncedSearch, setDebouncedSearch] = React.useState("");

  React.useEffect(() => {
    const timer = setTimeout(() => setDebouncedSearch(search), 300);
    return () => clearTimeout(timer);
  }, [search]);

  return (
    <div className="border border-gray-300 dark:border-neutral-700 rounded-md p-3 max-h-96 overflow-y-auto">
      <Input
        type="text"
        placeholder={t("watermarkTool.searchPlaceholder")}
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="mb-2 w-full"
      />
      <p className={classNames("text-xs mb-2", SECONDARY_TEXT_COLOR_DIM)}>
        {t("watermarkTool.selectHint")}
      </p>
      <ImageGallery
        createdById={userId}
        nameContains={debouncedSearch || undefined}
        itemRenderer={itemRenderer}
      />
    </div>
  );
}
