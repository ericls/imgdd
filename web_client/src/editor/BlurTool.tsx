import React from "react";
import { Button } from "~src/ui/button";
import { InputWithLabel } from "~src/ui/input";
import classNames from "classnames";
import { TEXT_COLOR, SECONDARY_TEXT_COLOR_DIM } from "~src/ui/classNames";
import { useTranslation } from "react-i18next";

export type BlurRegion = {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
};

export type BlurSettings = {
  region: BlurRegion | null;
  radius: number;
};

type BlurToolProps = {
  settings: BlurSettings;
  onSettingsChange: (settings: BlurSettings) => void;
  onApply: () => void;
  applying: boolean;
};

export function BlurTool({
  settings,
  onSettingsChange,
  onApply,
  applying,
}: BlurToolProps) {
  const { t } = useTranslation();
  const hasRegion =
    settings.region !== null &&
    settings.region.x2 > settings.region.x1 &&
    settings.region.y2 > settings.region.y1;

  return (
    <div className="space-y-4">
      <p className={classNames("text-sm", SECONDARY_TEXT_COLOR_DIM)}>
        {t("blurTool.hint")}
      </p>

      {hasRegion && settings.region ? (
        <div className={classNames("text-sm", TEXT_COLOR)}>
          {t("blurTool.regionSelected", {
            x1: Math.round(settings.region.x1 * 100),
            y1: Math.round(settings.region.y1 * 100),
            x2: Math.round(settings.region.x2 * 100),
            y2: Math.round(settings.region.y2 * 100),
          })}
        </div>
      ) : (
        <div className={classNames("text-sm", SECONDARY_TEXT_COLOR_DIM)}>
          {t("blurTool.noRegion")}
        </div>
      )}

      <InputWithLabel
        label={t("blurTool.radius", { value: settings.radius })}
        type="range"
        min="1"
        max="50"
        step="1"
        value={settings.radius}
        onChange={(e) =>
          onSettingsChange({ ...settings, radius: parseInt(e.target.value) })
        }
      />

      <Button
        onClick={onApply}
        disabled={applying || !hasRegion}
        className="w-full"
      >
        {applying ? t("blurTool.applying") : t("blurTool.applyBlur")}
      </Button>
    </div>
  );
}
