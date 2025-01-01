import React from "react";

import {
  HiOutlineCloudArrowUp,
  HiOutlineClipboardDocumentList,
  HiCheck,
} from "react-icons/hi2";
import classNames from "classnames";
import { TEXT_COLOR } from "~src/ui/classNames";
import { uniqueId } from "lodash-es";
import copy from "copy-to-clipboard";
import { Button } from "~src/ui/button";
import { UploaderItem } from "./uploaderItem";
import { useTranslation } from "react-i18next";
import { addSessionHeaderToXMLHttpRequest } from "~src/lib/sessionToken";

const FILE_SIZE_LIMIT = 5e6;

export type UploadingFile = {
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

export function Uplodaer() {
  const { t } = useTranslation();
  const [uploadingFiles, setUploadingFiles] = React.useState<UploadingFile[]>(
    []
  );
  const inputRef = React.useRef<HTMLInputElement>(null);
  const uploadSingle = React.useCallback(
    (file: File) => {
      const id = uniqueId("file-");
      const newUploadingFile: UploadingFile = {
        id,
        file,
        totalSize: 1,
        loadedSize: 0,
        errored: false,
        aborted: false,
        loaded: false,
      };
      setUploadingFiles((current) => [...current, newUploadingFile]);
      function setCurrentFile(data: Omit<Partial<UploadingFile>, "id">) {
        setUploadingFiles((current) => {
          return current.map((item) => {
            if (item.id === id) {
              return { ...item, ...data };
            }
            return item;
          });
        });
      }
      function progressHandler(e: ProgressEvent) {
        const { total, loaded } = e;
        setCurrentFile({ totalSize: total, loadedSize: loaded });
      }
      const request = new XMLHttpRequest();
      request.upload.addEventListener("progress", progressHandler, {
        once: false,
      });
      request.addEventListener("load", () => {
        const resp = request.responseText;
        let filename: string;
        let url: string;
        try {
          const payload: { filename: string; url: string } = JSON.parse(resp);
          filename = payload.filename;
          url = payload.url;
        } catch {
          setCurrentFile({ errored: true });
          return;
        }
        if (filename && url) {
          // if URL is full URL, use that URL
          // otherwise build one based on current URL
          if (!url.startsWith("http")) {
            url = `${window.location.origin}${url}`;
            console.log(url);
          }
          setCurrentFile({
            loaded: true,
            uploadedFileName: filename,
            url,
          });
        } else {
          setCurrentFile({ errored: true });
        }
      });
      request.addEventListener(
        "error",
        () => {
          console.log("error");
          setCurrentFile({ errored: true });
        },
        { once: false }
      );
      request.addEventListener(
        "abort",
        () => setCurrentFile({ aborted: true }),
        { once: false }
      );
      const formdata = new FormData();
      formdata.append("image", file, file.name);
      request.open("POST", "/upload");
      addSessionHeaderToXMLHttpRequest(request);
      request.send(formdata);
    },
    [setUploadingFiles]
  );
  const onFileChange = React.useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const files = e.target.files;
      if (!files) return;
      let i = 0;
      let c = 0;
      while (i < files.length && c < 5) {
        const file = files.item(i);
        if (!file) {
          continue;
        }
        if (file.size <= FILE_SIZE_LIMIT && !file.type.includes("svg")) {
          uploadSingle(file);
          c += 1;
        }
        i += 1;
      }
    },
    [uploadSingle]
  );
  const onAreaClick = React.useCallback(() => {
    inputRef.current?.click();
  }, []);
  const [dragActive, setDragActive] = React.useState(false);
  const handleDrag = React.useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  }, []);
  const handlePaste = React.useCallback(
    (e: Event) => {
      if (!(e instanceof ClipboardEvent)) return;
      const files = e.clipboardData?.files;
      if (!files?.length) return;
      let count = 0;
      for (const file of files) {
        if (count >= 5) {
          break;
        }
        if (
          file.type.startsWith("image/") &&
          file.size <= FILE_SIZE_LIMIT &&
          !file.type.includes("svg")
        ) {
          uploadSingle(file);
          count += 1;
        }
      }
    },
    [uploadSingle]
  );
  React.useEffect(() => {
    const isFirefox = navigator.userAgent.toLowerCase().indexOf("firefox") > -1;
    if (isFirefox) return; // Firefox have strange behavior with pasted images
    window.document.body.addEventListener("paste", handlePaste);
    return () => window.document.body.removeEventListener("paste", handlePaste);
  }, [handlePaste]);
  React.useEffect(() => {
    window.document.body.focus();
  }, []);
  const handleDrop = React.useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setDragActive(false);
      if (e.dataTransfer.files) {
        let count = 0;
        for (const file of e.dataTransfer.files) {
          if (count >= 5) {
            break;
          }
          if (
            file.type.startsWith("image/") &&
            file.size <= FILE_SIZE_LIMIT &&
            !file.type.includes("svg")
          ) {
            uploadSingle(file);
            count += 1;
          }
        }
      }
    },
    [uploadSingle]
  );
  const hasMoreThanOneUploaded = React.useMemo(() => {
    let count = 0;
    for (const f of uploadingFiles) {
      if (f.loaded) {
        count += 1;
      }
      if (count == 2) return true;
    }
    return false;
  }, [uploadingFiles]);
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
  const onCopyAllURLs = React.useCallback(() => {
    const urls = [];
    for (const f of uploadingFiles) {
      if (f.url) {
        urls.push(f.url);
      }
    }
    const text = urls.join("\n");
    copy(text);
    setJustCopied(true);
    scheduleSetNotJustCopied();
  }, [scheduleSetNotJustCopied, uploadingFiles]);
  return (
    <div
      className="uploader max-w-[96%] min-[500px]:max-w-[480px] focus:outline-none focus:ring-0"
      tabIndex={-1}
    >
      <div
        className={classNames(
          TEXT_COLOR,
          "rounded-md mb-6 duration-200",
          "uploding-area cursor-pointer p-10",
          "bg-neutral-50 dark:bg-neutral-800",
          "hover:bg-neutral-100 hover:dark:bg-neutral-900",
          "border border-dashed border-neutral-300",
          { ["blur-[1px]"]: dragActive }
        )}
        onClick={onAreaClick}
        onDragEnter={handleDrag}
        onDragOver={handleDrag}
        onDragLeave={handleDrag}
        onDrop={handleDrop}
      >
        <div className="flex flex-col items-center">
          <input
            type="file"
            ref={inputRef}
            className="hidden"
            onChange={onFileChange}
            accept="image/*"
            multiple
          />
          <HiOutlineCloudArrowUp className="mb-2" size={48} />
          <p className="text-center">{t("uploader.mainHelpText")}</p>
        </div>
      </div>
      <div>
        <div className="mb-2 flex justify-end">
          {hasMoreThanOneUploaded && (
            <Button
              variant={justCopied ? "green" : "indigo"}
              roundLevel="md"
              onClick={onCopyAllURLs}
              className="flex items-center gap-2 px-2 py-1"
            >
              {justCopied ? <HiCheck /> : <HiOutlineClipboardDocumentList />}
              {t("uploader.copyAllURLs")}
            </Button>
          )}
        </div>
        {uploadingFiles.map((uf) => {
          return <UploaderItem key={uf.id} uploadingFile={uf} />;
        })}
      </div>
    </div>
  );
}
