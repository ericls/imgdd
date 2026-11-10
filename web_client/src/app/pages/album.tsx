import React from "react";
import { Link, useParams } from "react-router";
import { useQuery } from "@apollo/client/react";
import classNames from "classnames";
import { useTranslation } from "react-i18next";
import { gql } from "~src/__generated__";
import { PublicAlbumImageFieldsFragment } from "~src/__generated__/graphql";
import {
  HEADING_2,
  LINK_COLOR,
  SECOND_LAYER,
  SECONDARY_TEXT_COLOR_DIM,
  TEXT_COLOR,
} from "~src/ui/classNames";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";
import { useHumanizeDateTime } from "~src/lib/humanizeDateTime";
import { absoluteURL } from "~src/lib/url";
import { routes } from "~src/routes";

// How often the page checks for new images while it is visible.
const POLL_INTERVAL_MS = 5000;
// How close to the bottom (in px) still counts as following the latest image.
const FOLLOW_THRESHOLD_PX = 200;

const _PublicAlbumImageFields = gql(`
  fragment PublicAlbumImageFields on PublicAlbumImage {
    caption
    image {
      id
      url
      name
      nominalWidth
      nominalHeight
      createdAt
    }
  }
`);

const PublicAlbumDoc = gql(`
  query PublicAlbum($id: ID!) {
    publicAlbum(id: $id) {
      id
      title
      createdAt
      truncated
      images {
        ...PublicAlbumImageFields
      }
    }
  }
`);

function isNearBottom() {
  const doc = document.documentElement;
  return (
    doc.scrollHeight - (window.scrollY + window.innerHeight) <
    FOLLOW_THRESHOLD_PX
  );
}

function scrollToBottom(behavior: ScrollBehavior) {
  window.scrollTo({ top: document.documentElement.scrollHeight, behavior });
}

export function AlbumPage() {
  const { albumId } = useParams<{ albumId: string }>();
  const { t } = useTranslation();
  const { data, error } = useQuery(PublicAlbumDoc, {
    variables: { id: albumId ?? "" },
    skip: !albumId,
    pollInterval: POLL_INTERVAL_MS,
    // Don't poll from background tabs.
    skipPollAttempt: () => document.hidden,
  });
  const album = data?.publicAlbum;
  const images = album?.images;
  const imageCount = images?.length ?? 0;

  // Follow new images like a log: keep the latest in view while the reader
  // is at the bottom, otherwise offer a button instead of jumping.
  const following = React.useRef(true);
  const seenCount = React.useRef<number | null>(null);
  const [unseenCount, setUnseenCount] = React.useState(0);
  React.useEffect(() => {
    const onScroll = () => {
      following.current = isNearBottom();
      if (following.current) setUnseenCount(0);
    };
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);
  React.useEffect(() => {
    if (images === undefined) return;
    const previous = seenCount.current;
    seenCount.current = images.length;
    if (previous === null) {
      // First load: start at the latest image.
      if (images.length > 0) scrollToBottom("auto");
      return;
    }
    const added = images.length - previous;
    if (added <= 0) return;
    if (following.current) {
      scrollToBottom("smooth");
    } else {
      setUnseenCount((n) => n + added);
    }
  }, [images]);

  if (!data && !error) return <FullScreenLoader />;
  if (error && !data)
    return (
      <div className={classNames("p-8", TEXT_COLOR)}>
        {t("album.errorLoading")}: {error.message}
      </div>
    );
  if (!album)
    return (
      <div className={classNames("p-8", TEXT_COLOR)}>{t("album.notFound")}</div>
    );

  return (
    <div className="w-full max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
      <header className="mb-6">
        <h1 className={classNames(HEADING_2, "font-poppins break-words")}>
          {album.title || t("album.untitled")}
        </h1>
        <div
          className={classNames(
            "mt-1 flex items-center gap-2 text-sm",
            SECONDARY_TEXT_COLOR_DIM,
          )}
        >
          <span>{t("album.imageCount", { count: imageCount })}</span>
          <span aria-hidden>&middot;</span>
          <span className="inline-flex items-center gap-1.5">
            <span className="relative flex h-2 w-2" aria-hidden>
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
            </span>
            {t("album.live")}
          </span>
        </div>
        {album.truncated && (
          <p className={classNames("mt-2 text-sm", SECONDARY_TEXT_COLOR_DIM)}>
            {t("album.truncated", { count: imageCount })}
          </p>
        )}
      </header>

      {imageCount === 0 ? (
        <p
          className={classNames(
            "rounded-md p-8 text-center",
            SECOND_LAYER,
            SECONDARY_TEXT_COLOR_DIM,
          )}
        >
          {t("album.empty")}
        </p>
      ) : (
        <ol className="space-y-6">
          {images?.map((item, index) => (
            <AlbumImage key={item.image.id} item={item} position={index + 1} />
          ))}
        </ol>
      )}

      {unseenCount > 0 && (
        <button
          type="button"
          onClick={() => {
            setUnseenCount(0);
            scrollToBottom("smooth");
          }}
          className="fixed bottom-6 left-1/2 -translate-x-1/2 rounded-full bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-lg hover:bg-indigo-500"
        >
          {t("album.newImages", { count: unseenCount })} &darr;
        </button>
      )}
    </div>
  );
}

function AlbumImage({
  item,
  position,
}: {
  item: PublicAlbumImageFieldsFragment;
  position: number;
}) {
  const { t } = useTranslation();
  const { image, caption } = item;
  const url = absoluteURL(image.url);
  const uploaded = useHumanizeDateTime({ datetimeStr: image.createdAt });
  return (
    <li
      id={`album-image-${image.id}`}
      className={classNames("rounded-md p-3", SECOND_LAYER)}
    >
      <div className="mb-2 flex items-baseline gap-3">
        <span
          className={classNames(
            "font-mono text-sm shrink-0",
            SECONDARY_TEXT_COLOR_DIM,
          )}
        >
          #{position}
        </span>
        {/* Captions are plain text: never rendered as HTML or links. */}
        <p
          className={classNames(
            "min-w-0 grow whitespace-pre-wrap break-words",
            TEXT_COLOR,
          )}
        >
          {caption}
        </p>
        <time
          dateTime={image.createdAt}
          title={new Date(image.createdAt).toLocaleString()}
          className={classNames("shrink-0 text-sm", SECONDARY_TEXT_COLOR_DIM)}
        >
          {uploaded}
        </time>
      </div>
      <a href={url} target="_blank" rel="noopener noreferrer">
        <img
          src={url}
          alt={caption || image.name}
          loading="lazy"
          width={image.nominalWidth}
          height={image.nominalHeight}
          className="mx-auto block h-auto max-w-full rounded ring-1 ring-neutral-200 dark:ring-neutral-700"
        />
      </a>
      <div className="mt-2 flex gap-4 text-sm">
        <a
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          className={LINK_COLOR}
        >
          {t("album.openFullSize")}
        </a>
        <Link to={routes.image(image.id)} className={LINK_COLOR}>
          {t("album.details")}
        </Link>
      </div>
    </li>
  );
}
