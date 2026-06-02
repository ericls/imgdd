import React from "react";
import { useParams, useNavigate } from "react-router";
import { useLazyQuery } from "@apollo/client/react";
import { gql } from "~src/__generated__";
import { Anchor } from "~src/__generated__/graphql";
import { EditorCanvas, OverlayState } from "./EditorCanvas";
import { WatermarkTool, WatermarkSettings } from "./WatermarkTool";
import { useApplyWatermark } from "./data";
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

export function ImageEditor() {
  const { imageId } = useParams<{ imageId: string }>();
  const navigate = useNavigate();
  const [fetchImage, { data, loading, error }] =
    useLazyQuery(ImageForEditorDoc);
  const { t } = useTranslation();
  const { execute: applyWatermark, loading: applying } = useApplyWatermark();

  const [settings, setSettings] = React.useState<WatermarkSettings>({
    overlayImageId: "",
    overlayImageUrl: "",
    opacity: 0.5,
    scale: 0.25,
    anchor: Anchor.Center,
    positionX: 0.5,
    positionY: 0.5,
  });

  const [overlayImg, setOverlayImg] = React.useState<HTMLImageElement | null>(
    null,
  );

  React.useEffect(() => {
    if (imageId) {
      fetchImage({ variables: { id: imageId } });
    }
  }, [imageId, fetchImage]);

  // Load overlay image element when URL changes
  React.useEffect(() => {
    if (!settings.overlayImageUrl) {
      return;
    }
    let cancelled = false;
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => {
      if (!cancelled) setOverlayImg(img);
    };
    img.src = absoluteURL(settings.overlayImageUrl);
    return () => {
      cancelled = true;
      setOverlayImg(null);
    };
  }, [settings.overlayImageUrl]);

  const image = data?.viewer.image;

  const overlay: OverlayState = React.useMemo(
    () => ({
      image: overlayImg,
      x: settings.positionX,
      y: settings.positionY,
      opacity: settings.opacity,
      scale: settings.scale,
      anchor: settings.anchor,
    }),
    [
      overlayImg,
      settings.positionX,
      settings.positionY,
      settings.opacity,
      settings.scale,
      settings.anchor,
    ],
  );

  const handlePositionChange = React.useCallback((x: number, y: number) => {
    setSettings((prev) => ({ ...prev, positionX: x, positionY: y }));
  }, []);

  const handleApply = React.useCallback(async () => {
    if (!imageId || !settings.overlayImageId) return;
    try {
      const result = await applyWatermark({
        variables: {
          input: {
            baseImageId: imageId,
            overlayImageId: settings.overlayImageId,
            position: {
              x: settings.positionX,
              y: settings.positionY,
            },
            anchor: settings.anchor,
            opacity: settings.opacity,
            scale: settings.scale,
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
  }, [t, imageId, settings, applyWatermark, navigate]);

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
              className="max-w-full mx-auto block"
            />
          </div>
        </div>

        <div className="w-full lg:w-80 flex-shrink-0">
          <div className={classNames("rounded-md p-4", SECOND_LAYER)}>
            <h2 className={classNames("text-lg font-medium mb-4", TEXT_COLOR)}>
              {t("imageEditor.watermark")}
            </h2>
            <WatermarkTool
              baseImageId={image.id}
              settings={settings}
              onSettingsChange={setSettings}
              onApply={handleApply}
              applying={applying}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
