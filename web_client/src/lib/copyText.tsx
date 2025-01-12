import copy from "copy-to-clipboard";

export function copyText(text: string, cb?: () => void) {
  if (!text) return;
  copy(text);
  cb?.();
}
