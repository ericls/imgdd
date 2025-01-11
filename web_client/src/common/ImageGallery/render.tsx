import React from "react";
import { ImageItem, ImageItemRenderer } from "./types";
import { useImagesQuery } from "./data";
import { humanFileSize } from "~src/lib/humanizeFileSize";
import { useHumanizeDateTime } from "~src/lib/humanizeDateTime";
import classNames from "~node_modules/classnames";
import { SECONDARY_TEXT_COLOR_DIM, TEXT_COLOR } from "~src/ui/classNames";
import { Button } from "~src/ui/button";
import {
  GrFormNext as NextPageIcon,
  GrFormPrevious as PrevPageIcon,
} from "react-icons/gr";
import { Loader } from "~src/ui/loader";

type DumbImageGalleryProps = {
  images: ImageItem[];
  itemRenderer: ImageItemRenderer;
  hasNext: boolean;
  hasPrev: boolean;
  loadNextPage: () => void;
  loadPrevPage: () => void;
  loading?: boolean;
};

export function DumbImageGallery({
  images,
  itemRenderer,
  hasNext,
  hasPrev,
  loadNextPage,
  loadPrevPage,
  loading,
}: DumbImageGalleryProps) {
  return (
    <div>
      <div className="grid gap-x-6 gap-y-2 grid-cols-[repeat(auto-fill,minmax(280px,1fr))]">
        {images.map((image) => (
          <React.Fragment key={image.id}>{itemRenderer(image)}</React.Fragment>
        ))}
      </div>
      <div>
        <div className="flex justify-center mt-4 gap-x-4">
          <Button
            variant={hasPrev && !loading ? "indigo" : "secondary"}
            onClick={loadPrevPage}
            disabled={!hasPrev || loading}
          >
            <PrevPageIcon className="w-5 h-5" />
          </Button>
          <Button
            onClick={loadNextPage}
            variant={hasNext && !loading ? "indigo" : "secondary"}
            disabled={!hasNext || loading}
          >
            <NextPageIcon className="w-5 h-5" />
          </Button>
        </div>
      </div>
    </div>
  );
}

export function ImageItemRenderer({ image }: { image: ImageItem }) {
  const { url, name, nominalWidth, nominalHeight, nominalByteSize, createdAt } =
    image;

  const humanizedDateTime = useHumanizeDateTime({ datetimeStr: createdAt });

  return (
    <div className="flex flex-col overflow-hidden ">
      {/* Image container with fixed aspect ratio (square) */}
      <div className="relative w-full pb-[80%] overflow-hidden bg-transparent rounded-md">
        <img
          src={url}
          alt={`preview of image file: ${name}`}
          className="absolute top-0 left-0 w-full h-full object-cover"
        />
      </div>
      <div className="flex items-center p-2 space-x-3">
        <div className={classNames("flex flex-col justify-center max-w-full")}>
          <span
            className={classNames(
              "text-base font-medium truncate w-100",
              TEXT_COLOR
            )}
            title={name}
          >
            {name}
          </span>
          <span
            className={classNames(
              "text-sm text-gray-600",
              SECONDARY_TEXT_COLOR_DIM
            )}
            title={new Date(createdAt).toLocaleString()}
          >
            {humanizedDateTime}
          </span>
        </div>
      </div>

      <div className="p-2 text-sm text-gray-700 hidden">
        <div>
          <strong>Dimensions:</strong> {nominalWidth} x {nominalHeight} px
        </div>
        <div>
          <strong>File Size:</strong> {humanFileSize(nominalByteSize)}
        </div>
      </div>
    </div>
  );
}

type ImageGalleryProps = {
  nameContains?: string;
  itemRenderer?: ImageItemRenderer;
};
export function ImageGallery({
  nameContains,
  itemRenderer,
}: ImageGalleryProps) {
  const { data, execute, hasNext, hasPrev, goNext, goPrev, loading } =
    useImagesQuery({
      nameContains,
    });
  React.useEffect(() => {
    execute();
  }, [execute]);
  const renderImage = React.useCallback(
    (image: ImageItem) => {
      if (itemRenderer) {
        return itemRenderer(image);
      }
      return <ImageItemRenderer image={image} />;
    },
    [itemRenderer]
  );
  const images = data?.viewer.images.edges.map((edge) => edge.node) ?? [];
  return (
    <DumbImageGallery
      images={images}
      itemRenderer={renderImage}
      hasNext={hasNext}
      hasPrev={hasPrev}
      loadNextPage={goNext}
      loadPrevPage={goPrev}
      loading={loading}
    />
  );
}
