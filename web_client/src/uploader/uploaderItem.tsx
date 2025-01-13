import React from "react";

import { HiOutlineClipboardDocumentList, HiCheck } from "react-icons/hi2";
import { IoEllipsisVertical } from "react-icons/io5";
import classNames from "classnames";
import { SECONDARY_TEXT_COLOR_DIM, SECOND_LAYER } from "~src/ui/classNames";
import copy from "copy-to-clipboard";
import { Button } from "~src/ui/button";
import { MenuWithTrigger } from "~src/ui/menuWithTrigger";
import { Loader } from "~src/ui/loader";
import { MenuSections } from "~src/ui/menu";
import { toast } from "react-toastify";
import { useTranslation } from "react-i18next";
import { copyText } from "~src/lib/copyText";

export type UploaderItem = {
  id: string;
  file: File;
  totalSize: number;
  loadedSize: number;
  loaded: boolean;
  errored: boolean;
  aborted: boolean;
  uploadedFileName?: string;
  url?: string;
};

export function UploaderItem({
  uploadingFile,
}: {
  uploadingFile: UploaderItem;
}) {
  const { t } = useTranslation();
  const { file, totalSize, loadedSize, loaded, errored, aborted, url } =
    uploadingFile;
  const lastCopyTs = React.useRef<number>(0);
  const [justCopied, setJustCopied] = React.useState(false);
  const scheduleSetNotJustCopied = React.useCallback(() => {
    setTimeout(() => {
      if (new Date().getTime() - lastCopyTs.current < 1000) {
        scheduleSetNotJustCopied();
        return;
      }
      setJustCopied(false);
    }, 1000);
  }, [setJustCopied]);
  const colorClass = React.useMemo(() => {
    if (errored || aborted) return "bg-red-600 dark:bg-red-500";
    if (justCopied) return "bg-green-600 dark:bg-green-500";
    if (loaded)
      return "bg-indigo-600 hover:bg-indigo-700 dark:bg-indigo-700 dark:hover:bg-indigo-600";
    return "bg-neutral-600 dark:bg-neutral-500";
  }, [errored, loaded, aborted, justCopied]);
  const widthPercentage = React.useMemo(() => {
    if (loaded) return 100;
    if (errored || aborted) return 100;
    return Math.floor((loadedSize * 100) / totalSize);
  }, [loaded, errored, aborted, loadedSize, totalSize]);
  const onClick = React.useCallback(() => {
    if (loaded && url) {
      if (!url) return;
      copy(url);
      lastCopyTs.current = new Date().getTime();
      setJustCopied(true);
      scheduleSetNotJustCopied();
    }
  }, [loaded, scheduleSetNotJustCopied, url]);

  const menuSections: MenuSections = React.useMemo(() => {
    if (!loaded || !url) {
      return { children: [] };
    }
    const fileName = file.name || "image";
    const copiedCb = () => toast(t("common.toast.copied"));
    const copyURL = () => copyText(url, copiedCb);
    const copyHTML = () =>
      copyText(`<img src="${url}" alt="${fileName}">`, copiedCb);
    const copyBB = () => copyText(`[img]${url}[/img]`, copiedCb);
    const copyMarkdown = () => copyText(`![${fileName}](${url})`, copiedCb);
    return {
      children: [
        {
          id: "main",
          items: [
            {
              id: "copyURL",
              action: copyURL,
              children: t("uploader.copyFormat", { format: "URL" }),
            },
            {
              id: "html",
              action: copyHTML,
              children: t("uploader.copyFormat", { format: "HTML" }),
            },
            {
              id: "bb",
              action: copyBB,
              children: t("uploader.copyFormat", { format: "BBCode" }),
            },
            {
              id: "markdown",
              action: copyMarkdown,
              children: t("uploader.copyFormat", { format: "Markdown" }),
            },
          ],
        },
      ],
    };
  }, [loaded, url, file.name, t]);

  return (
    <div className="uploder-item mb-2 flex items-center gap-1 relative">
      <div
        className={classNames(
          "w-full bg-gray-200 rounded-md h-8 relative flex items-center p-3 ",
          SECOND_LAYER,
          { "cursor-pointer": loaded }
        )}
        onClick={onClick}
      >
        <div
          className={classNames(
            "h-8 rounded-md absolute top-0 left-0 duration-200",
            colorClass
          )}
          style={{ width: widthPercentage + "%" }}
        ></div>
        <div className="text-white z-10 max-w-full pr-4 flex items-center gap-2 pointer-events-none shrink-0">
          {loaded ? (
            justCopied ? (
              <HiCheck className="shrink-0" />
            ) : (
              <HiOutlineClipboardDocumentList className="shrink-0" />
            )
          ) : null}
          <span className="whitespace-nowrap text-ellipsis overflow-hidden ">
            {errored || aborted ? "[Error]" : ""}
            {file.name}
          </span>
        </div>
      </div>
      <div className="uplodaer-item-actions h-full absolute right-0">
        {loaded ? (
          <MenuWithTrigger
            containerClassName="h-full"
            trigger={
              <Button
                variant="transparent"
                noPadding
                className={classNames(
                  "px-[2px] py-[2px]",
                  "hover:bg-white/20 h-full",
                  "text-white"
                )}
              >
                <IoEllipsisVertical size={20} />
              </Button>
            }
            menuSections={menuSections}
          />
        ) : !errored && !aborted ? (
          <div
            className={SECONDARY_TEXT_COLOR_DIM + " h-full flex items-center"}
          >
            <Loader />
          </div>
        ) : null}
      </div>
    </div>
  );
}
