export function absoluteURL(url: string) {
  if (url.startsWith("http")) {
    return url;
  }
  return `${window.location.origin}${url}`;
}
