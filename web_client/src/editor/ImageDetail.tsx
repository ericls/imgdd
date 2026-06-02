import React from "react";
import { useParams, useNavigate, Link } from "react-router";
import { useLazyQuery } from "@apollo/client/react";
import { gql } from "~src/__generated__";
import classNames from "classnames";
import {
  HEADING_2,
  SECOND_LAYER,
  TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIM,
  BASE_LAYER,
} from "~src/ui/classNames";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";
import { Button } from "~src/ui/button";
import { absoluteURL } from "~src/lib/url";
import { routes } from "~src/routes";
import { copyText } from "~src/lib/copyText";
import { useTranslation } from "react-i18next";

const ImageDetailDoc = gql(`
  query ImageDetail($id: ID!) {
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
        createdAt
        changes
        lineage {
          id
          url
          name
          changes
          createdAt
        }
      }
    }
  }
`);

export function ImageDetail() {
  const { imageId } = useParams<{ imageId: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [fetchImage, { data, loading, error }] = useLazyQuery(ImageDetailDoc);

  React.useEffect(() => {
    if (imageId) {
      fetchImage({ variables: { id: imageId } });
    }
  }, [imageId, fetchImage]);

  const image = data?.viewer.image;

  if (loading) return <FullScreenLoader />;
  if (error)
    return (
      <div className={classNames("p-8", TEXT_COLOR)}>
        {t("imageDetail.errorLoading")}: {error.message}
      </div>
    );
  if (!image)
    return (
      <div className={classNames("p-8", TEXT_COLOR)}>
        {t("imageDetail.imageNotFound")}
      </div>
    );

  const lineage = image.lineage;
  const currentIndex = lineage.findIndex((img) => img.id === image.id);

  return (
    <div className="mx-8 my-4 max-w-full">
      <div className="flex items-center justify-between mb-4 gap-4">
        <h1 className={classNames(HEADING_2, "font-poppins min-w-0 break-all")}>
          {image.name}
        </h1>
        <div className="flex gap-2 flex-shrink-0">
          <Button
            variant="secondary"
            onClick={() => navigate(routes.profile.editImage(image.id))}
          >
            {t("common.buttonLabel.edit")}
          </Button>
          <Button variant="secondary" onClick={() => navigate(-1)}>
            {t("common.buttonLabel.back")}
          </Button>
        </div>
      </div>

      <div className="flex flex-col lg:flex-row gap-6">
        <div className="flex-1 min-w-0">
          <div
            className={classNames(
              "rounded-md p-2 overflow-hidden",
              SECOND_LAYER,
            )}
          >
            <img
              src={absoluteURL(image.url)}
              alt={image.name}
              className="max-w-full h-auto rounded mx-auto block"
              style={{
                width: Math.round(
                  image.nominalWidth / (window.devicePixelRatio || 1),
                ),
              }}
            />
          </div>
          <div className={classNames("mt-2 text-sm", SECONDARY_TEXT_COLOR_DIM)}>
            {image.nominalWidth} x {image.nominalHeight} px &middot;{" "}
            {image.MIMEType}
          </div>
          <CopyURLs url={absoluteURL(image.url)} name={image.name} />
        </div>

        {lineage.length > 1 && (
          <div className="w-full lg:w-80 flex-shrink-0">
            <div className={classNames("rounded-md p-4", SECOND_LAYER)}>
              <h2
                className={classNames("text-lg font-medium mb-4", TEXT_COLOR)}
              >
                {t("imageDetail.editHistory")}
              </h2>
              <div className="space-y-3">
                {lineage.map((ancestor, i) => {
                  const isCurrent = i === currentIndex;
                  const changeSet = ancestor.changes
                    ? parseChangeType(ancestor.changes)
                    : null;
                  return (
                    <div
                      key={ancestor.id}
                      className={classNames(
                        "flex items-center gap-3 p-2 rounded-md",
                        {
                          "ring-2 ring-indigo-500": isCurrent,
                        },
                      )}
                    >
                      <img
                        src={absoluteURL(ancestor.url)}
                        alt={ancestor.name}
                        className="w-12 h-12 object-cover rounded flex-shrink-0"
                      />
                      <div className="min-w-0 flex-1">
                        {isCurrent ? (
                          <span
                            className={classNames(
                              "text-sm font-medium truncate block",
                              TEXT_COLOR,
                            )}
                          >
                            {ancestor.name}
                          </span>
                        ) : (
                          <Link
                            to={routes.profile.image(ancestor.id)}
                            className="text-sm font-medium truncate block text-indigo-500 hover:text-indigo-400"
                          >
                            {ancestor.name}
                          </Link>
                        )}
                        <span
                          className={classNames(
                            "text-xs",
                            SECONDARY_TEXT_COLOR_DIM,
                          )}
                        >
                          {i === 0
                            ? t("imageDetail.original")
                            : changeSet
                              ? changeSet
                              : t("imageDetail.editFallback")}
                        </span>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function CopyURLs({ url, name }: { url: string; name: string }) {
  const { t } = useTranslation();
  const formats = React.useMemo(
    () => [
      { label: "URL", value: url },
      { label: "HTML", value: `<img src="${url}" alt="${name}">` },
      { label: "Markdown", value: `![${name}](${url})` },
      { label: "BBCode", value: `[img]${url}[/img]` },
    ],
    [url, name],
  );

  return (
    <div className="mt-4">
      <h2 className={classNames("text-lg font-medium mb-3", TEXT_COLOR)}>
        {t("imageDetail.copyReferences")}
      </h2>
      <div className={classNames("space-y-2 rounded-md p-4", SECOND_LAYER)}>
        {formats.map((fmt) => (
          <CopyRow key={fmt.label} label={fmt.label} value={fmt.value} />
        ))}
      </div>
    </div>
  );
}

function CopyRow({ label, value }: { label: string; value: string }) {
  const { t } = useTranslation();
  const [justCopied, setJustCopied] = React.useState(false);
  const timeoutRef = React.useRef<ReturnType<typeof setTimeout> | undefined>(
    undefined,
  );

  const handleCopy = React.useCallback(() => {
    copyText(value, () => {
      setJustCopied(true);
      clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(() => setJustCopied(false), 1500);
    });
  }, [value]);

  React.useEffect(() => {
    return () => clearTimeout(timeoutRef.current);
  }, []);

  return (
    <div className="flex items-center gap-2">
      <span
        className={classNames(
          "text-sm w-20 flex-shrink-0",
          SECONDARY_TEXT_COLOR_DIM,
        )}
      >
        {label}
      </span>
      <code
        className={classNames(
          "text-xs flex-1 min-w-0 truncate px-2 py-1 rounded " + BASE_LAYER,
          TEXT_COLOR,
        )}
        title={value}
      >
        {value}
      </code>
      <Button
        variant={justCopied ? "green" : "secondary"}
        noPadding
        className="px-2 py-1 text-xs flex-shrink-0"
        onClick={handleCopy}
        disabled={justCopied}
      >
        {justCopied ? t("common.toast.copied") : t("imageDetail.copy")}
      </Button>
    </div>
  );
}

function parseChangeType(changesJson: string): string | null {
  try {
    const parsed = JSON.parse(changesJson);
    if (parsed.type) {
      return parsed.type.charAt(0).toUpperCase() + parsed.type.slice(1);
    }
  } catch {
    // ignore
  }
  return null;
}
