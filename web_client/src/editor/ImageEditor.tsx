import React from "react";
import { useParams, useNavigate } from "react-router";
import { useLazyQuery } from "@apollo/client/react";
import { gql } from "~src/__generated__";
import { Anchor } from "~src/__generated__/graphql";
import { EditorCanvas, OverlayState } from "./EditorCanvas";
import { WatermarkTool, WatermarkSettings } from "./WatermarkTool";
import { BlurTool, BlurSettings, BlurRegion } from "./BlurTool";
import { useApplyWatermark, useApplyBlur } from "./data";
import { toast } from "react-toastify";
import { useTranslation } from "react-i18next";
import classNames from "classnames";
import {
  HEADING_2,
  SECOND_LAYER,
  TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIM,
} from "~src/ui/classNames";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";
import { Button } from "~src/ui/button";
import { absoluteURL } from "~src/lib/url";
import { routes } from "~src/routes";

const ImageForEditorDoc = gql(`
  query ImageForEditor($id: ID!) {
    viewer {
      id
      image(id: $id) {
        id
        url
        name
        identifier
        nominalWidth
        nominalHeight
        MIMEType
        parent {
          id
          name
        }
        changes
      }
    }
  }
`);

type EditorTab = "watermark" | "blur";

export function ImageEditor() {
  const { imageId } = useParams<{ imageId: string }>();
  const navigate = useNavigate();
  const [fetchImage, { data, loading, error }] =
    useLazyQuery(ImageForEditorDoc);
  const { t } = useTranslation();
  const { execute: applyWatermark, loading: applyingWatermark } =
    useApplyWatermark();
  const { execute: applyBlur, loading: applyingBlur } = useApplyBlur();

  const [activeTab, setActiveTab] = React.useState<EditorTab>("watermark");

  const [watermarkSettings, setWatermarkSettings] =
    React.useState<WatermarkSettings>({
      overlayImageId: "",
      overlayImageUrl: "",
      opacity: 0.5,
      scale: 0.25,
      anchor: Anchor.Center,
      positionX: 0.5,
      positionY: 0.5,
    });

  const [blurSettings, setBlurSettings] = React.useState<BlurSettings>({
    region: null,
    radius: 10,
  });

  const [overlayImg, setOverlayImg] = React.useState<HTMLImageElement | null>(
    null,
  );

  React.useEffect(() => {
    if (imageId) {
      fetchImage({ variables: { id: imageId } });
    }
  }, [imageId, fetchImage]);

  React.useEffect(() => {
    if (!watermarkSettings.overlayImageUrl) {
      return;
    }
    let cancelled = false;
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => {
      if (!cancelled) setOverlayImg(img);
    };
    img.src = absoluteURL(watermarkSettings.overlayImageUrl);
    return () => {
      cancelled = true;
      setOverlayImg(null);
    };
  }, [watermarkSettings.overlayImageUrl]);

  const image = data?.viewer.image;

  const overlay: OverlayState = React.useMemo(
    () => ({
      image: overlayImg,
      x: watermarkSettings.positionX,
      y: watermarkSettings.positionY,
      opacity: watermarkSettings.opacity,
      scale: watermarkSettings.scale,
      anchor: watermarkSettings.anchor,
    }),
    [
      overlayImg,
      watermarkSettings.positionX,
      watermarkSettings.positionY,
      watermarkSettings.opacity,
      watermarkSettings.scale,
      watermarkSettings.anchor,
    ],
  );

  const handlePositionChange = React.useCallback((x: number, y: number) => {
    setWatermarkSettings((prev) => ({ ...prev, positionX: x, positionY: y }));
  }, []);

  const handleBlurRegionChange = React.useCallback(
    (region: BlurRegion | null) => {
      setBlurSettings((prev) => ({ ...prev, region }));
    },
    [],
  );

  const handleApplyWatermark = React.useCallback(async () => {
    if (!imageId || !watermarkSettings.overlayImageId) return;
    try {
      const result = await applyWatermark({
        variables: {
          input: {
            baseImageId: imageId,
            overlayImageId: watermarkSettings.overlayImageId,
            position: {
              x: watermarkSettings.positionX,
              y: watermarkSettings.positionY,
            },
            anchor: watermarkSettings.anchor,
            opacity: watermarkSettings.opacity,
            scale: watermarkSettings.scale,
          },
        },
      });
      const newImage = result.data?.applyWatermark.image;
      if (newImage) {
        toast(t("imageEditor.watermarkApplied"));
        navigate(routes.profile.image(newImage.id), { replace: true });
      }
    } catch (_err) {
      toast.error(t("imageEditor.watermarkFailed"));
    }
  }, [t, imageId, watermarkSettings, applyWatermark, navigate]);

  const handleApplyBlur = React.useCallback(async () => {
    if (!imageId || !blurSettings.region) return;
    try {
      const result = await applyBlur({
        variables: {
          input: {
            baseImageId: imageId,
            region: {
              x1: blurSettings.region.x1,
              y1: blurSettings.region.y1,
              x2: blurSettings.region.x2,
              y2: blurSettings.region.y2,
            },
            radius: blurSettings.radius,
          },
        },
      });
      const newImage = result.data?.applyBlur.image;
      if (newImage) {
        toast(t("imageEditor.blurApplied"));
        navigate(routes.profile.image(newImage.id), { replace: true });
      }
    } catch (_err) {
      toast.error(t("imageEditor.blurFailed"));
    }
  }, [t, imageId, blurSettings, applyBlur, navigate]);

  if (loading) return <FullScreenLoader />;
  if (error)
    return (
      <div className={classNames("p-8", TEXT_COLOR)}>
        {t("imageEditor.errorLoading")}: {error.message}
      </div>
    );
  if (!image)
    return (
      <div className={classNames("p-8", TEXT_COLOR)}>
        {t("imageEditor.imageNotFound")}
      </div>
    );

  return (
    <div className="mx-8 my-4 max-w-full">
      <div className="flex items-center justify-between mb-4">
        <h1 className={classNames(HEADING_2, "font-poppins", "wrap-anywhere")}>
          {t("imageEditor.title", { name: image.name })}
        </h1>
        <Button variant="secondary" onClick={() => navigate(-1)}>
          {t("common.buttonLabel.back")}
        </Button>
      </div>

      {image.parent && (
        <p className={classNames("text-sm mb-3", SECONDARY_TEXT_COLOR_DIM)}>
          {t("imageEditor.derivedFrom", { name: image.parent.name })}
        </p>
      )}

      <div className="flex flex-col lg:flex-row gap-6">
        <div className="flex-1 min-w-0">
          <div
            className={classNames(
              "rounded-md p-2 overflow-hidden",
              SECOND_LAYER,
            )}
          >
            <EditorCanvas
              baseImageUrl={absoluteURL(image.url)}
              overlay={overlay}
              onPositionChange={handlePositionChange}
              blurRegion={blurSettings.region}
              blurRadius={blurSettings.radius}
              onBlurRegionChange={handleBlurRegionChange}
              mode={activeTab}
              className="max-w-full mx-auto block"
            />
          </div>
        </div>

        <div className="w-full lg:w-80 flex-shrink-0">
          <div className={classNames("rounded-md p-4", SECOND_LAYER)}>
            <div className="flex gap-2 mb-4">
              <button
                type="button"
                aria-pressed={activeTab === "watermark"}
                onClick={() => setActiveTab("watermark")}
                className={classNames(
                  "flex-1 py-1.5 rounded-md text-sm font-medium transition-colors",
                  activeTab === "watermark"
                    ? "bg-indigo-600 text-white"
                    : classNames(
                        "hover:bg-indigo-100 dark:hover:bg-neutral-700",
                        TEXT_COLOR,
                      ),
                )}
              >
                {t("imageEditor.watermark")}
              </button>
              <button
                type="button"
                aria-pressed={activeTab === "blur"}
                onClick={() => setActiveTab("blur")}
                className={classNames(
                  "flex-1 py-1.5 rounded-md text-sm font-medium transition-colors",
                  activeTab === "blur"
                    ? "bg-indigo-600 text-white"
                    : classNames(
                        "hover:bg-indigo-100 dark:hover:bg-neutral-700",
                        TEXT_COLOR,
                      ),
                )}
              >
                {t("imageEditor.blur")}
              </button>
            </div>

            {activeTab === "watermark" && (
              <WatermarkTool
                baseImageId={image.id}
                settings={watermarkSettings}
                onSettingsChange={setWatermarkSettings}
                onApply={handleApplyWatermark}
                applying={applyingWatermark}
              />
            )}

            {activeTab === "blur" && (
              <BlurTool
                settings={blurSettings}
                onSettingsChange={setBlurSettings}
                onApply={handleApplyBlur}
                applying={applyingBlur}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
